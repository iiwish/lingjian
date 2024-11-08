package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(r *gin.RouterGroup) {
	user := r.Group("/user")
	{
		user.GET("/profile", GetUserProfile)
	}
}

// GetUserProfile 获取用户信息
// @Summary      获取用户信息
// @Description  获取当前登录用户的信息
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Router       /user/profile [get]
func GetUserProfile(c *gin.Context) {
	// 检查是否是OAuth2访问
	if isOAuth, _ := c.Get("is_oauth"); isOAuth.(bool) {
		clientID, _ := c.Get("client_id")
		utils.Success(c, gin.H{
			"client_id": clientID,
			"type":      "oauth2",
		})
		return
	}

	// JWT用户访问
	userId := c.GetUint("user_id")
	if userId == 0 {
		utils.Error(c, 401, "未授权")
		return
	}

	// 获取用户信息
	var user model.User
	err := model.DB.Get(&user, "SELECT id, username, email, phone, status FROM users WHERE id = ?", userId)
	if err != nil {
		utils.Error(c, 401, "用户不存在")
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
