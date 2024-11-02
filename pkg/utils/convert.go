package utils

import "strconv"

// ParseUint 将字符串转换为uint
func ParseUint(s string) uint {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(n)
}

// ParseInt 将字符串转换为int
func ParseInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
