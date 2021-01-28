package counter

import "time"

func init() {
	done = make(chan bool, 1)
	tickMsec := time.NewTicker(time.Millisecond)
	go func() {
		for {
			select {
			case <-tickMsec.C:
				globNow = (time.Now().UnixNano() / 1e6) % 1000
				for i := 0; i < len(registry); i++ {
					registry[i].reset(globNow)
				}
			case <-done:
				return
			}
		}
	}()
}
