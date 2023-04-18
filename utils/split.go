package utils

import (
	"strconv"
	"strings"
)

func SplitInt(str string, sep string) []int {
	strs := make([]int, 0)
	if str == "" {
		return strs
	}
	splits := strings.Split(str, sep)
	for _, split := range splits {
		id, _ := strconv.Atoi(split)
		strs = append(strs, id)
	}
	return strs
}

func SplitInt64(str string, sep string) []int64 {
	strs := make([]int64, 0)
	if str == "" {
		return strs
	}
	splits := strings.Split(str, sep)
	for _, split := range splits {
		id, _ := strconv.ParseInt(split, 10, 64)
		strs = append(strs, id)
	}
	return strs
}
