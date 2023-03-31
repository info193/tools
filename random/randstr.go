package random

import (
	"math/rand"
	"time"
)

func RandomStr(n int) string {
	letters := []byte("aeFfghi7jklnXopqr123stuvwQxy6zABCD8EGbcd5HIJKL9MNOPRmST0UVWYZ4")
	rand.Seed(time.Now().UnixNano())
	str := make([]byte, n)
	for i := 0; i < n; i++ {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}
func RandomInt(n int) string {
	letters := []byte("5062847193")
	rand.Seed(time.Now().UnixNano())
	str := make([]byte, n)
	for i := 0; i < n; i++ {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}
func RandomInt64(n int) string {
	letters := []byte("5129470836")
	rand.Seed(time.Now().UnixNano())
	str := make([]byte, n)
	for i := 0; i < n; i++ {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}
