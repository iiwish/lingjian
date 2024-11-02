package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/service"
	"github.com/iiwish/lingjian/pkg/utils"
)

// RegisterRBACRoutes 注册RBAC相关路由
func RegisterRBACRoutes(r *gin.RouterGroup) {
	rbac := r.Group("/rbac")
	{
		rbac.POST("/roles", CreateRole)
		rbac.POST("/permissions", CreatePermission)
		rbac.POST("/users/:user_id/roles/:role_id", AssignRoleToUser)
		rbac.POST("/roles/:role_id/permissions/:permission_id", AssignPermissionToRole)
		rbac.GET("/users/:user_id/roles", GetUserRoles)
		rbac.GET("/roles/:role_id/permissions", GetRolePermissions)
	}
}

type CreateRoleRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

// CreateRole 创建角色
func CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	rbacService := &service.RBACService{}
	if err := rbacService.CreateRole(req.Name, req.Code); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

type CreatePermissionRequest struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Type   string `json:"type" binding:"required"`
	Path   string `json:"path" binding:"required"`
	Method string `json:"method" binding:"required"`
}

// CreatePermission 创建权限
func CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	rbacService := &service.RBACService{}
	if err := rbacService.CreatePermission(req.Name, req.Code, req.Type, req.Path, req.Method); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// AssignRoleToUser 为用户分配角色
func AssignRoleToUser(c *gin.Context) {
	userID := utils.ParseUint(c.Param("user_id"))
	roleID := utils.ParseUint(c.Param("role_id"))

	rbacService := &service.RBACService{}
	if err := rbacService.AssignRoleToUser(userID, roleID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// AssignPermissionToRole 为角色分配权限
func AssignPermissionToRole(c *gin.Context) {
	roleID := utils.ParseUint(c.Param("role_id"))
	permissionID := utils.ParseUint(c.Param("permission_id"))

	rbacService := &service.RBACService{}
	if err := rbacService.AssignPermissionToRole(roleID, permissionID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// GetUserRoles 获取用户的所有角色
func GetUserRoles(c *gin.Context) {
	userID := utils.ParseUint(c.Param("user_id"))

	rbacService := &service.RBACService{}
	roles, err := rbacService.GetUserRoles(userID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, roles)
}

// GetRolePermissions 获取角色的所有权限
func GetRolePermissions(c *gin.Context) {
	roleID := utils.ParseUint(c.Param("role_id"))

	rbacService := &service.RBACService{}
	permissions, err := rbacService.GetRolePermissions(roleID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, permissions)
}
