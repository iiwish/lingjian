package middleware

import (
	"log"
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
			log.Printf("Auth: 存储实例未初始化")
			utils.Error(c, 500, "存储实例未初始化")
			c.Abort()
			return
		}

		auth := c.GetHeader("Authorization")
		if auth == "" {
			log.Printf("Auth: 未提供Authorization头")
			utils.Error(c, 401, "未授权")
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			log.Printf("Auth: 无效的Authorization格式")
			utils.Error(c, 401, "无效的授权格式")
			c.Abort()
			return
		}

		token := parts[1]
		log.Printf("Auth: 开始验证token: %s", token)

		// 先尝试作为JWT token验证
		claims, err := utils.ParseToken(token, utils.AccessToken)
		if err == nil {
			// 验证token是否在存储中
			userId, err := globalStore.VerifyToken(token, "access")
			if err != nil {
				log.Printf("Auth: token验证失败 - %v", err)
				utils.Error(c, 401, "无效的token")
				c.Abort()
				return
			}

			// 验证token中的用户ID是否匹配
			if claims.UserID != userId {
				log.Printf("Auth: token中的用户ID不匹配 - token: %d, store: %d", claims.UserID, userId)
				utils.Error(c, 401, "无效的token")
				c.Abort()
				return
			}

			// 将用户信息存储到上下文中
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)

			// 如果token中包含角色信息，也存储到上下文中
			if claims.RoleCode != "" {
				log.Printf("Auth: 设置角色信息 - UserID: %d, RoleCode: %s", claims.UserID, claims.RoleCode)
				c.Set("role_code", claims.RoleCode)
			} else {
				log.Printf("Auth: token中没有角色信息 - UserID: %d", claims.UserID)
				// 检查是否是切换角色的请求，如果是则允许通过
				if c.Request.URL.Path == "/api/v1/auth/switch-role" && c.Request.Method == "POST" {
					log.Printf("Auth: 允许切换角色请求通过")
					c.Next()
					return
				}
			}

			c.Next()
			return
		} else {
			log.Printf("Auth: JWT验证失败 - %v", err)
		}

		// 如果JWT验证失败，尝试作为OAuth2 token验证
		if clientID, _, err := globalStore.GetRefreshToken(token); err == nil {
			// 这是一个有效的OAuth2 token
			log.Printf("Auth: OAuth2 token验证成功 - ClientID: %s", clientID)
			c.Set("client_id", clientID)
			c.Set("is_oauth", true)
			c.Next()
			return
		}

		log.Printf("Auth: token验证完全失败")
		utils.Error(c, 401, "无效的token")
		c.Abort()
	}
}
