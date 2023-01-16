package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestContainsSliceString(t *testing.T) {
	dst := []string{"A", "B", "C", "D"}
	fmt.Println(utils.ContainsSliceString(dst, "E"))
}
