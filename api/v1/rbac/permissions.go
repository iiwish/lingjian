package rbac

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service/rbac"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      获取权限列表
// @Description  获取所有权限列表
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /permissions [get]
func ListPermissions(c *gin.Context) {
	permissionService := &rbac.PermissionService{}
	permissions, err := permissionService.ListPermissions()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, permissions)
}

// @Summary      创建权限
// @Description  创建新的权限
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        request body model.Permission true "创建权限请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /permissions [post]
func CreatePermission(c *gin.Context) {
	var req model.Permission
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	operatorID := c.GetUint("user_id")
	if operatorID == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	permissionService := &rbac.PermissionService{}
	if err := permissionService.CreatePermission(operatorID, &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      更新权限
// @Description  更新指定权限的信息
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        permission_id path int true "权限ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /permissions/{permission_id} [put]
func UpdatePermission(c *gin.Context) {
	permissionID := utils.ParseUint(c.Param("permission_id"))
	var req model.Permission
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	operatorID := c.GetUint("user_id")
	if operatorID == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	permissionService := &rbac.PermissionService{}
	if err := permissionService.UpdatePermission(operatorID, permissionID, &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      删除权限
// @Description  删除指定权限
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        permission_id path int true "权限ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /permissions/{permission_id} [delete]
func DeletePermission(c *gin.Context) {
	permissionID := utils.ParseUint(c.Param("permission_id"))

	operatorID := c.GetUint("user_id")
	if operatorID == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	permissionService := &rbac.PermissionService{}
	if err := permissionService.DeletePermission(operatorID, permissionID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}
