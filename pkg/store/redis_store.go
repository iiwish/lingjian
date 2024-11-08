package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/iiwish/lingjian/pkg/redis"
)

type RedisStore struct{}

var (
	// 验证码相关
	captchaKeyPrefix = "captcha:"
	captchaTTL       = 5 * time.Minute

	// token相关
	accessTokenKeyPrefix  = "token:access:"
	refreshTokenKeyPrefix = "token:refresh:"
	userTokensKeyPrefix   = "user:tokens:"
	accessTokenTTL        = 2 * time.Hour
	refreshTokenTTL       = 7 * 24 * time.Hour

	// OAuth2相关
	authCodeKeyPrefix    = "oauth:code:"
	oauthTokenKeyPrefix  = "oauth:token:"
	authCodeTTL          = 10 * time.Minute
	oauthAccessTokenTTL  = 2 * time.Hour
	oauthRefreshTokenTTL = 24 * time.Hour

	// 全局单例
	globalStore *RedisStore
)

// NewRedisStore 获取RedisStore实例
func NewRedisStore() *RedisStore {
	if globalStore == nil {
		globalStore = &RedisStore{}
	}
	return globalStore
}

// Set 实现验证码存储接口
func (s *RedisStore) Set(id string, value string) error {
	key := captchaKeyPrefix + id
	return redis.Set(context.Background(), key, value, int(captchaTTL.Seconds()))
}

// Get 实现验证码存储接口
func (s *RedisStore) Get(id string, clear bool) string {
	key := captchaKeyPrefix + id
	ctx := context.Background()
	val, err := redis.Get(ctx, key)
	if err != nil {
		return ""
	}
	if clear {
		redis.Del(ctx, key)
	}
	return val
}

// Verify 实现验证码存储接口
func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	v := s.Get(id, clear)
	return v == answer
}

// StoreAccessToken 存储访问令牌
func (s *RedisStore) StoreAccessToken(userId uint, token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", accessTokenKeyPrefix, token)
	userKey := fmt.Sprintf("%s%d", userTokensKeyPrefix, userId)

	// 存储token
	if err := redis.Set(ctx, key, userId, int(accessTokenTTL.Seconds())); err != nil {
		return err
	}

	// 将token添加到用户的token列表
	return redis.Set(ctx, userKey, token, int(accessTokenTTL.Seconds()))
}

// StoreRefreshToken 存储刷新令牌
func (s *RedisStore) StoreRefreshToken(userId uint, token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", refreshTokenKeyPrefix, token)
	userKey := fmt.Sprintf("%s%d", userTokensKeyPrefix, userId)

	// 存储token
	if err := redis.Set(ctx, key, userId, int(refreshTokenTTL.Seconds())); err != nil {
		return err
	}

	// 将token添加到用户的token列表
	return redis.Set(ctx, userKey, token, int(refreshTokenTTL.Seconds()))
}

// VerifyToken 验证令牌
func (s *RedisStore) VerifyToken(token, tokenType string) (uint, error) {
	ctx := context.Background()
	var key string

	switch tokenType {
	case "access":
		key = fmt.Sprintf("%s%s", accessTokenKeyPrefix, token)
	case "refresh":
		key = fmt.Sprintf("%s%s", refreshTokenKeyPrefix, token)
	default:
		return 0, fmt.Errorf("invalid token type")
	}

	val, err := redis.Get(ctx, key)
	if err != nil {
		return 0, err
	}

	var userId uint
	_, err = fmt.Sscanf(val, "%d", &userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

// RemoveUserTokens 移除用户的所有令牌
func (s *RedisStore) RemoveUserTokens(userId uint) error {
	ctx := context.Background()
	userKey := fmt.Sprintf("%s%d", userTokensKeyPrefix, userId)

	// 获取用户的token
	token, err := redis.Get(ctx, userKey)
	if err != nil {
		return nil // 如果没有找到token，视为成功
	}

	// 删除access token和refresh token
	accessKey := fmt.Sprintf("%s%s", accessTokenKeyPrefix, token)
	refreshKey := fmt.Sprintf("%s%s", refreshTokenKeyPrefix, token)

	if err := redis.Del(ctx, accessKey); err != nil {
		return err
	}
	if err := redis.Del(ctx, refreshKey); err != nil {
		return err
	}

	// 删除用户的token记录
	return redis.Del(ctx, userKey)
}

// StoreAuthCode 存储授权码
func (s *RedisStore) StoreAuthCode(code, clientID, scope string, expiry int) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", authCodeKeyPrefix, code)

	data := map[string]string{
		"client_id": clientID,
		"scope":     scope,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return redis.Set(ctx, key, string(jsonData), expiry)
}

// GetAuthCode 获取授权码信息
func (s *RedisStore) GetAuthCode(code string) (clientID string, scope string, err error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", authCodeKeyPrefix, code)

	jsonData, err := redis.Get(ctx, key)
	if err != nil {
		return "", "", fmt.Errorf("授权码不存在或已过期")
	}

	// 使用后立即删除授权码（一次性使用）
	redis.Del(ctx, key)

	var data map[string]string
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return "", "", err
	}

	return data["client_id"], data["scope"], nil
}

// StoreOAuthToken 存储OAuth令牌
func (s *RedisStore) StoreOAuthToken(accessToken, refreshToken, clientID, scope string) error {
	ctx := context.Background()

	// 存储访问令牌信息
	accessKey := fmt.Sprintf("%s%s", oauthTokenKeyPrefix, accessToken)
	accessData := map[string]string{
		"client_id":     clientID,
		"scope":         scope,
		"type":          "access",
		"refresh_token": refreshToken,
	}
	accessJsonData, err := json.Marshal(accessData)
	if err != nil {
		return err
	}

	// 存储刷新令牌信息
	refreshKey := fmt.Sprintf("%s%s", oauthTokenKeyPrefix, refreshToken)
	refreshData := map[string]string{
		"client_id":    clientID,
		"scope":        scope,
		"type":         "refresh",
		"access_token": accessToken,
	}
	refreshJsonData, err := json.Marshal(refreshData)
	if err != nil {
		return err
	}

	// 使用pipeline存储两个令牌
	if err := redis.Set(ctx, accessKey, string(accessJsonData), int(oauthAccessTokenTTL.Seconds())); err != nil {
		return err
	}
	return redis.Set(ctx, refreshKey, string(refreshJsonData), int(oauthRefreshTokenTTL.Seconds()))
}

// GetRefreshToken 获取刷新令牌信息
func (s *RedisStore) GetRefreshToken(refreshToken string) (clientID string, scope string, err error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", oauthTokenKeyPrefix, refreshToken)

	jsonData, err := redis.Get(ctx, key)
	if err != nil {
		return "", "", fmt.Errorf("刷新令牌不存在或已过期")
	}

	var data map[string]string
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return "", "", err
	}

	if data["type"] != "refresh" {
		return "", "", fmt.Errorf("无效的刷新令牌")
	}

	return data["client_id"], data["scope"], nil
}

// UpdateOAuthAccessToken 更新访问令牌
func (s *RedisStore) UpdateOAuthAccessToken(refreshToken, newAccessToken string) error {
	ctx := context.Background()
	refreshKey := fmt.Sprintf("%s%s", oauthTokenKeyPrefix, refreshToken)

	// 获取刷新令牌信息
	jsonData, err := redis.Get(ctx, refreshKey)
	if err != nil {
		return fmt.Errorf("刷新令牌不存在或已过期")
	}

	var data map[string]string
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return err
	}

	// 删除旧的访问令牌
	if oldAccessToken := data["access_token"]; oldAccessToken != "" {
		oldKey := fmt.Sprintf("%s%s", oauthTokenKeyPrefix, oldAccessToken)
		redis.Del(ctx, oldKey)
	}

	// 创建新的访问令牌
	newAccessKey := fmt.Sprintf("%s%s", oauthTokenKeyPrefix, newAccessToken)
	accessData := map[string]string{
		"client_id":     data["client_id"],
		"scope":         data["scope"],
		"type":          "access",
		"refresh_token": refreshToken,
	}
	accessJsonData, err := json.Marshal(accessData)
	if err != nil {
		return err
	}

	// 更新刷新令牌中的访问令牌信息
	data["access_token"] = newAccessToken
	refreshJsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 存储新的访问令牌
	if err := redis.Set(ctx, newAccessKey, string(accessJsonData), int(oauthAccessTokenTTL.Seconds())); err != nil {
		return err
	}

	// 更新刷新令牌信息
	return redis.Set(ctx, refreshKey, string(refreshJsonData), int(oauthRefreshTokenTTL.Seconds()))
}
