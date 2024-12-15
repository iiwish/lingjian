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

// convertBytesToString 递归地将 map 中的 []byte 转换为字符串
func ConvertBytesToString(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[key] = ConvertBytesToString(value)
		}
	case []interface{}:
		for i, value := range v {
			v[i] = ConvertBytesToString(value)
		}
	case []map[string]interface{}:
		for i, value := range v {
			v[i] = ConvertBytesToString(value).(map[string]interface{})
		}
	case []byte:
		return string(v)
	}
	return data
}

// Contains 检查字符串切片中是否包含指定字符串
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// HashPassword 密码哈希
func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
