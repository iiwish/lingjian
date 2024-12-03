package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

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

// HashPassword 密码哈希
func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
