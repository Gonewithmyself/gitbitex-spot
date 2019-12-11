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
	"testing"

	"github.com/gitbitex/gitbitex-spot/models"
	"github.com/gitbitex/gitbitex-spot/models/mysql"
	"github.com/siddontang/go-log/log"
)

func TestAddOffsetFills(t *testing.T) {
	bills := []*models.OffsetBill{
		// &models.OffsetBill{Offset: models.Offset{"test", 1, 1, 0}},
		&models.OffsetBill{Offset: models.Offset{"test", 0, 1, 1}},
		&models.OffsetBill{Offset: models.Offset{"test", 2, 1, 1}},
	}
	err := AddOffsetBills("", bills)
	log.Println(err)

	// a, err := mysql.SharedStore().
	// st := mysql.SharedStore()
	// t.Log(st.GetOffsetForUpdate("123", 0))
	//t.Log(st.GetAccount(41, "ETH"))

	//t.Log(st.AddOffset(&models.Offset{Group: "test"}))
	// t.Log(a)
	t.Error()
}
func TestXXXX(t *testing.T) {
	// bills := []*models.OffsetBill{
	// 	// &models.OffsetBill{Offset: models.Offset{"test", 1, 1, 0}},
	// 	&models.OffsetBill{Offset: models.Offset{"test", 0, 1, 1}},
	// }
	// err := AddOffsetFills(bills)
	// log.Println(err)

	// a, err := mysql.SharedStore().
	st := mysql.SharedStore(db_account)
	// t.Log(st.GetOffsetForUpdate("123", 0))
	t.Log(st.GetAccount(41, "ETH"))

	//t.Log(st.AddOffset(&models.Offset{Group: "test"}))
	// t.Log(a)
	t.Error()
}

func Test(t *testing.T) {
	orig := []int64{1, 2, 3, 4, 5, 6}
	min := len(orig) - 1
	t.Log(orig[:min])
	left := orig[min:]
	orig = append(orig[:0], left...)
	t.Error(orig)
}
