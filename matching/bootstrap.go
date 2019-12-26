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
	"github.com/gitbitex/gitbitex-spot/conf"
	"github.com/gitbitex/gitbitex-spot/service"
	"github.com/siddontang/go-log/log"
	"sync"
)

var engines []*Engine

func StartEngine() {
	gbeConfig := conf.GetConfig()
	products, err := service.GetProducts()
	if err != nil {
		panic(err)
	}
	for _, product := range products {
		orderReader := NewKafkaOrderReader(product.Id, gbeConfig.Kafka.Brokers)
		snapshotStore := NewRedisSnapshotStore(product.Id)
		logStore := NewKafkaLogStore(product.Id, gbeConfig.Kafka.Brokers)
		matchEngine := NewEngine(product, orderReader, logStore, snapshotStore)
		matchEngine.Start()
		engines = append(engines, matchEngine)
	}

	log.Info("match engine ok")
}

func StopEngine() {
	var wg sync.WaitGroup
	wg.Add(len(engines))
	for i := range engines {
		go func(i int) {
			engines[i].Stop()
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Info("stop engine ok")
}
