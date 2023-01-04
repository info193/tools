package test

import (
	"fmt"
	"github.com/info193/tools/random"
	"testing"
)

func TestGenerateNickname(t *testing.T) {
	fmt.Println(random.GenerateNickname(0)) // 随机男或女
	fmt.Println(random.GenerateNickname(1)) // 男
	fmt.Println(random.GenerateNickname(2)) // 女
}
