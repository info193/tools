package test

import (
	"fmt"
	"github.com/info193/tools/random"
	"testing"
)

func TestRandStr(t *testing.T) {
	fmt.Println(random.RandomStr(6, 19999999999))
	fmt.Println(random.RandomInt(6, 19999999999))
	fmt.Println(random.RandomInt64(6, 19999999999))
	fmt.Println(random.RandomStrSeed("?.,<>/:1"))
	fmt.Println(random.RandomIntSeed("19999999999"))
}
