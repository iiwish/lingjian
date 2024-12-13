package element

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      获取维度配置项树
// @Description  获取指定维度的配置项树
// @Tags         DimensionItem
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        id  path int true "维度ID"
// @Success      200 {array}   model.TreeConfigDimensionItem
// @Failure      400 {object}  utils.Response
// @Failure      500 {object}  utils.Response
// @Router       /dimensions/{dim_id} [get]
func (api *ElementAPI) TreeDimensionItems(c *gin.Context) {
	// 获取id参数
	id := c.Param("dim_id")
	if id == "" {
		utils.Error(c, http.StatusBadRequest, "dim_id不能为空")
		return
	}

	items, err := api.elementService.TreeDimensionItems(utils.ParseUint(id))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, items)
}

// @Summary      批量创建维度配置项
// @Description  批量创建新的维度配置项
// @Tags         DimensionItem
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        id          path int                        true  "维度ID"
// @Param        dimensions  body []model.ConfigDimensionItem true  "批量创建维度配置项的请求参数"
// @Success      201         {object}  utils.Response
// @Failure      400         {object}  utils.Response
// @Failure      500         {object}  utils.Response
// @Router       /dimensions/{dim_id} [post]
func (api *ElementAPI) CreateDimensionItems(c *gin.Context) {
	// 获取id参数
	dimID := c.Param("dim_id")
	if dimID == "" {
		utils.Error(c, http.StatusBadRequest, "dim_id不能为空")
		return
	}

	// 获取请求参数
	var dimensions []*model.ConfigDimensionItem
	if err := c.ShouldBindJSON(&dimensions); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 校验请求参数
	for _, dimension := range dimensions {
		if dimension.Code == "" {
			utils.Error(c, http.StatusBadRequest, "code不能为空")
			return
		}
		if dimension.Name == "" {
			utils.Error(c, http.StatusBadRequest, "name不能为空")
			return
		}
	}

	userID := uint(c.GetInt64("user_id"))
	err := api.elementService.BatchCreateDimensionItems(dimensions, userID, utils.ParseUint(dimID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{})
}

// @Summary      更新维度配置项
// @Description  更新已存在的维度配置项
// @Tags         DimensionItem
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        dim_id     path int                      true  "维度ID"
// @Param        id         path int                      true  "配置项ID"
// @Param        dimension  body model.ConfigDimensionItem true  "更新维度配置项的请求参数"
// @Success      200        {object}  model.ConfigDimensionItem
// @Failure      400        {object}  utils.Response
// @Failure      500        {object}  utils.Response
// @Router       /dimensions/{dim_id} [put]
func (api *ElementAPI) UpdateDimensionItems(c *gin.Context) {
	// 获取id参数
	id := c.Param("id")
	if id == "" {
		utils.Error(c, http.StatusBadRequest, "id不能为空")
		return
	}

	dim_id := c.Param("dim_id")
	if dim_id == "" {
		utils.Error(c, http.StatusBadRequest, "dim_id不能为空")
		return
	}

	// 获取请求参数
	var dimension model.ConfigDimensionItem
	if err := c.ShouldBindJSON(&dimension); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	dimension.ID = utils.ParseUint(id)

	// 校验请求参数
	if dimension.Code == "" {
		utils.Error(c, http.StatusBadRequest, "code不能为空")
		return
	}
	if dimension.Name == "" {
		utils.Error(c, http.StatusBadRequest, "name不能为空")
		return
	}

	userID := uint(c.GetInt64("user_id"))
	if err := api.elementService.UpdateDimensionItem(&dimension, userID, utils.ParseUint(dim_id)); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, dimension)
}

// @Summary      批量删除维度配置项
// @Description  批量删除指定的维度配置项
// @Tags         DimensionItem
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        id   path int    true  "维度ID"
// @Param        ids  body []uint true  "配置项ID列表"
// @Success      204  {object}  nil
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /dimensions/{dim_id} [delete]
func (api *ElementAPI) DeleteDimensionItems(c *gin.Context) {
	// 获取dimID参数
	dimID := c.Param("dim_id")
	if dimID == "" {
		utils.Error(c, http.StatusBadRequest, "dim_id不能为空")
		return
	}

	// 获取请求参数
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := uint(c.GetInt64("user_id"))
	if err := api.elementService.BatchDeleteDimensionItems(userID, utils.ParseUint(dimID), ids); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
