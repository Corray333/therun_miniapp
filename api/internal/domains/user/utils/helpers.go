package utils

import (
	"time"

	"golang.org/x/exp/rand"
)

func GenerateRefreshToken() string {
	randomizer := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	symbols := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := ""
	for i := 0; i < 8; i++ {
		b += string(symbols[randomizer.Intn(len(symbols))])
	}
	return b
}
