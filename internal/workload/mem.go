package workload

import "time"

func MemAlloc(mb int, hold time.Duration) int {
	if mb <= 0 {
		return 0
	}
	size := mb * 1024 * 1024
	b := make([]byte, size)

	for i := 0; i < len(b); i += 4096 {
		b[i] = byte(i)
	}

	if hold > 0 {
		time.Sleep(hold)
	}

	return len(b)
}
