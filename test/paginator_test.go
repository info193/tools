package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestPaginator(t *testing.T) {
	paginator := utils.NewPaginator(1, 300, 11111111)
	fmt.Println(paginator.Page(), paginator.Prepage(), paginator.TotalPage(), paginator.Count())
}
