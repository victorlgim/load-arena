package workload

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

func CPU(n int, algo string) string {
	sum := sha256.Sum256([]byte("seed"))
	for i := 0; i < n; i++ {
		b := append(sum[:], []byte(strconv.Itoa(i))...)
		sum = sha256.Sum256(b)
	}
	return hex.EncodeToString(sum[:])
}
