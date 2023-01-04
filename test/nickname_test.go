package test

import (
	"fmt"
	"testing"
	"tools/random"
)

func TestGenerateNickname(t *testing.T) {
	fmt.Println(random.GenerateNickname(0)) // 随机男或女
	fmt.Println(random.GenerateNickname(1)) // 男
	fmt.Println(random.GenerateNickname(2)) // 女
}
