package testpg

import (
	"math/rand"
	"time"
)

const (
	minLength = 512
	maxLength = 2048
)

func randomString() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	strLen := rand.Intn(maxLength-minLength+1) + minLength

	result := make([]byte, strLen)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
