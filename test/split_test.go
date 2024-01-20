package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestSplitStr64(t *testing.T) {
	fmt.Println(utils.SplitStr64([]int64{}))
	fmt.Println(utils.SplitStr([]int{1, 2, 3}))
}

