package rbac

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service/rbac"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      获取用户列表
// @Description  获取所有用户列表
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /users [get]
func GetUsers(c *gin.Context) {
	userService := &rbac.UserService{}
	users, err := userService.GetUsers()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, users)
}

// @Summary      创建用户
// @Description  创建新的用户
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        request body model.User true "创建用户请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /users [post]
func CreateUser(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	operatorID := c.GetUint("user_id")
	if operatorID == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	userService := &rbac.UserService{}
	if err := userService.CreateUser(operatorID, &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      更新用户
// @Description  更新指定用户的信息
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        user_id path int true "用户ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /users/{user_id} [put]
func UpdateUser(c *gin.Context) {
	userID := utils.ParseUint(c.Param("user_id"))
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	operatorID := c.GetUint("user_id")
	if operatorID == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	userService := &rbac.UserService{}
	if err := userService.UpdateUser(operatorID, userID, &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      删除用户
// @Description  删除指定用户
// @Tags         RBAC
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        user_id path int true "用户ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /users/{user_id} [delete]
func DeleteUser(c *gin.Context) {
	userID := utils.ParseUint(c.Param("user_id"))

	operatorID := c.GetUint("user_id")
	if operatorID == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	userService := &rbac.UserService{}
	if err := userService.DeleteUser(operatorID, userID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}
