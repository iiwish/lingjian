package store

import (
	"context"
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
