package utils

import (
	"math/rand"
	"time"
)

func GenerateCode() int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	min := 100000
	max := 999999
	return rng.Intn(max-min+1) + min
}
