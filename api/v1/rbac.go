package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/api/v1/rbac"
)

// RegisterRBACRoutes 注册RBAC相关路由
func RegisterRBACRoutes(r *gin.RouterGroup) {
	// 用户相关路由
	r.GET("/users", rbac.GetUsers)
	r.POST("/users", rbac.CreateUser)
	r.GET("/users/:user_id", rbac.GetUser)
	r.PUT("/users/:user_id", rbac.UpdateUser)
	r.DELETE("/users/:user_id", rbac.DeleteUser)

	// 角色相关路由
	r.GET("/roles", rbac.ListRoles)
	r.POST("/roles", rbac.CreateRole)
	r.PUT("/roles/:role_id", rbac.UpdateRole)
	r.DELETE("/roles/:role_id", rbac.DeleteRole)

	// 权限相关路由
	r.GET("/permissions", rbac.ListPermissions)
	r.POST("/permissions", rbac.CreatePermission)
	r.PUT("/permissions/:permission_id", rbac.UpdatePermission)
	r.DELETE("/permissions/:permission_id", rbac.DeletePermission)

	// 用户角色相关路由
	r.GET("/users/:user_id/roles", rbac.GetUserRoles)
	r.PATCH("/users/:user_id/roles", rbac.AssignRoleToUser)

	// 角色权限相关路由
	r.GET("/roles/:role_id/permissions", rbac.GetRolePermissions)
	r.PATCH("/roles/:role_id/permissions", rbac.AssignPermissionsToRole)
}
