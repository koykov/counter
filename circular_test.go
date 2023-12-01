package counter

import (
	"context"
	"testing"
	"time"
)

func TestCircular(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		c := NewCircular()
		c.Inc()
		if c.Sum() != 1 {
			t.FailNow()
		}
		time.Sleep(time.Millisecond * 1500)
		if c.Sum() != 0 {
			t.FailNow()
		}
	})
	t.Run("parallel", func(t *testing.T) {
		c := NewCircular()
		ctx, cancel := context.WithCancel(context.Background())
		for i := 0; i < 10; i++ {
			go func(ctx context.Context) {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						c.Inc()
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
