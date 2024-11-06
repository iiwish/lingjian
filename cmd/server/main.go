package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	v1 "github.com/iiwish/lingjian/api/v1"
	_ "github.com/iiwish/lingjian/docs" // 导入swagger文档
	"github.com/iiwish/lingjian/internal/middleware"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/queue"
	"github.com/iiwish/lingjian/pkg/redis"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           灵简服务端API
// @version         1.0
// @description     灵简服务端API文档

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func init() {
	// 加载配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 初始化数据库连接
	if err := model.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化Redis连接
	redis.InitRedis()

	// 初始化RabbitMQ连接
	if err := queue.InitRabbitMQ(); err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}

	// 初始化认证服务
	v1.InitAuthService()
}

func main() {
	// 设置gin模式
	gin.SetMode(viper.GetString("server.mode"))

	// 创建gin引擎
	r := gin.Default()

	// 基础中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API路由
	api := r.Group("/api")
	{
		// v1版本API
		v1Group := api.Group("/v1")
		{
			// 注册认证相关路由
			v1.RegisterAuthRoutes(v1Group)

			// 需要认证的路由
			authorized := v1Group.Group("/")
			authorized.Use(middleware.AuthMiddleware())
			{
				// 需要RBAC权限控制的路由
				rbacProtected := authorized.Group("/")
				rbacProtected.Use(middleware.RBACMiddleware())
				{
					// 注册RBAC相关路由
					v1.RegisterRBACRoutes(rbacProtected)
					// 注册应用相关路由
					v1.RegisterAppRoutes(rbacProtected)
					// 注册配置相关路由
					v1.RegisterConfigRoutes(rbacProtected)
					// 注册任务相关路由
					v1.RegisterTaskRoutes(rbacProtected)
				}
			}
		}
	}

	// 启动服务器
	port := viper.GetString("server.port")
	log.Printf("Server is running on http://localhost:%s", port)
	log.Printf("Swagger文档地址: http://localhost:%s/swagger/index.html", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
