package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestIntersectString(t *testing.T) {
	str := []string{"A", "B", "C"}
	old := []string{"B", "D", "A"}

	fmt.Println(utils.IntersectString(old, str))
}
