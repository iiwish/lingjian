package config

import (
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
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables [post]
func (api *ConfigAPI) CreateTable(c *gin.Context) {
	var req model.CreateTableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	userID := c.GetUint("user_id")
	appID := c.GetUint("app_id")
	id, err := api.configService.CreateTable(&req, userID, appID)

	if err != nil {
		utils.ServerError(c, err)
		return
	}

	utils.Success(c, gin.H{"id": id})
}

// @Summary      更新数据表配置
// @Description  更新指定数据表的配置信息,包括基本信息、字段信息、索引信息和功能信息
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        table_id path int true "配置ID"
// @Param        request body model.TableUpdateReq true "更新数据表配置请求参数"
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

	var req model.TableUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	userID := c.GetUint("user_id")
	appID := c.GetUint("app_id")
	err := api.configService.UpdateTable(id, &req, userID, appID)

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

	utils.Success(c, nil)
}
