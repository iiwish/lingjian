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
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	if _, err := api.configService.CreateMenu(&req, userID); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, Response{})
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

	userID := uint(c.GetInt64("user_id"))
	if err := api.configService.UpdateMenu(&menu, userID); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
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

	c.JSON(http.StatusOK, menu)
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

	c.Status(http.StatusNoContent)
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
// @Success      200  {object}  []model.ConfigMenu
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

	c.JSON(http.StatusOK, menus)
}
