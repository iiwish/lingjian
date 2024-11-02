package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/pkg/utils"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			utils.Error(c, 401, "未授权")
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Error(c, 401, "无效的授权格式")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1], utils.AccessToken)
		if err != nil {
			utils.Error(c, 401, "无效的token")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
