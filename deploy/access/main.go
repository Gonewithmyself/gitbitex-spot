package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gitbitex/gitbitex-spot/rest"
	"github.com/gitbitex/gitbitex-spot/service"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)
	rest.StartServer()
	go signalHandler()
	wg.Wait()
}

func signalHandler() {
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	<-notifier

	rest.StopServer()
	signal.Stop(notifier)
	wg.Done()
}

var minutes = []int64{1, 3, 5, 15, 30, 60, 120, 240, 360, 720, 1440}

func init() { service.InitTbmap(minutes) }
