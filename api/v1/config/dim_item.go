package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      创建维度配置项
// @Description  创建新的维度配置项
// @Tags         ConfigDimensionItem
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        id         path int                      true  "维度ID"
// @Param        dimension  body model.ConfigDimensionItem true  "创建维度配置项的请求参数"
// @Success      201        {object}  Response
// @Failure      400        {object}  Response
// @Failure      500        {object}  Response
// @Router       /config/dimensions/{dim_id}/items [post]
func (api *ConfigAPI) CreateDimensionItem(c *gin.Context) {
	// 获取id参数
	dimID := c.Param("dim_id")
	if dimID == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "dim_id不能为空"})
		return
	}

	// 获取请求参数
	var dimension model.ConfigDimensionItem
	if err := c.ShouldBindJSON(&dimension); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	// 校验请求参数
	if dimension.Code == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "code不能为空"})
		return
	}
	if dimension.Name == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "name不能为空"})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	id, err := api.configService.CreateDimensionItem(&dimension, userID, utils.ParseUint(dimID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ID": id})
}

// @Summary      批量创建维度配置项
// @Description  批量创建新的维度配置项
// @Tags         ConfigDimensionItem
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        id          path int                        true  "维度ID"
// @Param        dimensions  body []model.ConfigDimensionItem true  "批量创建维度配置项的请求参数"
// @Success      201         {object}  Response
// @Failure      400         {object}  Response
// @Failure      500         {object}  Response
// @Router       /config/dimensions/{dim_id}/items/batch [post]
func (api *ConfigAPI) BatchCreateDimensionItems(c *gin.Context) {
	// 获取id参数
	dimID := c.Param("dim_id")
	if dimID == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "dim_id不能为空"})
		return
	}

	// 获取请求参数
	var dimensions []*model.ConfigDimensionItem
	if err := c.ShouldBindJSON(&dimensions); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	// 校验请求参数
	for _, dimension := range dimensions {
		if dimension.Code == "" {
			c.JSON(http.StatusBadRequest, Response{Error: "code不能为空"})
			return
		}
		if dimension.Name == "" {
			c.JSON(http.StatusBadRequest, Response{Error: "name不能为空"})
			return
		}
	}

	userID := uint(c.GetInt64("user_id"))
	err := api.configService.BatchCreateDimensionItems(dimensions, userID, utils.ParseUint(dimID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// @Summary      更新维度配置项
// @Description  更新已存在的维度配置项
// @Tags         ConfigDimensionItem
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        dim_id     path int                      true  "维度ID"
// @Param        id         path int                      true  "配置项ID"
// @Param        dimension  body model.ConfigDimensionItem true  "更新维度配置项的请求参数"
// @Success      200        {object}  model.ConfigDimensionItem
// @Failure      400        {object}  Response
// @Failure      500        {object}  Response
// @Router       /config/dimensions/{dim_id}/items/{id} [put]
func (api *ConfigAPI) UpdateDimensionItem(c *gin.Context) {
	// 获取id参数
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "id不能为空"})
		return
	}

	dim_id := c.Param("dim_id")
	if dim_id == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "dim_id不能为空"})
		return
	}

	// 获取请求参数
	var dimension model.ConfigDimensionItem
	if err := c.ShouldBindJSON(&dimension); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}
	dimension.ID = utils.ParseUint(id)

	// 校验请求参数
	if dimension.Code == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "code不能为空"})
		return
	}
	if dimension.Name == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "name不能为空"})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	if err := api.configService.UpdateDimensionItem(&dimension, userID, utils.ParseUint(dim_id)); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dimension)
}

// @Summary      获取维度配置项树
// @Description  获取指定维度的配置项树
// @Tags         ConfigDimensionItem
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        id  path int true "维度ID"
// @Success      200 {array}   model.ConfigDimensionItem
// @Failure      400 {object}  Response
// @Failure      500 {object}  Response
// @Router       /config/dimensions/{dim_id}/items [get]
func (api *ConfigAPI) TreeDimensionItems(c *gin.Context) {
	// 获取id参数
	id := c.Param("dim_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "dim_id不能为空"})
		return
	}

	items, err := api.configService.TreeDimensionItems(utils.ParseUint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// @Summary      删除维度配置项
// @Description  删除指定的维度配置项
// @Tags         ConfigDimensionItem
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        dim_id  path int true "维度ID"
// @Param        id      path int true "配置项ID"
// @Success      204     {object}  nil
// @Failure      400     {object}  Response
// @Failure      500     {object}  Response
// @Router       /config/dimensions/{dim_id}/items/{id} [delete]
func (api *ConfigAPI) DeleteDimensionItem(c *gin.Context) {
	// 获取id参数
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "id不能为空"})
		return
	}

	dim_id := c.Param("dim_id")
	if dim_id == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "dim_id不能为空"})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	if err := api.configService.DeleteDimensionItem(utils.ParseUint(id), userID, utils.ParseUint(dim_id)); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary      批量删除维度配置项
// @Description  批量删除指定的维度配置项
// @Tags         ConfigDimensionItem
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        id   path int    true  "维度ID"
// @Param        ids  body []uint true  "配置项ID列表"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions/{dim_id}/items/batch [delete]
func (api *ConfigAPI) BatchDeleteDimensionItems(c *gin.Context) {
	// 获取dimID参数
	dimID := c.Param("dim_id")
	if dimID == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "dim_id不能为空"})
		return
	}

	// 获取请求参数
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	if err := api.configService.BatchDeleteDimensionItems(userID, utils.ParseUint(dimID), ids); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
