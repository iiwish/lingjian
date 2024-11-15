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
// @Param        request body config.CreateMenuRequest true "创建菜单配置请求参数"
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

// @Summary      获取菜单配置列表
// @Description  获取指定应用的菜单配置列表
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Param        app_id query int true "应用ID"
// @Success      200  {array}   model.ConfigMenu
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menus [get]
func (api *ConfigAPI) ListMenus(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Query("app_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid app_id"})
		return
	}

	operatorID := c.GetUint("user_id")
	if operatorID == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	menus, err := api.configService.ListMenus(uint(appID), operatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, menus)
}

// @Summary      获取菜单配置详情
// @Description  获取指定菜单配置的详细信息
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {object}  model.ConfigMenu
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menus/{id} [get]
func (api *ConfigAPI) GetMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	menu, err := api.configService.GetMenu(uint(id))
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
// @Param        id path int true "配置ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menus/{id} [delete]
func (api *ConfigAPI) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	if err := api.configService.DeleteMenu(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// treeMenus
// @Summary      获取菜单树
// @Description  获取指定应用的菜单树
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Param        app_id query int true "应用ID"
// @Success      200  {object}  []model.TreeConfigMenu
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menus/tree [get]
func (api *ConfigAPI) TreeMenus(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Query("app_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid app_id"})
		return
	}

	operatorID := c.GetUint("user_id")
	if operatorID == 0 {
		utils.Error(c, 403, "未授权")
		return
	}

	menus, err := api.configService.TreeMenus(uint(appID), operatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, menus)
}
