package counter

import "sync/atomic"

type Counter struct {
	now  *int64
	msec [1000]uint32
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
		sum += (c.val(i) + c.val(i+1)) + (c.val(i+2) + c.val(i+3)) + (c.val(i+4) + c.val(i+5)) + (c.val(i+6) + c.val(i+7)) + (c.val(i+8) + c.val(i+9))
	}
	return sum
}

func (c *Counter) StopAll() {
	done <- true
}

func (c *Counter) reset(idx int64) {
	c.msec[idx] = 0
}

func (c *Counter) val(idx int) uint32 {
	prev := idx - 1
	if prev < 0 {
		prev = 999
	}
	if *c.now == int64(prev) || *c.now == int64(idx) {
		return atomic.LoadUint32(&c.msec[idx])
	}
	return c.msec[idx]
}
