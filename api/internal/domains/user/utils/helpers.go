package utils

import "golang.org/x/exp/rand"

func GenerateRefreshToken() string {
	symbols := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := ""
	for i := 0; i < 8; i++ {
		b += string(symbols[rand.Intn(len(symbols))])
	}
	return b
}
