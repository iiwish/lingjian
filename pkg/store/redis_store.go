package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

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
)

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
	}
}

// 验证码相关方法
func (s *RedisStore) Set(id string, value string) error {
	key := captchaKeyPrefix + id
	return s.client.Set(context.Background(), key, value, captchaTTL).Err()
}

func (s *RedisStore) Get(id string, clear bool) string {
	key := captchaKeyPrefix + id
	ctx := context.Background()
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	if clear {
		s.client.Del(ctx, key)
	}
	return val
}

func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	v := s.Get(id, clear)
	return v == answer
}

// token相关方法
func (s *RedisStore) StoreAccessToken(userId uint, token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", accessTokenKeyPrefix, token)
	userKey := fmt.Sprintf("%s%d", userTokensKeyPrefix, userId)

	// 存储token
	if err := s.client.Set(ctx, key, userId, accessTokenTTL).Err(); err != nil {
		return err
	}

	// 将token添加到用户的token列表
	return s.client.SAdd(ctx, userKey, token).Err()
}

func (s *RedisStore) StoreRefreshToken(userId uint, token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", refreshTokenKeyPrefix, token)
	userKey := fmt.Sprintf("%s%d", userTokensKeyPrefix, userId)

	// 存储token
	if err := s.client.Set(ctx, key, userId, refreshTokenTTL).Err(); err != nil {
		return err
	}

	// 将token添加到用户的token列表
	return s.client.SAdd(ctx, userKey, token).Err()
}

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

	userId, err := s.client.Get(ctx, key).Uint64()
	if err != nil {
		return 0, err
	}

	return uint(userId), nil
}

func (s *RedisStore) RemoveUserTokens(userId uint) error {
	ctx := context.Background()
	userKey := fmt.Sprintf("%s%d", userTokensKeyPrefix, userId)

	// 获取用户的所有token
	tokens, err := s.client.SMembers(ctx, userKey).Result()
	if err != nil {
		return err
	}

	// 删除所有token
	for _, token := range tokens {
		accessKey := fmt.Sprintf("%s%s", accessTokenKeyPrefix, token)
		refreshKey := fmt.Sprintf("%s%s", refreshTokenKeyPrefix, token)
		s.client.Del(ctx, accessKey, refreshKey)
	}

	// 删除用户的token集合
	return s.client.Del(ctx, userKey).Err()
}
