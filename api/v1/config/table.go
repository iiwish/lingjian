package config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
)

// @Summary      创建数据表配置
// @Description  创建新的数据表配置
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        request body config.CreateTableRequest true "创建数据表配置请求参数"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables [post]
func (api *ConfigAPI) CreateTable(c *gin.Context) {
	var req model.CreateTableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	id, err := api.configService.CreateTable(&req, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ID": id})
}

// @Summary      更新数据表配置
// @Description  更新已存在的数据表配置
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        request body model.ConfigTable true "更新数据表配置请求参数"
// @Success      200  {object}  model.ConfigTable
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{id} [put]
func (api *ConfigAPI) UpdateTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	var table model.ConfigTable
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}
	table.ID = uint(id)

	userID := uint(c.GetInt64("user_id"))
	if err := api.configService.UpdateTable(&table, userID); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

// @Summary      获取数据表配置列表
// @Description  获取指定应用的所有数据表配置
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        app_id query int true "应用ID"
// @Success      200  {array}   model.ConfigTable
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables [get]
func (api *ConfigAPI) ListTables(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Query("app_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid app_id"})
		return
	}

	tables, err := api.configService.ListTables(uint(appID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, tables)
}

// @Summary      获取数据表配置详情
// @Description  获取指定数据表配置的详细信息
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {object}  model.ConfigTable
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{id} [get]
func (api *ConfigAPI) GetTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	table, err := api.configService.GetTable(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

// @Summary      删除数据表配置
// @Description  删除指定的数据表配置
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{id} [delete]
func (api *ConfigAPI) DeleteTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	if err := api.configService.DeleteTable(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
