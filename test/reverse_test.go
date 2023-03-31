package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestReverse(t *testing.T) {
	fmt.Println(utils.Reverse("1234a556."))
	//fmt.Println()
	fmt.Println(utils.ReverseRune("你12345..6"))
}
