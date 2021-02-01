package counter

import (
	"sync"
	"sync/atomic"
	"time"
)

// Realtime counter.
// Please note, you shouldn't init counter directly by using `counter.Counter{}`. Use `counter.NewCounter()` instead.
type Counter struct {
	// Pointer to global now variable.
	now *int64
	// Array of counters for each millisecond.
	msec [1000]uint32
}

var (
	// Global now variable.
	// Contains current millisecond value and updates every millisecond in a cycle.
	globNow int64
	// Counters registry and protecting mutex.
	mux      sync.RWMutex
	registry []*Counter
	// Channel to stop all counters.
	done chan bool
)

// Create new counter and registry it.
func NewCounter() *Counter {
	c := &Counter{}
	// Take address of global now in milliseconds.
	c.now = &globNow
	// Registry new counter.
	mux.Lock()
	registry = append(registry, c)
	mux.Unlock()
	return c
}

// Increase counter.
func (c *Counter) Inc() {
	// Get current millisecond.
	now := atomic.LoadInt64(c.now)
	// Increase counter of current millisecond.
	atomic.AddUint32(&c.msec[now], 1)
}

// Get current value of the counter.
func (c *Counter) Sum() uint32 {
	var sum uint32
	// Roll up the loop with chunks of size 10.
	for i := 0; i < 1000; i += 10 {
		// Use brackets to break the data associativity.
		sum += (atomic.LoadUint32(&c.msec[i]) + atomic.LoadUint32(&c.msec[i+1])) +
			(atomic.LoadUint32(&c.msec[i+2]) + atomic.LoadUint32(&c.msec[i+3])) +
			(atomic.LoadUint32(&c.msec[i+4]) + atomic.LoadUint32(&c.msec[i+5])) +
			(atomic.LoadUint32(&c.msec[i+6]) + atomic.LoadUint32(&c.msec[i+7])) +
			(atomic.LoadUint32(&c.msec[i+8]) + atomic.LoadUint32(&c.msec[i+9]))
	}
	return sum
}

// Stop all counters.
func (c *Counter) StopAll() {
	done <- true
}

// Reset counter for given millisecond value.
func (c *Counter) reset(idx int64) {
	atomic.StoreUint32(&c.msec[idx], 0)
}

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
