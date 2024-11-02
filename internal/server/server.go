package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/pkg/config"
	"github.com/iiwish/lingjian/pkg/database"
	"github.com/iiwish/lingjian/pkg/redis"
)

type Server struct {
	engine *gin.Engine
	addr   string
}

func NewServer() *Server {
	// 设置gin模式
	gin.SetMode(config.GlobalConfig.Server.Mode)

	engine := gin.Default()

	// 基础中间件
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	// 健康检查
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return &Server{
		engine: engine,
		addr:   fmt.Sprintf(":%s", config.GlobalConfig.Server.Port),
	}
}

// 初始化路由
func (s *Server) initRoutes() {
	// TODO: 添加路由组
	api := s.engine.Group("/api/v1")
	{
		// 后续添加具体路由
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})
	}
}

// 运行服务器
func (s *Server) Run() error {
	s.initRoutes()

	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.engine,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Printf("服务器启动在%s端口", config.GlobalConfig.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭数据库连接
	database.CloseMySQL()
	redis.CloseRedis()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭失败: %v", err)
	}

	log.Println("服务器已关闭")
	return nil
}
