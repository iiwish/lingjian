package element

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      获取维度项树
// @Description  获取指定维度的项树
// @Tags         Dimension
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        dim_id  path int true "维度ID"
// @Param        id 	query int     false "节点ID，不指定则返回整个维度配置项树"
// @Param 	 	 type      query    string  false  "菜单类型，可选值为 'children'、'descendants'、'leaves' , 默认为 'descendants'"
// @Param 	  	 level     query    int     false  "树的层级，可选值为 0、1、2、3， 默认为 0不指定层级"
// @Success      200 {array}   []model.TreeDimensionItem
// @Failure      400 {object}  utils.Response
// @Failure      500 {object}  utils.Response
// @Router       /dimension/{dim_id} [get]
func (api *ElementAPI) GetDimensionItems(c *gin.Context) {
	// 获取id参数
	id := c.Param("dim_id")
	if id == "" {
		utils.Error(c, http.StatusBadRequest, "dim_id不能为空")
		return
	}
	// 获取请求参数
	node_id := c.Query("id")
	if node_id == "" {
		node_id = "0"
	}

	// 获取query参数
	query_level := c.Query("level")

	// 获取type参数,默认为descendants
	query_type := c.Query("type")
	if query_type == "" {
		query_type = "descendants"
	}

	userID := c.GetUint("user_id")
	items, err := api.elementService.TreeDimensionItems(userID, utils.ParseUint(id), utils.ParseUint(node_id), query_type, utils.ParseUint(query_level))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, items)
}

// @Summary      批量创建维度项
// @Description  批量创建新的维度项
// @Tags         Dimension
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        id          path int                        true  "维度ID"
// @Param        dimensions  body model.CreateDimensionItemReq true  "创建维度项的请求参数"
// @Success      201         {object}  utils.Response
// @Failure      400         {object}  utils.Response
// @Failure      500         {object}  utils.Response
// @Router       /dimension/{dim_id} [post]
func (api *ElementAPI) CreateDimensionItem(c *gin.Context) {
	// 获取id参数
	dimID := c.Param("dim_id")
	if dimID == "" {
		utils.Error(c, http.StatusBadRequest, "dim_id不能为空")
		return
	}

	// 获取请求参数
	var dimensionItem *model.CreateDimensionItemReq
	if err := c.ShouldBindJSON(&dimensionItem); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 校验请求参数
	if dimensionItem.Code == "" {
		utils.Error(c, http.StatusBadRequest, "code不能为空")
		return
	}
	if dimensionItem.Name == "" {
		utils.Error(c, http.StatusBadRequest, "name不能为空")
		return
	}

	userID := c.GetUint("user_id")
	ids, err := api.elementService.CreateDimensionItem(dimensionItem, userID, utils.ParseUint(dimID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{"ids": ids})
}

// @Summary      更新维度项
// @Description  更新已存在的维度项
// @Tags         Dimension
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        dim_id     path int                      true  "维度ID"
// @Param        id         path int                      true  "配置项ID"
// @Param        dimension  body model.UpdateDimensionItemReq true  "更新维度项的请求参数"
// @Success      200        {object}  utils.Response
// @Failure      400        {object}  utils.Response
// @Failure      500        {object}  utils.Response
// @Router       /dimension/{dim_id} [put]
func (api *ElementAPI) UpdateDimensionItem(c *gin.Context) {
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
	var dimItem *model.UpdateDimensionItemReq
	if err := c.ShouldBindJSON(&dimItem); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	dimItem.ID = utils.ParseUint(id)

	// 校验请求参数
	if dimItem.Code == "" {
		utils.Error(c, http.StatusBadRequest, "code不能为空")
		return
	}
	if dimItem.Name == "" {
		utils.Error(c, http.StatusBadRequest, "name不能为空")
		return
	}

	userID := c.GetUint("user_id")
	if err := api.elementService.UpdateDimensionItem(dimItem, userID, utils.ParseUint(dim_id)); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      批量删除维度项
// @Description  批量删除指定的维度项
// @Tags         Dimension
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        id   path int    true  "维度ID"
// @Param        ids  body []uint true  "配置项ID列表"
// @Success      204  {object}  nil
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /dimension/{dim_id} [delete]
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

	userID := c.GetUint("user_id")
	if err := api.elementService.DeleteDimensionItems(userID, utils.ParseUint(dimID), ids); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary      更新维度项排序和父节点
// @Description  更新维度项的排序和父节点
// @Tags         Dimension
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        dim_id  path int true "维度ID"
// @Param        id      path int true "配置项ID"
// @Param        parent  query int false "父节点ID"
// @Param        sort    query int false "排序值"
// @Success      200  {object}  nil	"成功"
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /dimension/{dim_id}/{id} [put]
func (api *ElementAPI) UpdateDimensionItemSort(c *gin.Context) {
	// 获取dimID参数
	dimID := c.Param("dim_id")
	if dimID == "" {
		utils.Error(c, http.StatusBadRequest, "dim_id不能为空")
		return
	}

	// 获取id参数
	id := c.Param("id")
	if id == "" {
		utils.Error(c, http.StatusBadRequest, "id不能为空")
		return
	}

	// 获取请求参数
	parent := c.Query("parent")
	if parent == "" {
		utils.Error(c, http.StatusBadRequest, "parent不能为空")
	}
	sort := c.Query("sort")
	if sort == "" {
		utils.Error(c, http.StatusBadRequest, "sort不能为空")
	}

	userID := c.GetUint("user_id")
	if err := api.elementService.UpdateDimensionItemSort(userID, utils.ParseUint(dimID), utils.ParseUint(id), utils.ParseUint(parent), utils.ParseInt(sort)); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}
