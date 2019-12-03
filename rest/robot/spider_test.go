package main

import (
	"testing"
	"time"

	"github.com/gitbitex/gitbitex-spot/worker"
)

func Test_response_get(t *testing.T) {
	args := []string{"BTC-USDT", "ETH-USDT"}

	for i := range args {
		t.Log(args[i], spider.get(args[i]))
	}

	t.Error()
}

func Test_getParams(t *testing.T) {
	t.Log(genParams("BTC-USDT"))
	t.Log(genParams("BTC-USDT"))
	t.Error()
}

func Test_place(t *testing.T) {
	// placeOrder()
	t.Error()
}

func Test_tickerUniqueKey(t *testing.T) {
	type args struct {
		product string
		idx     int
		ts      int64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		// TODO: Add test cases.
		{"1", args{"BTC-USDT", 1, time.Now().Unix()}, 1},
		{"1", args{"EOS-USDT", 1, time.Now().Unix()}, 1},
		{"1", args{"BTC-USDT", 30, time.Now().Unix()}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := worker.TickerUniqueKey(tt.args.product, tt.args.idx, tt.args.ts); got != tt.want {
				t.Errorf("tickerUniqueKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
