package utils

import (
	"math/rand"
	"time"
)

func RandStringBytes(n int) string {
	letterBytes := "1234567890"

	rands := rand.New(rand.NewSource(time.Now().Unix()))

	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[rands.Intn(len(letterBytes))]
	}

	return string(b)
}
