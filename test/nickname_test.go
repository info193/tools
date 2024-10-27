package test

import (
	"fmt"
	"github.com/info193/tools/random"
	"testing"
)

func TestGenerateNickname(t *testing.T) {
	fmt.Println(random.BatchGenerate(10))
	fmt.Println(random.BatchShortGenerate(10))
	//fmt.Println(random.BatchShortGenerate(3))
	//fmt.Println(random.BatchGenerate(3))
	//fmt.Println(random.Generate())
	//fmt.Println(random.GenerateShort())
}
