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
)

// var ids = map[string]uint64{}

func GetProductById(id string) (*models.Product, error) {
	return mysql.SharedStore(db_account).GetProductById(id)
}

func GetProducts() ([]*models.Product, error) {
	return mysql.SharedStore(db_account).GetProducts()
}

// func ProductID(id string) uint64 {
// 	return ids[id]
// }

// func init() {
// 	pds, err := mysql.SharedStore().GetProducts()
// 	if err != nil {
// 		log.Fatal("load products", err)
// 	}
// 	for i := range pds {
// 		ids[pds[i].Id] = uint64(i)
// 	}
// }
