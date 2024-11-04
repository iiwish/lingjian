package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/service"
	"github.com/iiwish/lingjian/pkg/utils"
)

var authService *service.AuthService

// InitAuthService 初始化认证服务
func InitAuthService() {
	authService = service.NewAuthService()
}

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.GET("/captcha", GetCaptcha)
		auth.POST("/login", Login)
		auth.POST("/refresh", RefreshToken)
		auth.POST("/logout", Logout)
	}
}

// GetCaptcha 获取验证码
func GetCaptcha(c *gin.Context) {
	resp, err := authService.GenerateCaptcha()
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, resp)
}

// Login 用户登录
func Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

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

	resp, err := authService.RefreshToken(refreshToken)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, resp)
}

// Logout 用户登出
func Logout(c *gin.Context) {
	// 从JWT中获取用户ID
	userId := c.GetUint("user_id")
	if userId == 0 {
		utils.Error(c, 401, "未授权")
		return
	}

	if err := authService.Logout(userId); err != nil {
		utils.Error(c, 400, "登出失败")
		return
	}

	utils.Success(c, gin.H{"message": "登出成功"})
}
