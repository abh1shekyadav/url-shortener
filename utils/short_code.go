package utils

import (
	"log"
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var src = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateShortCode() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[src.Intn(len(letters))]
	}
	code := string(b)
	log.Println("GenerateShortCode output:", code)
	return code
}
