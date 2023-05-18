package utils

import (
	"math/rand"
	"strings"
)

func RandomString(n int) string {
	const characters = "abcdefghijklmnopqrstuvwxyz"
	var s strings.Builder

	for i := 0; i < n; i++ {
		randomInt := rand.Intn(len(characters))
		s.WriteByte(characters[randomInt])
	}

	return s.String()
}

func RandomMoney() int64 {
	return rand.Int63n(1001)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "INR"}
	return currencies[rand.Intn(len(currencies))]
}
