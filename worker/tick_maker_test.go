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
	"testing"
	"time"

	"github.com/gitbitex/gitbitex-spot/matching"
	"github.com/gitbitex/gitbitex-spot/models"
	"github.com/gitbitex/gitbitex-spot/service"
	"github.com/shopspring/decimal"
)

func TestSaveTicks(t *testing.T) {
	tk := &TickMaker{
		ticks: make(map[int64]*models.Tick),
	}

	ticks := genticks(tk, 10)
	// utils.Pjson(ticks)

	tm := make(map[int64][]*models.Tick)
	for i := range ticks {
		tm[ticks[i].Granularity] = append(tm[ticks[i].Granularity], &ticks[i])
	}

	err := service.AddTicks("BTC-USDT", tm)
	t.Log(len(ticks), len(minutes))
	t.Error(err)
}

func genticks(tk *TickMaker, n int) (ticks []models.Tick) {
	now := time.Now()
	for i := 0; i < n; i++ {
		log := &matching.MatchLog{
			Price: decimal.NewFromFloat(99),
			Size:  decimal.NewFromFloat(1),
			Base: matching.Base{
				Time:     now,
				Sequence: int64(i + 1),
			},
		}
		ticks = append(ticks, tk.OnMatch(log, int64(i))...)
		now = now.Add(time.Second * 120)
	}
	return
}

func (t *TickMaker) OnMatch(log *matching.MatchLog, offset int64) (ticks []models.Tick) {

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

		ticks = append(ticks, *tick)
	}
	return
}
func init() {
	// cfg, err := conf.GetConfig()
	// if err != nil {
	// 	panic(err)
	// }

	// url := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local",
	// 	cfg.DataSource.User, cfg.DataSource.Password, cfg.DataSource.Addr, cfg.DataSource.Database)
	// gdb, err = gorm.Open(cfg.DataSource.DriverName, url)
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	log.Println("success")
	// }
	// log.Println(url)

}
