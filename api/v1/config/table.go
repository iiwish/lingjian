package config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

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
	id, err := api.configService.CreateTable(&req, userID)

	if err != nil {
		utils.ServerError(c, err)
		return
	}

	utils.Success(c, gin.H{"ID": id})
}

// @Summary      更新数据表配置
// @Description  更新已存在的数据表配置
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        table_id path int true "配置ID"
// @Param        request body model.ConfigTable true "更新数据表配置请求参数"
// @Success      200  {object}  model.ConfigTable
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{table_id} [put]
func (api *ConfigAPI) UpdateTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("table_id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "invalid id")
		return
	}

	var table model.ConfigTable
	if err := c.ShouldBindJSON(&table); err != nil {
		utils.ValidationError(c, err)
		return
	}
	table.ID = uint(id)

	userID := uint(c.GetInt64("user_id"))
	if err := api.configService.UpdateTable(&table, userID); err != nil {
		utils.ServerError(c, err)
		return
	}

	utils.Success(c, table)
}

// @Summary      获取数据表配置列表
// @Description  获取指定应用的所有数据表配置
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Success      200  {array}   model.ConfigTable
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables [get]
func (api *ConfigAPI) ListTables(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Query("app_id"), 10, 64)
	if err != nil {
		utils.Error(c, 400, "invalid app_id")
		return
	}

	tables, err := api.configService.ListTables(uint(appID))
	if err != nil {
		utils.ServerError(c, err)
		return
	}

	utils.Success(c, tables)
}

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
	id, err := strconv.ParseUint(c.Param("table_id"), 10, 64)
	if err != nil {
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
