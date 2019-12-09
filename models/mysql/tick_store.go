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

package mysql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gitbitex/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetTicksByProductId(productId string, granularity int64, limit int) ([]*models.Tick, error) {
	db := s.db.Where("product_id =?", productId).Where("granularity=?", granularity).
		Order("time DESC").Limit(limit)
	var ticks []*models.Tick
	err := db.Find(&ticks).Error
	return ticks, err
}

func (s *Store) GetLastTickByProductId(productId string, granularity int64) (*models.Tick, error) {
	var tick models.Tick
	err := s.db.Raw("SELECT * FROM g_tick WHERE product_id=? AND granularity=? ORDER BY time DESC LIMIT 1",
		productId, granularity).Scan(&tick).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &tick, err
}

func (s *Store) AddTicks(ticks []*models.Tick) error {
	if len(ticks) == 0 {
		return nil
	}

	// ticks = uniqueTicks(ticks)
	var valueStrings []string
	for _, tick := range ticks {
		valueString := fmt.Sprintf("(%v, %v, %v, %v, %v, %v, %v,%v,%v)",
			"CURRENT_TIMESTAMP", tick.Time, tick.Open, tick.Low, tick.High, tick.Close,
			tick.Volume, tick.LogOffset, tick.LogSeq)
		valueStrings = append(valueStrings, valueString)
	}
	sql := fmt.Sprintf("insert INTO g_tick (created_at, time,open,low,high,close,"+
		"volume,log_offset,log_seq) VALUES %s"+
		"ON DUPLICATE KEY UPDATE `open`=VALUES(`open`), `low`=VALUES(`low`), `high`=VALUES(`high`), `close`=VALUES(`close`),`volume`=VALUES(`volume`),",
		strings.Join(valueStrings, ","))
	return s.db.Exec(sql).Error
}

func uniqueTicks(ticks []*models.Tick) []*models.Tick {
	if len(ticks) == 1 {
		return ticks
	}

	m := map[string]struct{}{}
	list := make([]*models.Tick, 0, len(ticks))
	for i := range ticks {
		k := ticks[i].ProductId + strconv.Itoa(int(ticks[i].Time)) + strconv.Itoa(int(ticks[i].Granularity))
		if _, ok := m[k]; !ok {
			list = append(list, ticks[i])
		}
	}
	return list
}
