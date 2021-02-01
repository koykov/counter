package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/koykov/counter"
)

var (
	wg sync.WaitGroup
	tick,
	tickDone *time.Ticker
	done uint32
)

func main() {
	c := counter.NewCounter()
	tick = time.NewTicker(time.Millisecond * 100)
	tickDone = time.NewTicker(time.Second * 10)
	go func() {
		for {
			select {
			case <-tick.C:
				fmt.Println("counter", c.Sum())
			case <-tickDone.C:
				atomic.StoreUint32(&done, 1)
				return
			}
		}
	}()

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			for {
				c.Inc()
				if atomic.LoadUint32(&done) == 1 {
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	c.StopAll()
}
