package config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service"
)

// @Summary      创建数据表配置
// @Description  创建新的数据表配置
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        request body service.CreateTableRequest true "创建数据表配置请求参数"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables [post]
func (api *ConfigAPI) CreateTable(c *gin.Context) {
	var req service.CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := uint(c.GetInt64("user_id"))

	if err := api.configService.CreateTable(&req, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "数据表配置创建成功"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var table model.ConfigTable
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	table.ID = uint(id)

	if err := api.configService.UpdateTable(&table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid app_id"})
		return
	}

	tables, err := api.configService.ListTables(uint(appID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	table, err := api.configService.GetTable(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := api.configService.DeleteTable(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary      获取数据表配置版本历史
// @Description  获取指定数据表配置的版本历史记录
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {array}   model.ConfigVersion
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{id}/versions [get]
func (api *ConfigAPI) GetTableVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	versions, err := api.configService.GetTableVersions(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// @Summary      回滚数据表配置
// @Description  将数据表配置回滚到指定版本
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        version query int true "目标版本号"
// @Success      200  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{id}/rollback [post]
func (api *ConfigAPI) RollbackTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	version, err := strconv.Atoi(c.Query("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version"})
		return
	}

	if err := api.configService.RollbackTable(uint(id), version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
