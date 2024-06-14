package utils

import (
	"fmt"
	"strings"
)

func UniqueSliceString(sli ...[]string) []string {
	dst := make([]string, 0)
	length := len(sli)
	for i := 0; i < length; i++ {
		dst = append(dst, sli[i]...)
	}

	arr := make([]string, 0)
	tmpM := make(map[string]int)
	for _, v := range dst {
		if _, ok := tmpM[v]; ok {
			continue
		}
		tmpM[v] = 1
		arr = append(arr, v)
	}
	return arr
}

func UniqueSliceInt(sli ...[]int) []int {
	dst := make([]int, 0)
	length := len(sli)
	for i := 0; i < length; i++ {
		dst = append(dst, sli[i]...)
	}

	arr := make([]int, 0)
	tmpM := make(map[int]int)
	for _, v := range dst {
		if _, ok := tmpM[v]; ok {
			continue
		}
		tmpM[v] = 1
		arr = append(arr, v)
	}
	return arr
}

func UniqueSliceInt32(sli ...[]int32) []int32 {
	dst := make([]int32, 0)
	length := len(sli)
	for i := 0; i < length; i++ {
		dst = append(dst, sli[i]...)
	}

	arr := make([]int32, 0)
	tmpM := make(map[int32]int)
	for _, v := range dst {
		if _, ok := tmpM[v]; ok {
			continue
		}
		tmpM[v] = 1
		arr = append(arr, v)
	}
	return arr
}

func UniqueSliceInt64(sli ...[]int64) []int64 {
	dst := make([]int64, 0)
	length := len(sli)
	for i := 0; i < length; i++ {
		dst = append(dst, sli[i]...)
	}

	arr := make([]int64, 0)
	tmpM := make(map[int64]int)
	for _, v := range dst {
		if _, ok := tmpM[v]; ok {
			continue
		}
		tmpM[v] = 1
		arr = append(arr, v)
	}
	return arr
}
func DeleteSliceInt(a []int, elem int) []int {
	tgt := a[:0]
	for _, v := range a {
		if v != elem {
			tgt = append(tgt, v)
		}
	}
	return tgt
}
func DeleteSliceInt32(a []int32, elem int32) []int32 {
	tgt := a[:0]
	for _, v := range a {
		if v != elem {
			tgt = append(tgt, v)
		}
	}
	return tgt
}
func DeleteSliceInt64(a []int64, elem int64) []int64 {
	tgt := a[:0]
	for _, v := range a {
		if v != elem {
			tgt = append(tgt, v)
		}
	}
	return tgt
}
func DeleteSliceString(a []string, elem string) []string {
	tgt := a[:0]
	for _, v := range a {
		if v != elem {
			tgt = append(tgt, v)
		}
	}
	return tgt
}

func JoinStringInt64(a []int64, str string) string {
	stringSlice := make([]string, len(a))
	for i, num := range a {
		stringSlice[i] = fmt.Sprintf("%d", num)
	}
	return strings.Join(stringSlice, str)
}

func JoinStringInt(a []int, str string) string {
	stringSlice := make([]string, len(a))
	for i, num := range a {
		stringSlice[i] = fmt.Sprintf("%d", num)
	}
	return strings.Join(stringSlice, str)
}
