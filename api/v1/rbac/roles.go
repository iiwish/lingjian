package rbac

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service/rbac"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      获取角色列表
// @Description  获取所有角色列表
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /roles [get]
func ListRoles(c *gin.Context) {
	rbacService := &rbac.RoleService{}
	roles, err := rbacService.ListRoles()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, roles)
}

// @Summary      创建角色
// @Description  创建新的角色
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Param        request body model.Role true "创建角色请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /roles [post]
func CreateRole(c *gin.Context) {
	var req model.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	operatorID := c.GetUint("user_id")
	if operatorID == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	rbacService := &rbac.RoleService{}
	if err := rbacService.CreateRole(operatorID, &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      更新角色
// @Description  更新指定角色的信息
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Param        role_id path int true "角色ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /roles/{role_id} [put]
func UpdateRole(c *gin.Context) {
	roleID := utils.ParseUint(c.Param("role_id"))
	var req model.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	operatorID := c.GetUint("user_id")
	if operatorID == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	rbacService := &rbac.RoleService{}
	if err := rbacService.UpdateRole(operatorID, roleID, &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      删除角色
// @Description  删除指定角色
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Param        role_id path int true "角色ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /roles/{role_id} [delete]
func DeleteRole(c *gin.Context) {
	roleID := utils.ParseUint(c.Param("role_id"))

	operatorID := c.GetUint("user_id")
	if operatorID == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	rbacService := &rbac.RoleService{}
	if err := rbacService.DeleteRole(operatorID, roleID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}
