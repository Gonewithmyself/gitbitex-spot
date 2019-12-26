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

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gitbitex/gitbitex-spot/conf"
	"github.com/gitbitex/gitbitex-spot/models"
	"github.com/gitbitex/gitbitex-spot/service"
	"github.com/go-redis/redis"
	"github.com/siddontang/go-log/log"
)

type BillExecutor struct {
	workerChs [fillWorkerNum]chan *models.Bill
	name      string
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewBillExecutor() *BillExecutor {
	f := &BillExecutor{
		workerChs: [fillWorkerNum]chan *models.Bill{},
		name:      "bill excutor",
	}
	f.ctx, f.cancel = context.WithCancel(context.Background())
	// 初始化和fillWorkersNum一样数量的routine，每个routine负责一个chan
	for i := 0; i < fillWorkerNum; i++ {
		f.workerChs[i] = make(chan *models.Bill, 256)
		go func(idx int) {
			f.wg.Add(1)
			defer func() {
				for len(f.workerChs[idx]) > 0 {
					select {
					case <-f.workerChs[idx]:
					default:
					}
				}
				f.wg.Done()
			}()
			for {
				select {
				case <-f.ctx.Done():
					return
				case bill := <-f.workerChs[idx]:
					err := service.ExecuteBill(bill.UserId, bill.Currency)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}(i)
	}
	return f
}

func (s *BillExecutor) Start() {
	// go s.runMqListener()
	go s.runInspector()
}

func (s *BillExecutor) Stop() {
	s.cancel()
	log.Info("stopping ...", s.name)
	s.wg.Wait()
	log.Info("stopped  ", s.name)
}

func (s *BillExecutor) runMqListener() {
	gbeConfig := conf.GetConfig()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     gbeConfig.Redis.Addr,
		Password: gbeConfig.Redis.Password,
		DB:       1,
	})

	for {
		ret := redisClient.BRPop(time.Second*1000, models.TopicBill)
		if ret.Err() != nil {
			log.Error(ret.Err())
			continue
		}

		var bill models.Bill
		err := json.Unmarshal([]byte(ret.Val()[1]), &bill)
		if err != nil {
			panic(ret.Err())
		}

		// 按userId进行sharding
		s.workerChs[bill.UserId%fillWorkerNum] <- &bill
	}
}

func (s *BillExecutor) runInspector() {
	s.wg.Add(1)
	defer s.wg.Done()
	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(1 * time.Second):
			bills, err := service.GetUnsettledBills()
			if err != nil {
				log.Error(err)
				continue
			}

			for _, bill := range bills {
				s.workerChs[bill.UserId%fillWorkerNum] <- bill
			}
		}
	}
}
