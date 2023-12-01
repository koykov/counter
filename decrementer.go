package counter

import "time"

type Decrementer interface {
	Decrement(c *PostponeReduced, d time.Duration)
}

type NativeDecrement struct{}

func (d NativeDecrement) Decrement(c *PostponeReduced, d_ time.Duration) {
	time.AfterFunc(d_, func() {
		c.Dec()
	})
}

var nativeDec = NativeDecrement{}
