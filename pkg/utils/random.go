package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}
