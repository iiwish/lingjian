package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// RBACMiddleware RBAC权限控制中间件
func RBACMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.Error(c, 401, "未授权")
			c.Abort()
			return
		}

		// 获取当前请求的路径和方法
		path := c.Request.URL.Path
		method := c.Request.Method

		// 查询用户的角色
		var roleIDs []uint
		err := model.DB.Select(&roleIDs, `
			SELECT role_id FROM user_roles WHERE user_id = ?
		`, userID)
		if err != nil {
			utils.Error(c, 500, "服务器错误")
			c.Abort()
			return
		}

		// 如果用户没有任何角色
		if len(roleIDs) == 0 {
			utils.Error(c, 403, "没有访问权限")
			c.Abort()
			return
		}

		// 构建角色ID的占位符
		placeholders := make([]string, len(roleIDs))
		args := make([]interface{}, len(roleIDs)+2) // +2 for path and method
		for i := range roleIDs {
			placeholders[i] = "?"
			args[i] = roleIDs[i]
		}
		args[len(roleIDs)] = path
		args[len(roleIDs)+1] = method

		// 查询角色的权限
		var count int
		query := fmt.Sprintf(`
			SELECT COUNT(*) FROM permissions p
			INNER JOIN role_permissions rp ON p.id = rp.permission_id
			WHERE rp.role_id IN (%s)
			AND p.path = ?
			AND p.method = ?
			AND p.status = 1
		`, strings.Join(placeholders, ","))

		err = model.DB.Get(&count, query, args...)
		if err != nil {
			utils.Error(c, 500, "服务器错误")
			c.Abort()
			return
		}

		// 如果没有找到匹配的权限
		if count == 0 {
			utils.Error(c, 403, "没有访问权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
