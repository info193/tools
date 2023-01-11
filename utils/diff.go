package utils

func DiffString(old, dst []string) []string {
	if len(old) == 0 || len(dst) == 0 {
		return []string{}
	}
	arr := make([]string, 0)
	for _, v := range old {
		exist := false
		for _, w := range dst {
			if v == w {
				exist = true
				break
			}
		}
		if !exist {
			arr = append(arr, v)
		}
	}
	return arr
}

func DiffInt(old, dst []int) []int {
	if len(old) == 0 || len(dst) == 0 {
		return []int{}
	}
	arr := make([]int, 0)
	for _, v := range old {
		exist := false
		for _, w := range dst {
			if v == w {
				exist = true
				break
			}
		}
		if !exist {
			arr = append(arr, v)
		}
	}
	return arr
}

func DiffInt32(old, dst []int32) []int32 {
	if len(old) == 0 || len(dst) == 0 {
		return []int32{}
	}
	arr := make([]int32, 0)
	for _, v := range old {
		exist := false
		for _, w := range dst {
			if v == w {
				exist = true
				break
			}
		}
		if !exist {
			arr = append(arr, v)
		}
	}
	return arr
}

func DiffInt64(old, dst []int64) []int64 {
	if len(old) == 0 || len(dst) == 0 {
		return []int64{}
	}
	arr := make([]int64, 0)
	for _, v := range old {
		exist := false
		for _, w := range dst {
			if v == w {
				exist = true
				break
			}
		}
		if !exist {
			arr = append(arr, v)
		}
	}
	return arr
}
