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

package rest

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitbitex/gitbitex-spot/pushing"
	"github.com/gitbitex/gitbitex-spot/rest/monitor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/siddontang/go-log/log"
	"golang.org/x/net/context"
)

type HttpServer struct {
	addr string
	s    *http.Server
}

func NewHttpServer(addr string) *HttpServer {
	return &HttpServer{
		addr: addr,
	}
}

func (server *HttpServer) Stop() {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*30))
	log.Info("stopping rest server")
	err := server.s.Shutdown(ctx)
	if err != nil {
		log.Error("stop timeout")
	}
}

func (server *HttpServer) Start() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.Default()
	r.Use(setCROSOptions)

	r.GET("/api/configs", GetConfigs)
	r.POST("/api/users", SignUp)
	r.POST("/api/users/accessToken", SignIn)
	r.POST("/api/users/token", GetToken)
	r.GET("/api/products", GetProducts)
	r.GET("/api/products/:productId/trades", GetProductTrades)
	r.GET("/api/products/:productId/book", GetProductOrderBook)
	r.GET("/api/products/:productId/candles", GetProductCandles)

	// monitoring
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// ws
	ws := pushing.NewBroker()
	r.GET("/ws", ws.Ws)

	private := r.Group("/", checkToken())
	{
		private.GET("/api/orders", GetOrders)
		private.POST("/api/orders", monitor.Wrap(PlaceOrder))
		private.DELETE("/api/orders/:orderId", CancelOrder)
		private.DELETE("/api/orders", CancelOrders)
		private.GET("/api/accounts", GetAccounts)
		private.GET("/api/users/self", GetUsersSelf)
		private.POST("/api/users/password", ChangePassword)
		private.DELETE("/api/users/accessToken", SignOut)
		private.GET("/api/wallets/:currency/address", GetWalletAddress)
		private.GET("/api/wallets/:currency/transactions", GetWalletTransactions)
		private.POST("/api/wallets/:currency/withdrawal", Withdrawal)
	}

	server.s = &http.Server{
		Addr:    server.addr,
		Handler: r,
	}

	err := server.s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func setCROSOptions(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Content-Type", "application/json")
}
