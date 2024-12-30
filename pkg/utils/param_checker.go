package utils

import (
	"regexp"
)

// IsUsername 校验用户名格式
func IsUsername(username string) bool {
	// 用户名正则：字母开头，允许字母数字下划线，4-16位
	var re = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{3,15}$`)
	return re.MatchString(username)
}

// IsPassword 校验密码格式
func IsPassword(password string) bool {
	// 密码正则：允许字母数字.-_，6-20位
	var re = regexp.MustCompile(`^[a-zA-Z0-9.-_]{6,20}$`)
	return re.MatchString(password)
}

// IsCode 校验验证码格式
func IsCode(code string) bool {
	// 验证码正则：字母开头，允许字母数字下划线，2-50位
	var re = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{1,49}$`)
	return re.MatchString(code)
}

// IsValidIdentifier 校验标识符格式
func IsValidIdentifier(identifier string) bool {
	// 标识符正则：字母开头，允许字母数字下划线，2-100位
	var re = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{1,99}$`)
	return re.MatchString(identifier)
}
