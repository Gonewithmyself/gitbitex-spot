package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gitbitex/gitbitex-spot/conf"
	"github.com/gitbitex/gitbitex-spot/matching"
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
	products, err := service.GetProducts()
	if err != nil {
		panic(err)
	}

	gbeConfig := conf.GetConfig()
	for _, product := range products {
		regWorker(
			worker.NewTradeMaker(matching.NewKafkaLogReader("tradeMaker", product.Id, gbeConfig.Kafka.Brokers)),
			worker.NewTickMaker(product.Id, matching.NewKafkaLogReader("tickMaker", product.Id, gbeConfig.Kafka.Brokers)),
			worker.NewDepthMaker(product.Id, matching.NewKafkaLogReader("depthMaker", product.Id, gbeConfig.Kafka.Brokers)),
		)
	}

	startAll()
	wg.Add(1)
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
