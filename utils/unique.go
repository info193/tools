package utils

func UniqueString(sli ...[]string) []string {
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

func UniqueInt(sli ...[]int) []int {
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

func UniqueInt32(sli ...[]int32) []int32 {
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

func UniqueInt64(sli ...[]int64) []int64 {
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
