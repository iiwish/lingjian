package middleware

import (
	"log"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// pathMatch 检查请求路径是否匹配权限路径
func pathMatch(pattern, path string) bool {
	// 将路径模式转换为正则表达式
	regexPattern := "^" + regexp.QuoteMeta(pattern) + "$"
	regexPattern = strings.Replace(regexPattern, `\*`, `.*`, -1)
	regexPattern = strings.Replace(regexPattern, `/:id`, `/\d+`, -1)
	regexPattern = strings.Replace(regexPattern, `/:dim_id`, `/\d+`, -1)
	regexPattern = strings.Replace(regexPattern, `/:menu_id`, `/\d+`, -1)
	regexPattern = strings.Replace(regexPattern, `/:table_id`, `/\d+`, -1)
	regexPattern = strings.Replace(regexPattern, `/:model_id`, `/\d+`, -1)
	regexPattern = strings.Replace(regexPattern, `/:form_id`, `/\d+`, -1)

	// 编译正则表达式
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		log.Printf("RBAC: 正则表达式编译失败 - %v", err)
		return false
	}

	// 检查路径是否匹配
	return regex.MatchString(path)
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

		c.Set("app_id", utils.ParseUint(appID))

		// 查询角色的所有权限
		var permissions []struct {
			Path   string
			Method string
		}
		query := `
            SELECT DISTINCT p.path, p.method FROM sys_permissions p
            INNER JOIN sys_role_permissions rp ON p.id = rp.permission_id
            INNER JOIN sys_user_roles ur ON rp.role_id = ur.role_id
            WHERE ur.user_id = ?
            AND p.method = ?
            AND p.status = 1
		`
		log.Printf("RBAC: 查询权限 - %v", method)

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
