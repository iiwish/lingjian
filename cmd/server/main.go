package main

import (
	"log"

	"github.com/iiwish/lingjian/internal/server"
	"github.com/iiwish/lingjian/pkg/config"
	"github.com/iiwish/lingjian/pkg/database"
	"github.com/iiwish/lingjian/pkg/redis"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化数据库连接
	database.InitMySQL()
	redis.InitRedis()

	// 创建并运行服务器
	srv := server.NewServer()
	if err := srv.Run(); err != nil {
		log.Fatalf("服务器运行错误: %v", err)
	}
}
