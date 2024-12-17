package config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      获取数据表配置详情
// @Description  获取指定数据表配置的详细信息
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        table_id path int true "配置ID"
// @Success      200  {object}  model.ConfigTable
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{table_id} [get]
func (api *ConfigAPI) GetTable(c *gin.Context) {
	id := utils.ParseUint(c.Param("table_id"))
	if id == 0 {
		utils.Error(c, 400, "invalid id")
		return
	}

	table, err := api.configService.GetTable(uint(id))
	if err != nil {
		utils.ServerError(c, err)
		return
	}

	utils.Success(c, table)
}

// @Summary      创建数据表配置
// @Description  创建新的数据表配置
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        request body model.CreateTableReq true "创建数据表配置请求参数"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables [post]
func (api *ConfigAPI) CreateTable(c *gin.Context) {
	var req model.CreateTableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	userID := uint(c.GetInt64("user_id"))
	appID := uint(c.GetInt64("app_id"))
	id, err := api.configService.CreateTable(&req, userID, appID)

	if err != nil {
		utils.ServerError(c, err)
		return
	}

	utils.Success(c, gin.H{"ID": id})
}

// @Summary      更新数据表基本信息
// @Description  更新指定数据表的基本信息
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        table_id path int true "配置ID"
// @Param        request body model.ConfigTable true "更新数据表配置请求参数"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{table_id} [put]
func (api *ConfigAPI) UpdateTable(c *gin.Context) {
	id := utils.ParseUint(c.Param("table_id"))
	if id == 0 {
		utils.Error(c, 400, "invalid id")
		return
	}

	var req model.ConfigTable
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}
	req.ID = uint(id)

	userID := uint(c.GetInt64("user_id"))
	appID := uint(c.GetInt64("app_id"))
	err := api.configService.UpdateTable(&req, userID, appID)

	if err != nil {
		utils.ServerError(c, err)
		return
	}

	utils.Success(c, nil)
}

// @Summary      更新数据表字段信息
// @Description  更新指定数据表的字段信息
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        table_id path int true "配置ID"
// @Param        request body []model.FieldUpdate true "更新数据表字段请求参数"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{table_id}/fields [put]
func (api *ConfigAPI) UpdateTableFields(c *gin.Context) {
	table_id := utils.ParseUint(c.Param("table_id"))
	if table_id == 0 {
		utils.Error(c, 400, "invalid id")
		return
	}

	var req []model.FieldUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	userID := uint(c.GetInt64("user_id"))
	appID := uint(c.GetInt64("app_id"))
	err := api.configService.UpdateTableFields(table_id, req, userID, appID)

	if err != nil {
		utils.ServerError(c, err)
		return
	}

	utils.Success(c, nil)
}

// @Summary      更新索引信息
// @Description  更新指定数据表的索引信息
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        table_id path int true "配置ID"
// @Param        request body []model.IndexUpdate true "更新数据表索引请求参数"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{table_id}/indexes [put]
func (api *ConfigAPI) UpdateTableIndexes(c *gin.Context) {
	table_id := utils.ParseUint(c.Param("table_id"))
	if table_id == 0 {
		utils.Error(c, 400, "invalid id")
		return
	}

	var req []model.IndexUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	userID := uint(c.GetInt64("user_id"))
	appID := uint(c.GetInt64("app_id"))
	err := api.configService.UpdateTableIndexes(table_id, req, userID, appID)

	if err != nil {
		utils.ServerError(c, err)
		return
	}

	utils.Success(c, nil)
}

// @Summary      更新 func 字段
// @Description  更新指定数据表的 func 字段
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        table_id path int true "配置ID"
// @Param        request body model.CreateTableReq true "更新数据表配置请求参数"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{table_id}/func [put]
func (api *ConfigAPI) UpdateTableFunc(c *gin.Context) {
	table_id := utils.ParseUint(c.Param("table_id"))
	if table_id == 0 {
		utils.Error(c, 400, "invalid id")
		return
	}

	var req model.CreateTableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}
	req.ID = uint(table_id)

	userID := uint(c.GetInt64("user_id"))
	appID := uint(c.GetInt64("app_id"))
	err := api.configService.UpdateTableFunc(&req, userID, appID)

	if err != nil {
		utils.ServerError(c, err)
		return
	}

	utils.Success(c, nil)
}

// @Summary      删除数据表配置
// @Description  删除指定的数据表配置
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        table_id path int true "配置ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{table_id} [delete]
func (api *ConfigAPI) DeleteTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("table_id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "invalid id")
		return
	}

	if err := api.configService.DeleteTable(uint(id)); err != nil {
		utils.ServerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
