package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestUniqueString(t *testing.T) {
	maps := make(map[string]string, 0)
	//maps["name"] = "sss"
	//maps["demo"] = "ddd"
	fmt.Println(len(maps))

	str := []string{"A", "B", "C", "A", "E", "C"}
	old := []string{"D", "A", "E", "C"}

	fmt.Println(utils.UniqueSliceString(str, []string{}, old))

	intStr := []int{1, 2, 3, 4, 5}
	intOld := []int{6, 8, 2, 3, 4}
	for _, v := range intOld {
		intStr = utils.DeleteSliceInt(intStr, v)
	}
	fmt.Println(intStr)
}
