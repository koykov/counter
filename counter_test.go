package counter

import (
	"sync"
	"testing"
)

func BenchmarkCounter(b *testing.B) {
	c := NewCounter()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var wg sync.WaitGroup
		for pb.Next() {
			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					c.Inc()
					wg.Done()
				}()
			}
			wg.Wait()
		}
	})
	b.Log(c.Sum())
}
