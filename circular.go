package counter

import (
	"sync"
	"sync/atomic"
	"time"
)

// Counter is an alias of Circular. Made for backward compatibility.
type Counter = Circular

// NewCounter id an alias of NewCircular.
func NewCounter() *Counter {
	return NewCircular()
}

// Circular is a circular realtime counter.
// Please don't init counter directly by using `counter.Circular{}`. Use `counter.NewCounter()` instead.
type Circular struct {
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
	registry []*Circular
	// Channel to stop all counters.
	done chan struct{}
)

// NewCircular makes new counter and registry it.
func NewCircular() *Circular {
	c := &Circular{}
	// Take address of global now in milliseconds.
	c.now = &globNow
	// Registry new counter.
	mux.Lock()
	registry = append(registry, c)
	mux.Unlock()
	return c
}

// Inc increases counter.
func (c *Circular) Inc() {
	// Get current millisecond.
	now := atomic.LoadInt64(c.now)
	// Increase counter of current millisecond.
	atomic.AddUint32(&c.msec[now], 1)
}

// Sum returns current value of the counter.
func (c *Circular) Sum() uint32 {
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

// Reset counter for given millisecond value.
func (c *Circular) reset(idx int64) {
	atomic.StoreUint32(&c.msec[idx], 0)
}

// StopAll stops all counters.
func StopAll() {
	done <- struct{}{}
}

func init() {
	// Prepare done channel.
	done = make(chan struct{}, 1)
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
				tickMsec.Stop()
				return
			}
		}
	}()
}
