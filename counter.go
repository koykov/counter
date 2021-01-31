package counter

import (
	"sync"
	"sync/atomic"
)

type Counter struct {
	now  *int64
	msec [1000]uint32
}

var (
	globNow  int64
	mux      sync.RWMutex
	registry []*Counter
	done     chan bool
)

func NewCounter() *Counter {
	c := &Counter{}
	c.now = &globNow
	mux.Lock()
	registry = append(registry, c)
	mux.Unlock()
	return c
}

func (c *Counter) Inc() {
	now := atomic.LoadInt64(c.now)
	atomic.AddUint32(&c.msec[now], 1)
}

func (c *Counter) Sum() uint32 {
	var sum uint32
	for i := 0; i < 1000; i += 10 {
		sum += (atomic.LoadUint32(&c.msec[i]) + atomic.LoadUint32(&c.msec[i+1])) +
			(atomic.LoadUint32(&c.msec[i+2]) + atomic.LoadUint32(&c.msec[i+3])) +
			(atomic.LoadUint32(&c.msec[i+4]) + atomic.LoadUint32(&c.msec[i+5])) +
			(atomic.LoadUint32(&c.msec[i+6]) + atomic.LoadUint32(&c.msec[i+7])) +
			(atomic.LoadUint32(&c.msec[i+8]) + atomic.LoadUint32(&c.msec[i+9]))
	}
	return sum
}

func (c *Counter) StopAll() {
	done <- true
}

func (c *Counter) reset(idx int64) {
	atomic.StoreUint32(&c.msec[idx], 0)
}
