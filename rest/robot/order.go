package main

import (
	"math"
	"math/rand"

	"github.com/gitbitex/gitbitex-spot/models"
	"github.com/shopspring/decimal"
)

var (
	user int64 = 1
)

// RandInt 范围整数
func RandInt(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

// RandFloat  范围指定精度的浮点数
func RandFloat(min, max float64, n int) float64 {
	base := math.Pow10(n)
	iMin, iMax := int64(min*base), int64(max*base)

	return float64(RandInt(iMin, iMax)) / base
}

func placeOrder(pair string) {
	side, size, price, funds := genParams(pair)
	req := placeOrderRequest{
		Side:      string(side),
		Type:      string(models.OrderTypeLimit),
		ProductId: pair,
	}

	req.Size, _ = size.Float64()
	req.Price, _ = price.Float64()
	req.Funds, _ = funds.Float64()

	err := post(req)
	if err != nil {
		panic(err)
	}
	// od, err := service.PlaceOrder(user, "robot", pair, models.OrderTypeLimit,
	// 	side, size, price, funds)
	// if err != nil {
	// 	log.Fatal("place order: ", err)
	// }

	// rest.SubmitOrder(od)
}

func genParams(pair string) (side models.Side, size, price, funds decimal.Decimal) {
	usdt := RandFloat(1, 100, 1)
	p := round(spider.get(pair))
	price = decimal.NewFromFloat(p)
	funds = decimal.NewFromFloat(usdt)
	side = models.SideBuy
	size = funds.Div(price)
	if rand.Intn(2) == 0 {
		side = models.SideSell
	}
	return
}

func round(f float64) float64 {
	delta := RandFloat(0.01, 0.99, 2)
	if rand.Intn(2) == 0 {
		f -= delta
	} else {
		f += delta
	}
	return f

	// fn := math.Floor
	// if rand.Intn(2) == 0 {
	// 	fn = math.Ceil
	// }
	// return fn(f)
}
