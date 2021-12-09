package testt

import (
	"math/rand"
	"time"
)

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
