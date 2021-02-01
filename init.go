package counter

import (
	"sync/atomic"
	"time"
)

func init() {
	// Prepare done channel.
	done = make(chan bool, 1)
	// Init ticker for milliseconds.
	tickMsec := time.NewTicker(time.Millisecond)
	go func() {
		for {
			select {
			case <-tickMsec.C:
				// New millisecond came.
				// Update global now with current millisecond.
				now := (time.Now().UnixNano() / 1e6) % 1000
				atomic.StoreInt64(&globNow, now)

				mux.RLock()
				for i := 0; i < len(registry); i++ {
					// Reset counter of millisecond for each counter.
					registry[i].reset(now)
				}
				mux.RUnlock()
			case <-done:
				// Done signal caught, exiting.
				return
			}
		}
	}()
}
