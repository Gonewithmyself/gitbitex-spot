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

package matching

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gitbitex/gitbitex-spot/models"
	"github.com/gitbitex/gitbitex-spot/service"
	"github.com/shopspring/decimal"
)

var count int64

func Test_orderBook_ApplyOrder(t *testing.T) {

	pd, _ := service.GetProductById("BTC-USDT")
	ob := NewOrderBook(pd)

	pjson(ob.ApplyOrder(sell(90, 0.1, "")))
	pjson(ob.ApplyOrder(buy(100, 0.2)))

	pjson(ob.ApplyOrder(buy(99, 0.1, "")))
	pjson(ob.ApplyOrder(sell(99, 0.2)))
	// pjson(ob.ApplyOrder(sell(20, 5)))
	t.Error()
	_ = ob
}

func Test_orderBook_CancelOrder(t *testing.T) {
	pd, _ := service.GetProductById("BTC-USDT")
	ob := NewOrderBook(pd)

	od := sell(99, 0.2, "")
	pjson(ob.ApplyOrder(od))
	pjson(ob.ApplyOrder(buy(99, 0.1, "")))
	pjson(ob.CancelOrder(od))

	t.Log(len(ob.depths["buy"].orders), len(ob.depths["sell"].orders))
	od = buy(99, 0.2, "")
	pjson(ob.ApplyOrder(od))
	pjson(ob.ApplyOrder(sell(99, 0.1, "")))
	pjson(ob.CancelOrder(od))

	t.Error()
}

func pjson(arg interface{}) {
	d, _ := json.Marshal(arg)
	fmt.Println(string(d))
}

func buy(price, size float64, typ ...string) *models.Order {
	count++
	o := &models.Order{
		Id:        count,
		UserId:    1,
		Type:      "market",
		Side:      "buy",
		ProductId: "BTC-USDT",
	}
	if len(typ) > 0 {
		o.Type = "limit"
		o.Size = decimal.NewFromFloat(size)
		o.Price = decimal.NewFromFloat(price)
	}

	o.Funds = decimal.NewFromFloat(price).Mul(decimal.NewFromFloat(size))
	return o
}

func sell(price, size float64, typ ...string) *models.Order {
	count++
	o := &models.Order{
		Id:        count,
		UserId:    1,
		Type:      "market",
		Side:      "sell",
		ProductId: "BTC-USDT",
	}
	o.Size = decimal.NewFromFloat(size)

	if len(typ) > 0 {
		o.Type = "limit"
		o.Funds = o.Size.Mul(o.Price)
		o.Price = decimal.NewFromFloat(price)
	}
	return o
}
