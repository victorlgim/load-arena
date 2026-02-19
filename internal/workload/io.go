package workload

import "time"

func IO(d time.Duration) time.Duration {
	if d <= 0 {
		return 0
	}
	time.Sleep(d)
	return d
}
