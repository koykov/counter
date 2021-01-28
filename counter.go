package counter

import "sync/atomic"

type Counter struct {
	now   *int64
	msec  [1000]uint32
	total uint32
}

var (
	globNow  int64
	registry []*Counter
	done     chan bool
)

func NewCounter() *Counter {
	c := &Counter{}
	c.now = &globNow
	registry = append(registry, c)
	return c
}

func (c *Counter) Inc() {
	now := *c.now
	atomic.AddUint32(&c.msec[now], 1)
}

func (c *Counter) Sum() uint32 {
	var sum uint32
	for i := 0; i < 1000; i += 10 {
		sum += (c.msec[i] + c.msec[i+1]) + (c.msec[i+2] + c.msec[i+3]) + (c.msec[i+4] + c.msec[i+5]) + (c.msec[i+6] + c.msec[i+7]) + (c.msec[i+8] + c.msec[i+9])
	}
	atomic.StoreUint32(&c.total, sum)
	return sum
}

func (c *Counter) reset(idx int64) {
	c.msec[idx] = 0
}

func (c *Counter) StopAll() {
	done <- true
}
