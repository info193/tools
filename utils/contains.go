package utils

func ContainsSliceString(dst []string, str string) bool {
	if len(dst) == 0 {
		return false
	}
	for _, value := range dst {
		if value == str {
			return true
		}
	}
	return false
}
func ContainsSliceInt(dst []int, str int) bool {
	if len(dst) == 0 {
		return false
	}
	for _, value := range dst {
		if value == str {
			return true
		}
	}
	return false
}

func ContainsSliceInt32(dst []int32, str int32) bool {
	if len(dst) == 0 {
		return false
	}
	for _, value := range dst {
		if value == str {
			return true
		}
	}
	return false
}

func ContainsSliceInt64(dst []int64, str int64) bool {
	if len(dst) == 0 {
		return false
	}
	for _, value := range dst {
		if value == str {
			return true
		}
	}
	return false
}
