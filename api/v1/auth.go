package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/service"
	"github.com/iiwish/lingjian/pkg/utils"
)

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", Login)
		auth.POST("/refresh", RefreshToken)
	}
}

// Login 用户登录
func Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	authService := &service.AuthService{}
	resp, err := authService.Login(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, resp)
}

// RefreshToken 刷新访问令牌
func RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		utils.Error(c, 400, "刷新令牌不能为空")
		return
	}

	authService := &service.AuthService{}
	resp, err := authService.RefreshToken(refreshToken)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, resp)
}
