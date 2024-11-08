package middleware

import (
	"fmt"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// pathMatch 检查请求路径是否匹配权限路径
func pathMatch(permPath, reqPath string) bool {
	// 如果权限路径包含通配符
	if strings.Contains(permPath, "*") {
		permParts := strings.Split(permPath, "/")
		reqParts := strings.Split(reqPath, "/")

		// 如果路径段数不同，且权限路径最后一段不是通配符，则不匹配
		if len(permParts) != len(reqParts) && permParts[len(permParts)-1] != "*" {
			return false
		}

		// 逐段比较
		for i := 0; i < len(permParts) && i < len(reqParts); i++ {
			if permParts[i] == "*" {
				continue
			}
			if permParts[i] != reqParts[i] {
				return false
			}
		}
		return true
	}

	// 不包含通配符时进行精确匹配
	return path.Clean(permPath) == path.Clean(reqPath)
}

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
		reqPath := c.Request.URL.Path
		method := c.Request.Method

		// 在测试环境中，admin用户拥有所有权限
		if gin.Mode() == gin.TestMode && userID.(uint) == 1 {
			c.Next()
			return
		}

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
		args := make([]interface{}, len(roleIDs)+1) // +1 for method
		for i := range roleIDs {
			placeholders[i] = "?"
			args[i] = roleIDs[i]
		}
		args[len(roleIDs)] = method

		// 查询角色的所有权限
		var permissions []struct {
			Path   string
			Method string
		}
		query := fmt.Sprintf(`
			SELECT DISTINCT p.path, p.method FROM permissions p
			INNER JOIN role_permissions rp ON p.id = rp.permission_id
			WHERE rp.role_id IN (%s)
			AND p.method = ?
			AND p.status = 1
		`, strings.Join(placeholders, ","))

		err = model.DB.Select(&permissions, query, args...)
		if err != nil {
			utils.Error(c, 500, "服务器错误")
			c.Abort()
			return
		}

		// 检查是否有匹配的权限
		hasPermission := false
		for _, perm := range permissions {
			if pathMatch(perm.Path, reqPath) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			utils.Error(c, 403, "没有访问权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
