package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestUniqueString(t *testing.T) {
	str := []string{"A", "B", "C", "A", "E", "C"}
	old := []string{"D", "A", "E", "C"}

	fmt.Println(utils.UniqueString(str, []string{}, old))
}
