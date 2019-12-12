// Copyright 2019 GitBitEx.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package worker

import (
	"math/rand"
	"time"

	"github.com/gitbitex/gitbitex-spot/matching"
	"github.com/gitbitex/gitbitex-spot/models"
	"github.com/gitbitex/gitbitex-spot/models/mysql"
	"github.com/gitbitex/gitbitex-spot/service"
	"github.com/shopspring/decimal"
	"github.com/siddontang/go-log/log"
)

var minutes = []int64{1, 3, 5, 15, 30, 60, 120, 240, 360, 720, 1440}

// func TickerUniqueKey(product string, idx int, ts int64) uint64 {
// 	id := service.ProductID(product)

// 	return id<<56 + uint64(idx)<<48 + uint64(ts)
// }

type TickMaker struct {
	ticks     map[int64]*models.Tick
	tickCh    chan []models.Tick
	logReader matching.LogReader
	logOffset int64
	logSeq    int64
}

func NewTickMaker(productId string, logReader matching.LogReader) *TickMaker {
	t := &TickMaker{
		ticks:     map[int64]*models.Tick{},
		tickCh:    make(chan []models.Tick, 100),
		logReader: logReader,
	}

	// 加载数据库中记录的最新tick
	for _, granularity := range minutes {
		tick, err := service.GetLastTickByProductId(productId, granularity)
		if err != nil {
			panic(err)
		}

		if tick != nil {
			log.Infof("load last tick: %v", tick)
			tick.Granularity = granularity
			t.ticks[granularity] = tick
		}
	}

	last, err := mysql.SharedStore(productId).GetLastOffset("kline_"+logReader.GetProductId(), 0)
	if err != nil {
		panic(err)
	}
	t.logOffset = last.LogOffset
	t.logSeq = last.LogSeq

	t.logReader.RegisterObserver(t)
	return t
}

func (t *TickMaker) Start() {
	if t.logOffset > 0 {
		t.logOffset++
	}
	go t.logReader.Run(t.logSeq, t.logOffset)
	go t.flusher()
}

func (t *TickMaker) OnOpenLog(log *matching.OpenLog, offset int64) {
	// do nothing
}

func (t *TickMaker) OnDoneLog(log *matching.DoneLog, offset int64) {
	// do nothing
}

func (t *TickMaker) OnMatchLog(log *matching.MatchLog, offset int64) {
	var list = make([]models.Tick, 0, len(minutes))
	for _, granularity := range minutes {
		tickTime := log.Time.UTC().Truncate(time.Duration(granularity) * time.Minute).Unix()

		tick, found := t.ticks[granularity]
		if !found || tick.Time != tickTime {
			tick = &models.Tick{
				Open:        log.Price,
				Close:       log.Price,
				Low:         log.Price,
				High:        log.Price,
				Volume:      log.Size,
				Time:        tickTime,
				Granularity: granularity,
				LogOffset:   offset,
				LogSeq:      log.Sequence,
			}
			t.ticks[granularity] = tick
		} else {
			tick.Close = log.Price
			tick.Low = decimal.Min(tick.Low, log.Price)
			tick.High = decimal.Max(tick.High, log.Price)
			tick.Volume = tick.Volume.Add(log.Size)
			tick.LogOffset = offset
			tick.LogSeq = log.Sequence
		}

		list = append(list, *tick)
	}
	t.tickCh <- list
}

func (t *TickMaker) flusher() {
	var tickM = make(map[int64][]*models.Tick, len(minutes))
	var count int
	pid := t.logReader.GetProductId()
	ticker := time.NewTicker(time.Second)

	flushfn := func() {
		for {
			err := service.AddTicks(pid, tickM)
			if err != nil {
				log.Error(err)
				// retry
				time.Sleep(time.Second + time.Duration(rand.Intn(2000)))
				continue
			}

			// TODO redis publish
			count = 0
			break
		}
	}
	for {
		select {
		case <-ticker.C:
			if count > 0 {
				flushfn()
			}

		case ticks := <-t.tickCh:
			count++
			for i := range ticks {
				tickM[ticks[i].Granularity] = append(tickM[ticks[i].Granularity], &ticks[i])
			}

			if count < 100 {
				continue
			}

			flushfn()
		}
	}
}

func init() {
	service.InitTbmap(minutes)
}
