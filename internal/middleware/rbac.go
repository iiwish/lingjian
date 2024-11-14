package middleware

import (
	"log"
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
		// 检查用户是否已认证
		userId, exists := c.Get("user_id")
		if !exists {
			log.Printf("RBAC: 用户未认证")
			utils.Error(c, 403, "未授权")
			c.Abort()
			return
		}

		// 获取当前请求的路径和方法
		reqPath := c.Request.URL.Path
		method := c.Request.Method

		// 获取应用 ID
		appID := c.GetHeader("App-ID")
		if appID == "" {
			log.Printf("RBAC: 缺少应用 ID")
			utils.Error(c, 400, "缺少应用 ID")
			c.Abort()
			return
		}

		// 查询角色的所有权限
		var permissions []struct {
			Path   string
			Method string
		}
		query := `
            SELECT DISTINCT p.path, p.method FROM sys_permissions p
            INNER JOIN sys_role_permissions rp ON p.id = rp.permission_id
            INNER JOIN sys_user_roles ur ON rp.role_id = ur.role_id
            INNER JOIN sys_config_menus m ON p.menu_id = m.id
            WHERE ur.user_id = ?
            AND p.method = ?
            AND p.status = 1
            AND m.app_id = ?
		`

		err := model.DB.Select(&permissions, query, userId, method)
		if err != nil {
			log.Printf("RBAC: 查询权限失败 - %v", err)
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
			log.Printf("RBAC: 用户 %v 访问 %s %s 被拒绝", userId, method, reqPath)
			utils.Error(c, 403, "没有访问权限")
			c.Abort()
			return
		}

		log.Printf("RBAC: 用户 %v 访问 %s %s 通过", userId, method, reqPath)
		c.Next()
	}
}
