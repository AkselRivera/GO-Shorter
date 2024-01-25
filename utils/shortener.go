package utils

import (
	"math/rand"
	"time"
)

const Dictionary = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetHash() string {
	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed))

	var hash string

	for i := 0; i < 6; i++ {
		hash += string(Dictionary[rand.Intn(len(Dictionary))])
	}

	return hash
}
