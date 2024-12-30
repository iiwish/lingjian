package config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      创建菜单配置
// @Description  创建新的菜单配置
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        request body model.ConfigMenu true "创建菜单配置请求参数"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menus [post]
func (api *ConfigAPI) CreateMenu(c *gin.Context) {
	var req model.ConfigMenu
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.MenuCode == "system" {
		utils.Error(c, http.StatusBadRequest, "不允许创建system菜单")
		return
	}

	userID := c.GetUint("user_id")
	appID := c.GetUint("app_id")
	req.AppID = appID
	menuID, err := api.configService.CreateMenu(&req, userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{"id": menuID})
}

// @Summary      更新菜单配置
// @Description  更新已存在的菜单配置
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        id path int true "配置ID"
// @Param        request body model.ConfigMenu true "更新菜单配置请求参数"
// @Success      200  {object}  model.ConfigMenu
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menus/{id} [put]
func (api *ConfigAPI) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	var menu model.ConfigMenu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}
	menu.ID = uint(id)

	if menu.MenuCode == "system" {
		c.JSON(http.StatusBadRequest, Response{Error: "不允许修改system菜单编码"})
		return
	}

	userID := c.GetUint("user_id")
	if err := api.configService.UpdateMenu(&menu, userID); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	utils.Success(c, menu)
}

// @Summary      获取菜单配置详情
// @Description  获取指定菜单配置的详细信息
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        id path int true "配置ID"
// @Success      200  {object}  model.ConfigMenu
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menus/{id} [get]
func (api *ConfigAPI) GetMenuByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	menu, err := api.configService.GetMenuByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	utils.Success(c, menu)
}

// @Summary      删除菜单配置
// @Description  删除指定的菜单配置
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        id path int true "配置ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menus/{id} [delete]
func (api *ConfigAPI) DeleteMenu(c *gin.Context) {
	id := utils.ParseUint(c.Param("id"))

	if err := api.configService.DeleteMenu(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	utils.Success(c, nil)
}

// @Summary      获取菜单列表
// @Description  根据可选的 level、parent_id 和 type 参数获取菜单列表
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        level     query    int     false  "菜单级别"
// @Param        parent_id query    uint    false  "父菜单ID"
// @Param        type      query    string  false  "菜单类型，可选值为 'children'、'descendants' , 默认为 'children'"
// @Success      200  {object}  []model.TreeConfigMenu
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menus [get]
func (api *ConfigAPI) GetMenus(c *gin.Context) {
	appID := c.GetUint("app_id")
	if appID == 0 {
		utils.Error(c, 400, "无效的 app_id 参数")
		return
	}

	operatorID := c.GetUint("user_id")
	if operatorID == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	// 获取可选参数
	levelStr := c.Query("level")
	parentIDStr := c.Query("parent_id")
	menuType := c.Query("type")

	if menuType == "" {
		menuType = "children"
	}

	var (
		level    *int
		parentID *uint
		err      error
	)

	if levelStr != "" {
		lvl := utils.ParseInt(levelStr)
		if levelStr != "0" && lvl == 0 {
			c.JSON(http.StatusBadRequest, Response{Error: "无效的 level 参数"})
			return
		}
		level = &lvl
	}

	if parentIDStr != "" {
		pid := utils.ParseUint(parentIDStr)
		if parentIDStr != "0" && pid == 0 {
			c.JSON(http.StatusBadRequest, Response{Error: "无效的 parent_id 参数"})
			return
		}
		parentID = &pid
	}

	if menuType != "children" && menuType != "descendants" {
		c.JSON(http.StatusBadRequest, Response{Error: "无效的 type 参数"})
		return
	}

	// 调用服务获取菜单列表
	menus, err := api.configService.GetMenus(appID, operatorID, level, parentID, menuType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	utils.Success(c, gin.H{"items": menus})
}

// @Summary      更新菜单项排序和父节点
// @Description  更新菜单项的排序和父节点
// @Tags         ConfigMenu
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
// @Router       /config/menu/{menu_id}/{id} [put]
func (api *ConfigAPI) UpdateMenuItemSort(c *gin.Context) {
	// 获取menuID参数
	menuID := c.Param("menu_id")
	if menuID == "" {
		utils.Error(c, http.StatusBadRequest, "menu_id不能为空")
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
		return
	}
	sort := c.Query("sort")
	if sort == "" {
		utils.Error(c, http.StatusBadRequest, "sort不能为空")
		return
	}

	userID := c.GetUint("user_id")
	if err := api.configService.UpdateMenuItemSort(userID, utils.ParseUint(menuID), utils.ParseUint(id), utils.ParseUint(parent), utils.ParseInt(sort)); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      获取系统菜单id
// @Description  获取系统菜单id
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Success      200  {object}  {id:1}
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menu/sysid [get]
func (api *ConfigAPI) GetSystemMenuID(c *gin.Context) {
	appid := c.GetUint("app_id")
	// 获取系统菜单id
	menuID, err := api.configService.GetSystemMenuID(appid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	utils.Success(c, gin.H{"id": menuID})
}
