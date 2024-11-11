package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/middleware"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service"
	"github.com/iiwish/lingjian/pkg/store"
	"github.com/iiwish/lingjian/pkg/utils"
)

var authService *service.AuthService

// InitAuthService 初始化认证服务
func InitAuthService(s store.Store) {
	authService = service.NewAuthService(s)
	middleware.SetStore(s) // 添加这行代码
}

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		// 不需要认证的路由
		auth.GET("/captcha", GetCaptcha)
		auth.POST("/login", Login)
		auth.POST("/refresh", RefreshToken)

		// 需要认证的路由
		authRequired := auth.Group("/", middleware.AuthMiddleware())
		{
			authRequired.POST("/logout", Logout)
			authRequired.GET("/userinfo", GetUserInfo)
		}

		// OAuth2.0相关路由
		oauth := auth.Group("/oauth")
		{
			oauth.GET("/authorize", AuthorizeHandler)
			oauth.POST("/authorize", ConfirmAuthorize)
			oauth.POST("/token", TokenHandler)
		}
	}
}

// @Summary      获取用户详细信息
// @Description  获取当前登录用户的详细信息，包括角色和权限
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Router       /auth/userinfo [get]
func GetUserInfo(c *gin.Context) {
	// 打印context中的信息
	fmt.Println("context:", c.Keys)

	// 从上下文中获取用户ID
	userId := c.GetUint("user_id")
	if userId == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	// 获取用户基本信息
	var user model.User
	err := model.DB.Get(&user, `
		SELECT id, username, nickname, avatar, email 
		FROM sys_users 
		WHERE id = ?`, userId)
	if err != nil {
		utils.Error(c, 403, "用户不存在")
		return
	}

	// 获取用户角色
	var roles []string
	err = model.DB.Select(&roles, `
		SELECT r.code 
		FROM sys_roles r 
		JOIN sys_user_roles ur ON r.id = ur.role_id 
		WHERE ur.user_id = ? AND r.status = 1`, userId)
	if err != nil {
		utils.Error(c, 500, "获取角色信息失败")
		return
	}

	// 获取用户权限
	var permissions []string
	err = model.DB.Select(&permissions, `
		SELECT DISTINCT p.code 
		FROM sys_permissions p 
		JOIN sys_role_permissions rp ON p.id = rp.permission_id 
		JOIN sys_user_roles ur ON rp.role_id = ur.role_id 
		WHERE ur.user_id = ? AND p.status = 1`, userId)
	if err != nil {
		utils.Error(c, 500, "获取权限信息失败")
		return
	}

	utils.Success(c, gin.H{
		"id":          user.ID,
		"username":    user.Username,
		"nickname":    user.Nickname,
		"avatar":      user.Avatar,
		"email":       user.Email,
		"roles":       roles,
		"permissions": permissions,
	})
}

// @Summary      获取验证码
// @Description  生成图形验证码
// @Tags         Auth
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
// @Tags         Auth
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
// @Tags         Auth
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
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /auth/logout [post]
func Logout(c *gin.Context) {
	// 从JWT中获取用户ID
	userId := c.GetUint("user_id")
	if userId == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	if err := authService.Logout(userId); err != nil {
		utils.Error(c, 400, "登出失败")
		return
	}

	utils.Success(c, gin.H{"message": "登出成功"})
}

// @Summary      OAuth2授权页面
// @Description  获取OAuth2授权页面
// @Tags         OAuth2
// @Accept       json
// @Produce      html
// @Param        client_id     query    string  true  "客户端ID"
// @Param        redirect_uri  query    string  true  "重定向URI"
// @Param        response_type query    string  true  "响应类型(code)"
// @Param        scope         query    string  true  "权限范围"
// @Param        state         query    string  false "状态"
// @Success      200  {string} string "授权页面HTML"
// @Failure      400  {object} utils.Response
// @Router       /auth/oauth/authorize [get]
func AuthorizeHandler(c *gin.Context) {
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")
	responseType := c.Query("response_type")
	scope := c.Query("scope")
	state := c.Query("state")

	// 验证参数
	if clientID == "" || redirectURI == "" || responseType != "code" {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	// 渲染授权页面
	html := `
		<html>
			<head>
				<title>授权页面</title>
			</head>
			<body>
				<h2>授权请求</h2>
				<p>应用 ` + clientID + ` 请求访问您的以下权限：</p>
				<p>权限范围：` + scope + `</p>
				<form method="post" action="/api/v1/auth/oauth/authorize">
					<input type="hidden" name="client_id" value="` + clientID + `">
					<input type="hidden" name="redirect_uri" value="` + redirectURI + `">
					<input type="hidden" name="response_type" value="` + responseType + `">
					<input type="hidden" name="scope" value="` + scope + `">
					<input type="hidden" name="state" value="` + state + `">
					<button type="submit" name="approved" value="true">同意授权</button>
					<button type="submit" name="approved" value="false">���绝授权</button>
				</form>
			</body>
		</html>
	`
	c.Header("Content-Type", "text/html")
	c.String(200, html)
}

// @Summary      确认OAuth2授权
// @Description  处理用户的授权确认
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @Param        request body service.AuthorizeRequest true "授权请求参数"
// @Success      302  {string} string "重定向到客户端"
// @Failure      400  {object} utils.Response
// @Router       /auth/oauth/authorize [post]
func ConfirmAuthorize(c *gin.Context) {
	var req service.AuthorizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	// 处理授权
	redirectURL, err := authService.HandleAuthorize(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	// 重定向到客户端
	c.Redirect(302, redirectURL)
}

// @Summary      OAuth2令牌
// @Description  获取访问令牌
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @Param        request body service.TokenRequest true "令牌请求参数"
// @Success      200  {object} utils.Response{data=service.TokenResponse}
// @Failure      400  {object}  utils.Response
// @Router       /auth/oauth/token [post]
func TokenHandler(c *gin.Context) {
	var req service.TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	resp, err := authService.HandleToken(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, resp)
}
