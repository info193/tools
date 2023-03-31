package test

import (
	"fmt"
	"github.com/info193/tools/random"
	"testing"
)

func TestRandStr(t *testing.T) {
	fmt.Println(random.RandomStr(5))
	fmt.Println(random.RandomInt(5))
	fmt.Println(random.RandomInt64(5))
}
