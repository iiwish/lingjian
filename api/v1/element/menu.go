package element

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      获取菜单明细树
// @Description  获取菜单明细树
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer Token"
// @Param        App-ID header string true "应用ID"
// @Param        menu_id  path int true "菜单ID"
// @Param        parent_id 	query int     false "节点ID，不指定则返回整个维度配置项树"
// @Param 	 	 type      query    string  false  "菜单类型，可选值为 'children'、'descendants'、'leaves' , 默认为 'descendants'"
// @Param 	  	 level     query    int     false  "树的层级，可选值为 0、1、2、3， 默认为 0不指定层级"
// @Success      200 {array}   []model.TreeMenuItem
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /menu/{menu_id} [get]
func (api *ElementAPI) GetMenuItems(c *gin.Context) {
	// 获取id参数
	menuID := utils.ParseUint(c.Param("menu_id"))
	if menuID == 0 {
		utils.Error(c, http.StatusBadRequest, "menu_id不能为空")
		return
	}

	// 获取请求参数
	nodeID := utils.ParseUint(c.Query("parent_id"))

	// 获取query参数
	queryLevel := utils.ParseUint(c.Query("level"))

	// 获取type参数,默认为descendants
	queryType := c.Query("type")
	if queryType == "" {
		queryType = "descendants"
	}

	userID := c.GetUint("user_id")
	menus, err := api.elementService.GetMenuItems(userID, menuID, nodeID, queryType, queryLevel)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, menus)
}

// @Summary      创建菜单明细
// @Description  创建新的菜单明细
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer Token"
// @Param        App-ID header string true "应用ID"
// @Param        request body model.CreateMenuItemReq true "创建菜单明细请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /menu/{menu_id} [post]
func (api *ElementAPI) CreateMenuItem(c *gin.Context) {
	// 获取菜单ID
	menuID := utils.ParseUint(c.Param("menu_id"))
	if menuID == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid menu_id")
		return
	}

	// 绑定请求参数
	var req model.CreateMenuItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	itemID, err := api.elementService.CreateMenuItem(userID, &req, menuID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{"id": itemID})
}

// @Summary      更新菜单明细
// @Description  更新已存在的菜单明细
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer Token"
// @Param        App-ID header string true "应用ID"
// @Param        menu_id path int true "菜单ID"
// @Param        id path int true "菜单明细ID"
// @Param        request body model.UpdateMenuItemReq true "更新菜单明细请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /menu/{menu_id}/{id} [put]
func (api *ElementAPI) UpdateMenuItem(c *gin.Context) {
	// 获取菜单ID和菜单明细ID
	menuID := utils.ParseUint(c.Param("menu_id"))
	if menuID == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid menu_id")
		return
	}

	itemID := utils.ParseUint(c.Param("id"))
	if itemID == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	// 绑定请求参数
	var req model.UpdateMenuItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	req.ID = itemID
	userID := c.GetUint("user_id")
	err := api.elementService.UpdateMenuItem(&req, userID, menuID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      删除菜单明细
// @Description  删除指定的菜单明细
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer Token"
// @Param        App-ID header string true "应用ID"
// @Param        menu_id path int true "菜单ID"
// @Param        id path int true "菜单明细ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /menu/{menu_id}/{id} [delete]
func (api *ElementAPI) DeleteMenuItem(c *gin.Context) {
	// 获取菜单ID和菜单明细ID
	menuID := utils.ParseUint(c.Param("menu_id"))
	if menuID == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid menu_id")
		return
	}

	itemID := utils.ParseUint(c.Param("id"))
	if itemID == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	userID := c.GetUint("user_id")
	err := api.elementService.DeleteMenuItem(userID, menuID, []uint{itemID})
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      更新菜单项排序和父节点
// @Description  更新菜单项的排序和父节点
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        menu_id  path int true "菜单ID"
// @Param        id      path int true "菜单项ID"
// @Param        parent  query int false "父节点ID"
// @Param        sort    query int false "排序值"
// @Success      200  {object}  nil	"成功"
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /menu/{menu_id}/{id}/sort [put]
func (api *ElementAPI) UpdateMenuItemSort(c *gin.Context) {
	// 获取menuID参数
	menuID := utils.ParseUint(c.Param("menu_id"))
	if menuID == 0 {
		utils.Error(c, http.StatusBadRequest, "menu_id不能为空")
		return
	}

	// 获取id参数
	id := utils.ParseUint(c.Param("id"))
	if id == 0 {
		utils.Error(c, http.StatusBadRequest, "id不能为空")
		return
	}

	// 获取请求参数
	parent := utils.ParseUint(c.Query("parent"))

	sort := utils.ParseInt(c.Query("sort"))
	if sort == 0 {
		utils.Error(c, http.StatusBadRequest, "sort不能为空")
	}

	userID := c.GetUint("user_id")
	if err := api.elementService.UpdateDimensionItemSort(userID, menuID, id, parent, sort); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      创建系统菜单
// @Description  创建系统菜单
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        request body model.CreateMenuItemReq true "创建系统菜单请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /menu [post]
func (api *ElementAPI) CreateSysMenuItem(c *gin.Context) {

	// 绑定请求参数
	var req model.CreateMenuItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	appID := c.GetUint("app_id")
	if err := api.elementService.CreateSysMenu(appID, userID, &req); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}
