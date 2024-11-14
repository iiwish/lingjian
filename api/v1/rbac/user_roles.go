package rbac

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service/rbac"
	"github.com/iiwish/lingjian/pkg/utils"
)

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

	userRoleService := &rbac.UserRoleService{}
	roles, err := userRoleService.GetUserRoles(userID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, roles)
}

// @Summary      为用户分配角色
// @Description  将指定角色分配给用户
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Param        user_id path int true "用户ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /users/{user_id}/roles/ [patch]
func AssignRoleToUser(c *gin.Context) {
	var patches model.PatchRolePerms
	if err := c.ShouldBindJSON(&patches); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	userID := utils.ParseUint(c.Param("user_id"))

	// 检查操作者权限
	operatorID := c.GetUint("user_id")
	if operatorID == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	userRoleService := &rbac.UserRoleService{}

	for _, patch := range patches {
		switch patch.Op {
		case "add":
			if err := userRoleService.AddRoleToUser(operatorID, userID, patch.Value); err != nil {
				utils.Error(c, 500, err.Error())
				return
			}
		case "remove":
			if err := userRoleService.RemoveRolesFromUser(operatorID, userID, patch.Value); err != nil {
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
