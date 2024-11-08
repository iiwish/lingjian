package test

import (
	"fmt"
	"log"
	"sync"

	"github.com/iiwish/lingjian/pkg/utils"
)

// MockStore 用于测试的存储实现
type MockStore struct {
	tokens sync.Map // 用于存储token
}

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
	log.Printf("存储访问令牌 - UserID: %d, Token: %s", userId, token)
	// 解析token以获取完整的claims信息
	claims, err := utils.ParseToken(token, utils.AccessToken)
	if err != nil {
		return err
	}
	s.tokens.Store(token, claims)
	return nil
}

// StoreRefreshToken 存储刷新令牌
func (s *MockStore) StoreRefreshToken(userId uint, token string) error {
	log.Printf("存储刷新令牌 - UserID: %d, Token: %s", userId, token)
	// 解析token以获取完整的claims信息
	claims, err := utils.ParseToken(token, utils.RefreshToken)
	if err != nil {
		return err
	}
	s.tokens.Store(token, claims)
	return nil
}

// VerifyToken 验证令牌
func (s *MockStore) VerifyToken(token, tokenType string) (uint, error) {
	if token == "invalid_token" {
		return 0, fmt.Errorf("invalid token")
	}

	// 解析token获取用户ID和角色信息
	var claims *utils.CustomClaims
	var err error

	if tokenType == "access" {
		claims, err = utils.ParseToken(token, utils.AccessToken)
	} else {
		claims, err = utils.ParseToken(token, utils.RefreshToken)
	}

	if err != nil {
		log.Printf("解析令牌失败 - Token: %s, Error: %v", token, err)
		return 0, err
	}

	// 验证token是否在存储中
	storedValue, exists := s.tokens.Load(token)
	if !exists {
		log.Printf("令牌未找到 - Token: %s", token)
		return 0, fmt.Errorf("token not found in store")
	}

	storedClaims, ok := storedValue.(*utils.CustomClaims)
	if !ok {
		log.Printf("存储的claims类型错误 - Token: %s", token)
		return 0, fmt.Errorf("invalid claims type in store")
	}

	// 验证存储的用户ID和角色是否与token中的匹配
	if storedClaims.UserID != claims.UserID {
		log.Printf("用户ID不匹配 - Stored: %d, Token: %d", storedClaims.UserID, claims.UserID)
		return 0, fmt.Errorf("user id mismatch")
	}

	if storedClaims.RoleCode != claims.RoleCode {
		log.Printf("角色不匹配 - Stored: %s, Token: %s", storedClaims.RoleCode, claims.RoleCode)
		return 0, fmt.Errorf("role code mismatch")
	}

	log.Printf("令牌验证成功 - UserID: %d, RoleCode: %s, Token: %s, TokenType: %s",
		claims.UserID, claims.RoleCode, token, tokenType)
	return claims.UserID, nil
}

// RemoveUserTokens 移除用户的所有令牌
func (s *MockStore) RemoveUserTokens(userId uint) error {
	log.Printf("移除用户令牌 - UserID: %d", userId)
	// 遍历并删除该用户的所有token
	s.tokens.Range(func(key, value interface{}) bool {
		if claims, ok := value.(*utils.CustomClaims); ok && claims.UserID == userId {
			s.tokens.Delete(key)
		}
		return true
	})
	return nil
}

// StoreAuthCode 存储授权码
func (s *MockStore) StoreAuthCode(code, clientID, scope string, expiry int) error {
	return nil
}

// GetAuthCode 获取授权码信息
func (s *MockStore) GetAuthCode(code string) (clientID string, scope string, err error) {
	if code == "invalid_code" {
		return "", "", fmt.Errorf("invalid code")
	}
	return "test_client", "read", nil
}

// StoreOAuthToken 存储OAuth令牌
func (s *MockStore) StoreOAuthToken(accessToken, refreshToken, clientID, scope string) error {
	return nil
}

// GetRefreshToken 获取刷新令牌信息
func (s *MockStore) GetRefreshToken(refreshToken string) (clientID string, scope string, err error) {
	if refreshToken == "invalid_refresh_token" {
		return "", "", fmt.Errorf("invalid refresh token")
	}
	return "test_client", "read", nil
}

// UpdateOAuthAccessToken 更新访问令牌
func (s *MockStore) UpdateOAuthAccessToken(refreshToken, newAccessToken string) error {
	if refreshToken == "invalid_refresh_token" {
		return fmt.Errorf("invalid refresh token")
	}
	return nil
}
