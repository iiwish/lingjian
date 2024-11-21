package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/middleware"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// RegisterSysVarsRoutes 注册系统配置相关路由
func RegisterSysVarsRoutes(r *gin.RouterGroup) {
	sysVars := r.Group("/sys_vars")
	{
		sysVars.GET("", GetSysVars)
		sysVars.GET("/:code", GetSysVarByCode)
		sysVarsRequired := sysVars.Group("/", middleware.AuthMiddleware())
		{
			sysVarsRequired.PUT("", UpdateSysVars)
			sysVarsRequired.POST("", CreateSysVar)
		}
	}
}

// @Summary      获取系统配置
// @Description  获取所有系统配置项
// @Tags         SysVars
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response{data=map[string]string}
// @Failure      500  {object}  utils.Response
// @Router       /sys_vars [get]
func GetSysVars(c *gin.Context) {
	var sysVars []model.SysVar
	err := model.DB.Select(&sysVars, "SELECT * FROM sys_vars")
	if err != nil {
		utils.Error(c, 500, "获取系统配置失败")
		return
	}

	utils.Success(c, sysVars)
}

// @Summary      根据code获取系统配置
// @Description  根据code获取单个系统配置项
// @Tags         SysVars
// @Accept       json
// @Produce      json
// @Param        code path string true "系统配置项的code"
// @Success      200  {object}  utils.Response{data=model.SysVar}
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /sys_vars/{code} [get]
func GetSysVarByCode(c *gin.Context) {
	code := c.Param("code")
	var sysVar model.SysVar
	err := model.DB.Get(&sysVar, "SELECT * FROM sys_vars WHERE key = ?", code)
	if err != nil {
		utils.Error(c, 404, "系统配置项不存在")
		return
	}

	utils.Success(c, sysVar)
}

// @Summary      修改系统配置
// @Description  修改系统配置项
// @Tags         SysVars
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        request body model.SysVar true "系统配置项"
// @Success      200  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /sys_vars [put]
func UpdateSysVars(c *gin.Context) {
	var req model.SysVar
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	// 校验权限
	userId := c.GetUint("user_id")
	if userId == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	req.UpdaterID = userId

	// 更新系统配置
	_, err := model.DB.Exec("UPDATE sys_vars SET name = :name, code = :code,value = :value, description = :description, updater_id = :updater_id, updated_at = NOW() WHERE id = :id", req)
	if err != nil {
		utils.Error(c, 500, "系统配置更新失败")
		return
	}

	utils.Success(c, gin.H{"message": "系统配置更新成功"})
}

// @Summary      新建系统配置
// @Description  新建系统配置项
// @Tags         SysVars
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        request body model.SysVar true "系统配置项"
// @Success      200  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /sys_vars [post]
func CreateSysVar(c *gin.Context) {
	var req model.SysVar
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	// 校验权限
	userId := c.GetUint("user_id")
	if userId == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	// 开启事务
	tx := model.DB.MustBegin()
	defer tx.Rollback()

	// 新建系统配置
	result, err := tx.NamedExec(`
		INSERT INTO sys_vars (name, code, value, description, creator_id, updater_id, created_at, updated_at)
		VALUES (:name, :code, :value, :description, :creator_id, :updater_id, NOW(), NOW())
	`, req)
	if err != nil {
		utils.Error(c, 500, "系统配置新建失败")
		return
	}

	// 获取新创建的系统配置ID
	id, err := result.LastInsertId()
	if err != nil {
		utils.Error(c, 500, "系统配置新建失败")
		return
	}

	utils.Success(c, gin.H{"id": id, "message": "系统配置新建成功"})
}
