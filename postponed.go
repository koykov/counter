package counter

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type PostponeReduced struct {
	c uint64
	d Decrementer

	once sync.Once
	id   uint64
}

var prID uint64

func NewPostponeReduced(d Decrementer) *PostponeReduced {
	c := &PostponeReduced{d: d}
	_ = c.ID()
	return c
}

func (c *PostponeReduced) init() {
	c.id = atomic.AddUint64(&prID, 1)
}

func (c *PostponeReduced) ID() uint64 {
	c.once.Do(c.init)
	return c.id
}

func (c *PostponeReduced) Inc() *PostponeReduced {
	c.once.Do(c.init)
	atomic.AddUint64(&c.c, 1)
	return c
}

func (c *PostponeReduced) Dec() *PostponeReduced {
	c.once.Do(c.init)
	atomic.AddUint64(&c.c, math.MaxUint64)
	return c
}

func (c *PostponeReduced) DecAfter(d time.Duration) *PostponeReduced {
	c.once.Do(c.init)
	dec := c.d
	if dec == nil {
		dec = nativeDec
	}
	c.d.Decrement(c, d)
	return c
}

func (c *PostponeReduced) Sum() uint64 {
	c.once.Do(c.init)
	return atomic.LoadUint64(&c.c)
}
