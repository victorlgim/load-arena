package workload

import (
	"math/rand"
	"os"
	"time"
)

type ChaosResult struct {
	Failed  bool
	Status  int
	Message string
}

func Chaos(rate float64, mode string) ChaosResult {
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}

	rand.Seed(time.Now().UnixNano())
	fail := rand.Float64() < rate

	if !fail {
		return ChaosResult{Failed: false, Status: 200, Message: "ok (no chaos)"}
	}

	switch mode {
	case "sleep":
		time.Sleep(2 * time.Second)
		return ChaosResult{Failed: true, Status: 504, Message: "chaos sleep timeout simulation"}
	case "exit":
		os.Exit(1)
		return ChaosResult{Failed: true, Status: 500, Message: "unreachable"}
	default: 
		return ChaosResult{Failed: true, Status: 500, Message: "chaos induced error"}
	}
}
