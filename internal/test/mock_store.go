package test

import "fmt"

// MockStore 用于测试的存储实现
type MockStore struct{}

// Set 存储验证码
func (s *MockStore) Set(id string, value string) error {
	return nil
}

// Get 获取验证码
func (s *MockStore) Get(id string, clear bool) string {
	return "1234" // 在测试中总是返回1234
}

// Verify 验证验证码
func (s *MockStore) Verify(id, answer string, clear bool) bool {
	return answer == "1234" // 在测试中验证码总是1234
}

// StoreAccessToken 存储访问令牌
func (s *MockStore) StoreAccessToken(userId uint, token string) error {
	return nil
}

// StoreRefreshToken 存储刷新令牌
func (s *MockStore) StoreRefreshToken(userId uint, token string) error {
	return nil
}

// VerifyToken 验证令牌
func (s *MockStore) VerifyToken(token, tokenType string) (uint, error) {
	if token == "invalid_refresh_token" {
		return 0, fmt.Errorf("invalid token")
	}
	return 1, nil // 在测试中总是返回用户ID 1
}

// RemoveUserTokens 移除用户的所有令牌
func (s *MockStore) RemoveUserTokens(userId uint) error {
	return nil
}
