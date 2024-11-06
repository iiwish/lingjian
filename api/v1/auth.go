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

// @Summary      获取验证码
// @Description  生成图形验证码
// @Tags         认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response{data=service.CaptchaResponse}
// @Failure      400  {object}  utils.Response
// @Router       /auth/captcha [get]
func GetCaptcha(c *gin.Context) {
	resp, err := authService.GenerateCaptcha()
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, resp)
}

// @Summary      用户登录
// @Description  使用用户名和密码登录
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request body service.LoginRequest true "登录请求参数"
// @Success      200  {object}  utils.Response{data=service.LoginResponse}
// @Failure      400  {object}  utils.Response
// @Router       /auth/login [post]
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

// @Summary      刷新令牌
// @Description  使用刷新令牌获取新的访问令牌
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        X-Refresh-Token header string true "刷新令牌"
// @Success      200  {object}  utils.Response{data=service.TokenResponse}
// @Failure      400  {object}  utils.Response
// @Router       /auth/refresh [post]
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

// @Summary      用户登出
// @Description  注销用户登录状态
// @Tags         认证
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /auth/logout [post]
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
