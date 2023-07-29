package random

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func RandomStrSeed(str string) int64 {
	var seed int64
	if str == "" {
		return seed
	}
	byteStr := []byte(str)

	for _, value := range byteStr {
		byteCode, _ := strconv.ParseInt(fmt.Sprintf("%d", value), 10, 64)
		seed += byteCode
	}
	return seed
}
func RandomIntSeed(str string) int64 {
	if str == "" {
		return 0
	}
	seed, _ := strconv.ParseInt(str, 10, 64)
	return seed
}

func RandomStr(n int, seed int64) string {
	letters := []byte("aeFfghi7jklnX2opqr13stuvwQxy6zABCD8EGbcd5HIJKL9MNOPRmST0UVWYZ4")
	rand.Seed(time.Now().UnixNano() + seed)
	str := make([]byte, n)
	for i := 0; i < n; i++ {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}
func RandomInt(n int, seed int64) string {
	letters := []byte("5062847193")
	rand.Seed(time.Now().UnixNano() + seed)
	str := make([]byte, n)
	for i := 0; i < n; i++ {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}
func RandomInt64(n int, seed int64) string {
	letters := []byte("5129470836")
	rand.Seed(time.Now().UnixNano() + seed)
	str := make([]byte, n)
	for i := 0; i < n; i++ {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}
