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

package worker

import "testing"

import "time"

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TickerUniqueKey(tt.args.product, tt.args.idx, tt.args.ts); got != tt.want {
				t.Errorf("tickerUniqueKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinte(t *testing.T) {
	now := time.Date(2019, 12, 8, 19, 20, 0, 0, time.Local)
	// now = time.Date(2019, 12, 02, 12,45, 0, 0, time.Local)
	for _, granularity := range minutes {
		tickTime := now.UTC().Truncate(time.Duration(granularity) * time.Minute)
		t.Log(tickTime.Local(), granularity)
	}
	t.Error(now)
}
