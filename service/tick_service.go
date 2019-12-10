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
	"github.com/gitbitex/gitbitex-spot/models"
	"github.com/gitbitex/gitbitex-spot/models/mysql"
	"math"
)

func GetLastTickByProductId(productId string, granularity int64) (*models.Tick, error) {
	return mysql.SharedStore(productId).GetLastTickByProductId(productId, granularity)
}

func GetTicksByProductId(productId string, granularity int64, limit int) ([]*models.Tick, error) {
	return mysql.SharedStore(productId).GetTicksByProductId(productId, granularity, limit)
}

func AddTicks(productId string, ticks map[int64][]*models.Tick) (err error) {
	// return mysql.SharedStore().AddTicks(ticks)
	// var min := math.
	var (
		min int = math.MaxInt32
	)

	for k := range ticks {
		if l := len(ticks[k]); l < min {
			min = l
		}
	}

	last := ticks[1][min]
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

	for k := range ticks {
		if err = tx.AddTicks(ticks[k][:min]); err != nil {
			return
		}
	}

	if err = tx.UpsertOffset(&models.Offset{
		Group:     "kline_" + last.ProductId,
		LogOffset: last.LogOffset,
		LogSeq:    last.LogSeq,
	}); err != nil {
		return err
	}

	return
}
