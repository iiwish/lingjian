package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      创建维度配置
// @Description  创建新的维度配置
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        request body config.CreateDimensionRequest true "创建维度配置请求参数"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions/:id/items [post]
func (api *ConfigAPI) CreateDimensionItem(c *gin.Context) {
	// 获取id参数
	dimID := c.Param("id")
	if dimID == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "id不能为空"})
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

func (api *ConfigAPI) BatchCreateDimensionItems(c *gin.Context) {
	// 获取id参数
	dimID := c.Param("id")
	if dimID == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "id不能为空"})
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

func (api *ConfigAPI) TreeDimensionItems(c *gin.Context) {
	// 获取id参数
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "id不能为空"})
		return
	}

	items, err := api.configService.TreeDimensionItems(utils.ParseUint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

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

func (api *ConfigAPI) BatchDeleteDimensionItems(c *gin.Context) {
	// 获取dimID参数
	dimID := c.Param("id")
	if dimID == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "id不能为空"})
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
