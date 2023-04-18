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

func SplitStr64(str []int64) string {
	strs := make([]string, 0)
	for _, split := range str {
		id := strconv.FormatInt(split, 10)
		strs = append(strs, id)
	}
	return strings.Join(strs, ",")
}

func SplitStr(str []int) string {
	strs := make([]string, 0)
	for _, split := range str {
		id := strconv.Itoa(split)
		strs = append(strs, id)
	}
	return strings.Join(strs, ",")
}
