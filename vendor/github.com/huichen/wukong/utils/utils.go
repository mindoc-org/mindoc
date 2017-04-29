package utils

func AbsInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
