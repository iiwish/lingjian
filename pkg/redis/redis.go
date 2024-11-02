package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var RDB *redis.Client

func InitRedis() {
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")

	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Fatalf("连接Redis失败: %v", err)
	}

	log.Println("Redis连接成功")
}

func CloseRedis() {
	if RDB != nil {
		if err := RDB.Close(); err != nil {
			log.Printf("关闭Redis连接失败: %v", err)
			return
		}
		log.Println("Redis连接已关闭")
	}
}

// Set 设置键值对
func Set(ctx context.Context, key string, value interface{}, ttl int) error {
	return RDB.Set(ctx, key, value, 0).Err()
}

// Get 获取值
func Get(ctx context.Context, key string) (string, error) {
	return RDB.Get(ctx, key).Result()
}

// Del 删除键
func Del(ctx context.Context, key string) error {
	return RDB.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := RDB.Exists(ctx, key).Result()
	return result > 0, err
}
