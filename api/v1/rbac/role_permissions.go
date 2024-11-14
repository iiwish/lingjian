package rbac

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service/rbac"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      获取角色权限
// @Description  获取指定角色的所有权限
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Param        role_id path string true "角色代码"
// @Param        app_code query string true "应用代码"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /roles/{role_id}/permissions [get]
func GetRolePermissions(c *gin.Context) {
	roleID := c.Param("role_id")

	if roleID == "" {
		utils.Error(c, 400, "角色代码不能为空")
		return
	}

	rolePermissionService := &rbac.RolePermissionService{}
	permissions, err := rolePermissionService.GetRolePermissions(roleID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items": permissions,
		"total": len(permissions),
	})
}

// @Summary      修改角色权限
// @Description  为指定角色添加或移除权限
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Param        role_id path string true "角色ID"
// @Param        request body model.PatchRolePerms true "权限修改操作" example:[{"op":"add","value":["权限1ID","权限2ID"]},{"op":"remove","value":["权限3ID","权限4ID"]}]
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /roles/{role_id}/permissions [patch]
func AssignPermissionsToRole(c *gin.Context) {
	var patches model.PatchRolePerms
	if err := c.ShouldBindJSON(&patches); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	roleID := c.Param("role_id")
	if roleID == "" {
		utils.Error(c, 400, "角色代码不能为空")
		return
	}

	operatorID := c.GetUint("user_id")
	if operatorID == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	rolePermissionService := &rbac.RolePermissionService{}

	for _, patch := range patches {
		switch patch.Op {
		case "add":
			if err := rolePermissionService.AddPermissionsToRole(operatorID, roleID, patch.Value); err != nil {
				utils.Error(c, 500, err.Error())
				return
			}
		case "remove":
			if err := rolePermissionService.RemovePermissionsFromRole(operatorID, roleID, patch.Value); err != nil {
				utils.Error(c, 500, err.Error())
				return
			}
		default:
			utils.Error(c, 400, "不支持的操作类型")
			return
		}
	}

	utils.Success(c, nil)
}
