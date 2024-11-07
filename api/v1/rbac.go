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

// @Summary      创建角色
// @Description  创建新的角色
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Param        request body CreateRoleRequest true "创建角色请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /rbac/roles [post]
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

// @Summary      创建权限
// @Description  创建新的权限
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Param        request body CreatePermissionRequest true "创建权限请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /rbac/permissions [post]
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

// @Summary      为用户分配角色
// @Description  将指定角色分配给用户
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Param        user_id path int true "用户ID"
// @Param        role_id path int true "角色ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /rbac/users/{user_id}/roles/{role_id} [post]
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

// @Summary      为角色分配权限
// @Description  为指定角色分配权限
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Param        role_id path int true "角色ID"
// @Param        permission_id path int true "权限ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /rbac/roles/{role_id}/permissions/{permission_id} [post]
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

// @Summary      获取用户角色
// @Description  获取指定用户的所有角色
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Param        user_id path int true "用户ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /rbac/users/{user_id}/roles [get]
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

// @Summary      获取角色权限
// @Description  获取指定角色的所有权限
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Param        role_id path int true "角色ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /rbac/roles/{role_id}/permissions [get]
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
