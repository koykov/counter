package counter

import (
	"sync/atomic"
	"time"
)

func init() {
	done = make(chan bool, 1)
	tickMsec := time.NewTicker(time.Millisecond)
	go func() {
		for {
			select {
			case <-tickMsec.C:
				now := (time.Now().UnixNano() / 1e6) % 1000
				atomic.StoreInt64(&globNow, now)
				mux.RLock()
				for i := 0; i < len(registry); i++ {
					registry[i].reset(now)
				}
				mux.RUnlock()
			case <-done:
				return
			}
		}
	}()
}
