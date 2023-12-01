package counter

import (
	"context"
	"testing"
	"time"
)

func TestPostponeReduced(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		c := NewPostponeReduced(NativeDecrement{})
		c.Inc().DecAfter(time.Millisecond * 100)
		if c.Sum() != 1 {
			t.FailNow()
		}
		time.Sleep(time.Millisecond * 150)
		if c.Sum() != 0 {
			t.FailNow()
		}
	})
	t.Run("parallel", func(t *testing.T) {
		c := NewPostponeReduced(NativeDecrement{})
		ctx, cancel := context.WithCancel(context.Background())
		for i := 0; i < 10; i++ {
			go func(ctx context.Context) {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						c.Inc().DecAfter(time.Second)
						time.Sleep(time.Microsecond * 500)
					}
				}
			}(ctx)
		}
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Second):
					t.Log(c.Sum())
				}
			}
		}(ctx)
		time.Sleep(time.Second * 10)
		cancel()
	})
}
