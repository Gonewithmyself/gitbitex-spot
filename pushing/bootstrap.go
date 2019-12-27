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

package pushing

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitbitex/gitbitex-spot/conf"
	"github.com/gitbitex/gitbitex-spot/matching"
	"github.com/gitbitex/gitbitex-spot/service"
	"github.com/gorilla/websocket"
	"github.com/siddontang/go-log/log"
)

func StartServer() {
	gbeConfig := conf.GetConfig()

	sub := newSubscription()

	newRedisStream(sub).Start()

	products, err := service.GetProducts()
	if err != nil {
		panic(err)
	}

	// TODO  get stream isolated?
	for _, product := range products {
		newTickerStream(product.Id, sub, matching.NewKafkaLogReader("tickerStream", product.Id, gbeConfig.Kafka.Brokers)).Start()
		newMatchStream(product.Id, sub, matching.NewKafkaLogReader("matchStream", product.Id, gbeConfig.Kafka.Brokers)).Start()
	}

	go NewServer(gbeConfig.PushServer.Addr, gbeConfig.PushServer.Path, sub).Run()

	log.Info("websocket server ok")
}

type Broker struct {
	sub *subscription
}

func NewBroker() *Broker {
	return &Broker{
		sub: newSubscription(),
	}
}

func (s *Broker) Start() {
	gbeConfig := conf.GetConfig()
	products, err := service.GetProducts()
	if err != nil {
		panic(err)
	}

	newRedisStream(s.sub).Start()
	for _, product := range products {
		newTickerStream(product.Id, s.sub, matching.NewKafkaLogReader("tickerStream", product.Id, gbeConfig.Kafka.Brokers)).Start()
		newMatchStream(product.Id, s.sub, matching.NewKafkaLogReader("matchStream", product.Id, gbeConfig.Kafka.Brokers)).Start()
	}
}

func (s *Broker) Ws(c *gin.Context) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error(err)
		return
	}

	NewClient(conn, s.sub).startServe()
}

func (s *Broker) Stop() {

}
