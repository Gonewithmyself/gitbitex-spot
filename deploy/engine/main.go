package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gitbitex/gitbitex-spot/matching"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)
	matching.StartEngine()
	go signalHandler()
	wg.Wait()
}

func signalHandler() {
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	<-notifier
	matching.StopEngine()
	wg.Done()
}
