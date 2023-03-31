package utils

import (
	"fmt"
	"strings"
)

func Reverse(s string) string {
	strSlice := strings.Split(s, "")
	var str string
	for i := len(strSlice); i > 0; i-- {
		str = str + fmt.Sprintf("%s", strSlice[i-1])
	}
	return str
}
func ReverseRune(s string) string {
	strSlice := []rune(s)
	len := len(strSlice)
	var str string
	for i := len; i > 0; i-- {
		str = str + fmt.Sprintf("%s", string(strSlice[i-1]))
	}
	return str
}
