package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/middleware"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(r *gin.RouterGroup) {
	user := r.Group("/user")
	{
		user.GET("/profile", middleware.AuthMiddleware(), GetUserProfile)
		user.PUT("/profile", middleware.AuthMiddleware(), UpdateUserProfile)
	}
}

// GetUserProfile 获取用户信息
// @Summary      获取用户信息
// @Description  获取当前登录用户的信息
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Success      200  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Router       /user/profile [get]
func GetUserProfile(c *gin.Context) {
	// JWT用户访问
	userId := c.GetUint("user_id")
	if userId == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	// 获取用户信息
	var user model.User
	err := model.DB.Get(&user, "SELECT id, username, nickname, email, phone, status FROM sys_users WHERE id = ?", userId)
	if err != nil {
		utils.Error(c, 403, "用户不存在")
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"phone":    user.Phone,
		"status":   user.Status,
		"type":     "jwt",
	})
}

// UpdateUserProfile 更新用户信息
// @Summary      更新用户信息
// @Description  更新当前登录用户的信息
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer Token"
// @Param        App-ID header string true "应用ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Router       /user/profile [put]
func UpdateUserProfile(c *gin.Context) {
	// JWT用户访问
	userId := c.GetUint("user_id")
	if userId == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	// 绑定请求参数
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}
	req.ID = userId

	// 更新用户信息
	_, err := model.DB.NamedExec(`
		UPDATE sys_users SET nickname=:nickname, email=:email, phone=:phone,updater_id=:id, updated_at=NOW()   WHERE id=:id
	`, req)
	if err != nil {
		utils.Error(c, 500, "更新用户信息失败")
		return
	}

	utils.Success(c, nil)
}
