package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service"
	"github.com/iiwish/lingjian/pkg/utils"
)

// RegisterRBACRoutes 注册RBAC相关路由
func RegisterRBACRoutes(r *gin.RouterGroup) {
	// 角色相关路由
	r.POST("/roles", CreateRole)
	r.GET("/roles/:role_code/permissions", GetRolePermissions)
	r.POST("/roles/:role_code/permissions", AssignPermissionsToRole)

	// 权限相关路由
	r.POST("/permissions", CreatePermission)

	// 用户角色相关路由
	r.POST("/users/:user_id/roles/:role_id", AssignRoleToUser)
	r.GET("/users/:user_id/roles", GetUserRoles)
}

// @Summary      创建角色
// @Description  创建新的角色
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Param        request body model.CreateRoleRequest true "创建角色请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /roles [post]
func CreateRole(c *gin.Context) {
	var req model.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	rbacService := &service.RBACService{}
	if err := rbacService.CreateRole(&req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      创建权限
// @Description  创建新的权限
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Param        request body model.CreatePermissionRequest true "创建权限请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /permissions [post]
func CreatePermission(c *gin.Context) {
	var req model.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	rbacService := &service.RBACService{}
	if err := rbacService.CreatePermission(&req); err != nil {
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
// @Router       /users/{user_id}/roles/{role_id} [post]
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
// @Param        role_code path string true "角色代码"
// @Param        request body model.AssignPermissionsRequest true "分配权限请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /roles/{role_code}/permissions [post]
func AssignPermissionsToRole(c *gin.Context) {
	var req model.AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	roleCode := c.Param("role_code")
	if roleCode == "" {
		utils.Error(c, 400, "角色代码不能为空")
		return
	}

	rbacService := &service.RBACService{}
	if err := rbacService.AssignPermissionsToRole(roleCode, req.AppCode, req.PermissionCodes); err != nil {
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
// @Router       /users/{user_id}/roles [get]
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
// @Param        role_code path string true "角色代码"
// @Param        app_code query string true "应用代码"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /roles/{role_code}/permissions [get]
func GetRolePermissions(c *gin.Context) {
	roleCode := c.Param("role_code")
	appCode := c.Query("app_code")

	if roleCode == "" || appCode == "" {
		utils.Error(c, 400, "角色代码和应用代码不能为空")
		return
	}

	rbacService := &service.RBACService{}
	permissions, err := rbacService.GetRolePermissions(roleCode, appCode)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, permissions)
}
