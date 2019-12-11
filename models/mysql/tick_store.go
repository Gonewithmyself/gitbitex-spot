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
	"strings"

	"github.com/gitbitex/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetTicksByProductId(productId, granularity string, limit int) ([]*models.Tick, error) {
	db := s.db.Order("time DESC").Limit(limit)
	var ticks []*models.Tick
	err := db.Table(granularity).Find(&ticks).Error
	return ticks, err
}

func (s *Store) GetLastTickByProductId(productId, granularity string) (*models.Tick, error) {
	var tick models.Tick
	sql := fmt.Sprintf("SELECT * FROM %s ORDER BY time DESC LIMIT 1", granularity)
	err := s.db.Raw(sql).Scan(&tick).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &tick, err
}

func (s *Store) AddTicks(tbname string, ticks []*models.Tick) error {
	if len(ticks) == 0 {
		return nil
	}

	var valueStrings []string
	for _, tick := range ticks {
		valueString := fmt.Sprintf("(%v, %v, %v, %v, %v, %v, %v)",
			"CURRENT_TIMESTAMP", tick.Time, tick.Open, tick.Low, tick.High, tick.Close,
			tick.Volume)
		valueStrings = append(valueStrings, valueString)
	}
	sql := fmt.Sprintf("insert INTO %s (created_at, time,open,low,high,close,"+
		"volume) VALUES %s"+
		"ON DUPLICATE KEY UPDATE `open`=VALUES(`open`), `low`=VALUES(`low`), `high`=VALUES(`high`), `close`=VALUES(`close`),`volume`=VALUES(`volume`)",
		tbname, strings.Join(valueStrings, ","))
	return s.db.Exec(sql).Error
}
