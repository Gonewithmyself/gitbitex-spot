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
	"testing"
	"time"

	"github.com/gitbitex/gitbitex-spot/service"
)

func TestNewEngine(t *testing.T) {
	p, err := service.GetProductById("EOS-USDT")
	if err != nil {
		t.Fatal(err)
	}

	snapshotStore := NewRedisSnapshotStore(p.Id)
	snap, err := snapshotStore.GetLatest()
	if err != nil {
		t.Fatal(err)
	}

	bk := NewOrderBook(p)
	bk.Restore(&snap.OrderBookSnapshot)

	// snapshot
	start := time.Now()
	sn := bk.Snapshot()
	t.Log(time.Now().Sub(start).Milliseconds(), len(sn.Orders))

	// compress
	t.Error(p)
}

func TestCompress(t *testing.T) {

}
