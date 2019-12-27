package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gitbitex/gitbitex-spot/conf"
	"github.com/gitbitex/gitbitex-spot/matching"
	"github.com/gitbitex/gitbitex-spot/models"
	"github.com/gitbitex/gitbitex-spot/service"
	"github.com/gitbitex/gitbitex-spot/worker"
)

var wg sync.WaitGroup
var workers []Worker

type Worker interface {
	Start()
	Stop()
}

func main() {
	wg.Add(1)

	go models.NewBinLogStream().Start()

	fill := worker.NewFillExecutor()
	bill := worker.NewBillExecutor()
	regWorker(fill, bill)
	products, err := service.GetProducts()
	if err != nil {
		panic(err)
	}

	gbeConfig := conf.GetConfig()
	for _, product := range products {
		fillmk := worker.NewFillMaker(matching.NewKafkaLogReader("fillMaker", product.Id, gbeConfig.Kafka.Brokers))
		billmk := worker.NewBillMaker(matching.NewKafkaLogReader("billMaker", product.Id, gbeConfig.Kafka.Brokers))
		regWorker(fillmk, billmk)
	}

	startAll()
	go signalHandler()
	wg.Wait()
}

func regWorker(s ...Worker) {
	workers = append(workers, s...)
}

func startAll() {
	for i := range workers {
		workers[i].Start()
	}
}
func stopAll() {
	for i := range workers {
		workers[i].Stop()
	}
}

func signalHandler() {
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	<-notifier
	stopAll()
	signal.Stop(notifier)
	wg.Done()
}
