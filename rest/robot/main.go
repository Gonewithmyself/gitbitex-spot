package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gitbitex/gitbitex-spot/service"
)

var done = make(chan struct{})
var wg sync.WaitGroup

func sinalNotify() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	close(done)
}

func main() {

	// placeOrder("ETH-USDT")
	pds, err := service.GetProducts()
	if nil != err {
		log.Fatal("get", err)
	}

	for i := range pds {
		if pds[i].Id == "EOS-USDT" {
			wg.Add(1)
			go loopPlaceOrder(pds[i].Id)
		}

		// break
	}

	wg.Wait()
	fmt.Println("exit normal")
}

func loopPlaceOrder(pair string) {
	defer wg.Done()
	t := time.NewTimer(time.Hour)
	for {
		t.Reset(time.Microsecond * time.Duration(RandInt(1500, 2000)))
		select {
		case <-done:
			return
		case <-t.C:
			for i := 0; i < rand.Intn(1)+1; i++ {
				placeOrder(pair)
			}
			//placeOrder(pair)
		}
	}
}
