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

package service

import (
	"math"
	"strconv"
	"sync"

	"github.com/gitbitex/gitbitex-spot/models"
	"github.com/gitbitex/gitbitex-spot/models/mysql"
)

func GetLastTickByProductId(productId string, granularity int64) (*models.Tick, error) {
	return mysql.SharedStore(productId).GetLastTickByProductId(productId, tbmap[granularity])
}

func GetTicksByProductId(productId string, granularity int64, limit int) ([]*models.Tick, error) {
	return mysql.SharedStore(productId).GetTicksByProductId(productId, tbmap[granularity], limit)
}

func AddTicks(productId string, ticks map[int64][]*models.Tick) (err error) {
	var (
		min int = math.MaxInt32
	)

	for k := range ticks {
		if l := len(ticks[k]); l < min {
			min = l
		}
	}

	last := ticks[1][min-1]
	tx, err := mysql.SharedStore(productId).BeginTx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.CommitTx()
		if err == nil {
			for k := range ticks {
				orig := ticks[k]
				left := orig[min:]
				ticks[k] = append(orig[:0], left...)
			}
		}

	}()

	list := make([]*models.Tick, 0, min)
	for k := range ticks {
		m := map[int64]struct{}{}
		ll := ticks[k][:min]
		for i := min - 1; i >= 0; i-- {
			if _, ok := m[ll[i].Time]; !ok {
				m[ll[i].Time] = struct{}{}
				list = append(list, ll[i])
			}
		}

		tb, ok := tbmap[k]
		if !ok {
			panic(k)
		}
		if err = tx.AddTicks(tb, list); err != nil {
			return
		}
		list = list[:0]
	}

	if err = tx.UpsertOffset(&models.Offset{
		Group:     "kline_" + productId,
		LogOffset: last.LogOffset,
		LogSeq:    last.LogSeq,
	}); err != nil {
		return err
	}

	return
}

var tbmap = map[int64]string{}
var one sync.Once

func InitTbmap(minutes []int64) {
	one.Do(func() {
		for _, m := range minutes {
			tbmap[m] = "g_tick_" + tickName(m)
		}
	})
}

func tickName(m int64) string {
	str := ""
	if m < 60 {

		if m != 1 {
			str = strconv.Itoa(int(m))
		}
		return "m" + str
	}

	h := m / 60
	if h < 24 {
		if h != 1 {
			str = strconv.Itoa(int(h))
		}
		return "h" + str
	}

	d := h / 24
	if d != 1 {
		str = strconv.Itoa(int(d))
	}
	return "d" + str
}
