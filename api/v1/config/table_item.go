package config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
)

// @Summary      创建数据表记录
// @Description  创建新的数据表记录
// @Tags         ConfigTableItem
// @Accept       json
// @Produce      json
// @Param        table_id path int true "表ID"
// @Param        request body map[string]interface{} true "创建数据表记录请求参数"
// @Success      201  {object}  gin.H{"ID": uint}
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{table_id}/items [post]
func (api *ConfigAPI) CreateTableItem(c *gin.Context) {
	// 获取表ID
	tableID, err := strconv.ParseUint(c.Param("table_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid table_id"})
		return
	}

	// 绑定请求参数
	var tableItem map[string]interface{}
	if err := c.ShouldBindJSON(&tableItem); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	id, err := api.configService.CreateTableItem(tableItem, userID, uint(tableID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ID": id})
}

// @Summary      批量创建数据表记录
// @Description  批量创建新的数据表记录
// @Tags         ConfigTableItem
// @Accept       json
// @Produce      json
// @Param        table_id path int true "表ID"
// @Param        request body []map[string]interface{} true "批量创建数据表记录请求参数"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{table_id}/items/batch [post]
func (api *ConfigAPI) BatchCreateTableItems(c *gin.Context) {
	// 获取表ID
	tableID, err := strconv.ParseUint(c.Param("table_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid table_id"})
		return
	}

	// 绑定请求参数
	var tableItems []map[string]interface{}
	if err := c.ShouldBindJSON(&tableItems); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	err = api.configService.BatchCreateTableItems(tableItems, userID, uint(tableID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary      更新数据表记录
// @Description  更新已存在的数据表记录
// @Tags         ConfigTableItem
// @Accept       json
// @Produce      json
// @Param        table_id path int true "表ID"
// @Param        id path int true "记录ID"
// @Param        request body map[string]interface{} true "更新数据表记录请求参数"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{table_id}/items/{id} [put]
func (api *ConfigAPI) UpdateTableItem(c *gin.Context) {
	// 获取表ID和记录ID
	tableID, err := strconv.ParseUint(c.Param("table_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid table_id"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	// 绑定请求参数
	var tableItem map[string]interface{}
	if err := c.ShouldBindJSON(&tableItem); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	err = api.configService.UpdateTableItem(tableItem, userID, uint(tableID), uint(itemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary      获取数据表记录
// @Description  获取指定数据表记录的详细信息
// @Tags         ConfigTableItem
// @Accept       json
// @Produce      json
// @Param        table_id path int true "表ID"
// @Param        id path int true "记录ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{table_id}/items/{id} [get]
func (api *ConfigAPI) GetTableItem(c *gin.Context) {
	// 获取表ID和记录ID
	tableID, err := strconv.ParseUint(c.Param("table_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid table_id"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	// 获取记录
	item, err := api.configService.GetTableItem(uint(tableID), uint(itemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary      获取数据表记录列表
// @Description  获取指定数据表的记录列表
// @Tags         ConfigTableItem
// @Accept       json
// @Produce      json
// @Param        table_id path int true "表ID"
// @Param        page query int false "页码"
// @Param        page_size query int false "每页数量"
// @Param        query body model.QueryCondition false "查询条件"
// @Success      200  {object}  gin.H{"items": []map[string]interface{}, "total": int}
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{table_id}/items [get]
func (api *ConfigAPI) ListTableItems(c *gin.Context) {
	// 获取表ID
	tableID, err := strconv.ParseUint(c.Param("table_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid table_id"})
		return
	}

	// 获取分页参数
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	// 获取查询条件
	var query model.QueryCondition
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	// 获取记录列表
	items, total, err := api.configService.ListTableItems(uint(tableID), page, pageSize, &query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": total,
	})
}

// @Summary      删除数据表记录
// @Description  删除指定的数据表记录
// @Tags         ConfigTableItem
// @Accept       json
// @Produce      json
// @Param        table_id path int true "表ID"
// @Param        id path int true "记录ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{table_id}/items/{id} [delete]
func (api *ConfigAPI) DeleteTableItem(c *gin.Context) {
	// 获取表ID和记录ID
	tableID, err := strconv.ParseUint(c.Param("table_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid table_id"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	err = api.configService.DeleteTableItem(userID, uint(tableID), uint(itemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary      批量删除数据表记录
// @Description  批量删除指定的数据表记录
// @Tags         ConfigTableItem
// @Accept       json
// @Produce      json
// @Param        table_id path int true "表ID"
// @Param        request body []uint true "记录ID列表"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/tables/{table_id}/items/batch [delete]
func (api *ConfigAPI) BatchDeleteTableItems(c *gin.Context) {
	// 获取表ID
	tableID, err := strconv.ParseUint(c.Param("table_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid table_id"})
		return
	}

	// 绑定请求参数
	var itemIDs []uint
	if err := c.ShouldBindJSON(&itemIDs); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	err = api.configService.BatchDeleteTableItems(userID, uint(tableID), itemIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
