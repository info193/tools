package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestDiffString(t *testing.T) {
	str := []string{"A", "B", "C"}
	old := []string{"D", "A"}

	fmt.Println(utils.DiffString(old, str))
}
