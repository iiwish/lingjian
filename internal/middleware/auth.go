package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/pkg/store"
	"github.com/iiwish/lingjian/pkg/utils"
)

var globalStore store.Store

// SetStore 设置全局存储实例
func SetStore(s store.Store) {
	globalStore = s
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if globalStore == nil {
			globalStore = store.NewRedisStore()
		}

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

		token := parts[1]

		// 从存储中验证token
		userId, err := globalStore.VerifyToken(token, "access")
		if err != nil {
			utils.Error(c, 401, "无效的token")
			c.Abort()
			return
		}

		// 解析JWT获取详细信息
		claims, err := utils.ParseToken(token, utils.AccessToken)
		if err != nil {
			utils.Error(c, 401, "无效的token")
			c.Abort()
			return
		}

		// 验证token中的用户ID是否匹配
		if claims.UserID != userId {
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
