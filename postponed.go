package counter

import (
	"math"
	"sync/atomic"
	"time"
)

type PostponeReduced struct {
	c uint64
	d Decrementer
}

func NewPostponeReduced(d Decrementer) *PostponeReduced {
	c := &PostponeReduced{d: d}
	return c
}

func (c *PostponeReduced) Inc() *PostponeReduced {
	atomic.AddUint64(&c.c, 1)
	return c
}

func (c *PostponeReduced) Dec() *PostponeReduced {
	atomic.AddUint64(&c.c, math.MaxUint64)
	return c
}

func (c *PostponeReduced) DecAfter(d time.Duration) *PostponeReduced {
	dec := c.d
	if dec == nil {
		dec = nativeDec
	}
	c.d.Decrement(c, d)
	return c
}

func (c *PostponeReduced) Sum() uint64 {
	return atomic.LoadUint64(&c.c)
}
