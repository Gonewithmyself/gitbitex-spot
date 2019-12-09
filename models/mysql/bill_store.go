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
	"time"

	"github.com/gitbitex/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetUnsettledBillsByUserId(userId int64, currency string) ([]*models.Bill, error) {
	db := s.db.Where("settled =?", 0).Where("user_id=?", userId).
		Where("currency=?", currency).Order("id ASC").Limit(100)

	var bills []*models.Bill
	err := db.Find(&bills).Error
	return bills, err
}

func (s *Store) GetLastBillByProductId(productId string) (*models.Bill, error) {
	var bill models.Bill
	err := s.db.Where("product_id =?", productId).Order("id DESC").Limit(1).Find(&bill).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &bill, err
}

func (s *Store) GetUnsettledBills() ([]*models.Bill, error) {
	db := s.db.Where("settled =?", 0).Order("id ASC").Limit(100)

	var bills []*models.Bill
	err := db.Find(&bills).Error
	return bills, err
}

func (s *Store) AddBills(bills []*models.Bill) error {
	if len(bills) == 0 {
		return nil
	}
	var valueStrings []string
	for _, bill := range bills {
		valueString := fmt.Sprintf("(%v,%v, '%v', %v, %v, '%v', %v, '%v')",
			"CURRENT_TIMESTAMP", bill.UserId, bill.Currency, bill.Available, bill.Hold, bill.Type, bill.Settled, bill.Notes)
		valueStrings = append(valueStrings, valueString)
	}
	sql := fmt.Sprintf("INSERT INTO g_bill (created_at, user_id,currency,available,hold, type,settled,notes) VALUES %s", strings.Join(valueStrings, ","))
	return s.db.Exec(sql).Error
}

func (s *Store) AddOffsetBills(bills []*models.OffsetBill) error {
	l := len(bills)
	if l == 0 {
		return nil
	}
	var valueStrings = make([]string, 0, l)
	for _, bill := range bills {
		valueString := fmt.Sprintf("(%v,%v, '%v', %v, %v, '%v', %v, '%v')",
			"CURRENT_TIMESTAMP", bill.UserId, bill.Currency, bill.Available, bill.Hold, bill.Type, bill.Settled, bill.Notes)
		valueStrings = append(valueStrings, valueString)
	}
	sql := fmt.Sprintf("INSERT INTO g_bill (created_at, user_id,currency,available,hold, type,settled,notes) VALUES %s", strings.Join(valueStrings, ","))
	return s.db.Exec(sql).Error
}

func (s *Store) GetOffsetForUpdate(group string, partition int64) (*models.Offset, error) {
	var offset models.Offset
	// sql := fmt.Sprintf("SELECT * FROM g_offset WHERE group='%s' AND partition=%d", group, partition)
	err := s.db.Raw("SELECT * FROM g_offset WHERE `group`=? AND `partition`=? FOR UPDATE",
		group, partition).Scan(&offset).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	// err := s.db.Where("group=?", group).Where("partition=?", partition).Find(&offset).Error
	return &offset, err
}

func (s *Store) AddOffset(account *models.Offset) error {
	return s.db.Create(account).Error
}

func (s *Store) UpdateOffset(account *models.Offset) error {
	return s.db.Save(account).Error
}

func (s *Store) GetLastOffset(group string, partition int64) (*models.Offset, error) {
	var offset models.Offset
	err := s.db.Where("`group`=? AND `partition=?`", group, partition).Find(&offset).Error
	if err == gorm.ErrRecordNotFound {
		return &offset, nil
	}
	return &offset, err
}

func (s *Store) UpdateBill(bill *models.Bill) error {
	bill.UpdatedAt = time.Now()
	return s.db.Save(bill).Error
}
