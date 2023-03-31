package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestMd5(t *testing.T) {
	fmt.Println(utils.Md5(utils.Md5("123456") + "aw92"))
	//fmt.Println()
	//fmt.Println(utils.ReverseRune("你12345..6"))
}
