package element

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      获取数据表记录列表
// @Description  获取指定数据表的记录列表
// @Tags         Table
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        table_id path int true "表ID"
// @Param        page query int false "页码"
// @Param        page_size query int false "每页数量"
// @Param        query body model.QueryCondition false "查询条件"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /table/{table_id} [get]
func (api *ElementAPI) GetTableItems(c *gin.Context) {
	// 获取表ID
	tableID := utils.ParseUint(c.Param("table_id"))
	if tableID == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid table_id")
		return
	}

	// 获取分页参数
	page := utils.ParseInt(c.Query("page"))
	if page <= 0 {
		page = 1
	}

	// 获取每页数量
	pageSize := utils.ParseInt(c.Query("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 1000 {
		pageSize = 1000
	}

	// 获取查询条件
	var query model.QueryCondition
	if err := c.ShouldBindJSON(&query); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 获取记录列表
	items, total, err := api.elementService.GetTableItems(tableID, page, pageSize, &query)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	data := map[string]interface{}{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"items":    items,
	}

	utils.Success(c, data)
}

// @Summary      创建数据表记录
// @Description  创建新的数据表记录
// @Tags         Table
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        table_id path int true "表ID"
// @Param        request body []map[string]interface{} true "创建数据表记录请求参数"
// @Success      201  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /table/{table_id} [post]
func (api *ElementAPI) CreateTableItems(c *gin.Context) {
	// 获取表ID
	tableID := utils.ParseUint(c.Param("table_id"))
	if tableID == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid table_id")
		return
	}

	// 绑定请求参数
	var tableItems []map[string]interface{}
	if err := c.ShouldBindJSON(&tableItems); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := uint(c.GetInt64("user_id"))
	err := api.elementService.CreateTableItems(tableItems, userID, tableID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      更新数据表记录
// @Description  更新已存在的数据表记录
// @Tags         Table
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        table_id path int true "表ID"
// @Param        id path int true "记录ID"
// @Param        request body model.UpdateTableItemsRequest true "更新数据表记录请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /table/{table_id} [put]
func (api *ElementAPI) UpdateTableItems(c *gin.Context) {
	// 获取表ID和记录ID
	tableID := utils.ParseUint(c.Param("table_id"))
	if tableID == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid table_id")
		return
	}

	// 绑定请求参数
	var reqItems model.UpdateTableItemsRequest
	if err := c.ShouldBindJSON(&reqItems); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := uint(c.GetInt64("user_id"))
	err := api.elementService.UpdateTableItems(reqItems, userID, uint(tableID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      批量删除数据表记录
// @Description  批量删除指定的数据表记录
// @Tags         Table
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        table_id path int true "表ID"
// @Param        request body []map[string]interface{} true "记录删除请求参数"
// @Success      204  {object}  nil
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /table/{table_id} [delete]
func (api *ElementAPI) DeleteTableItems(c *gin.Context) {
	// 获取表ID
	tableID := utils.ParseUint(c.Param("table_id"))
	if tableID == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid table_id")
		return
	}

	// 绑定请求参数
	var req []map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := uint(c.GetInt64("user_id"))
	err := api.elementService.DeleteTableItems(userID, tableID, req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
