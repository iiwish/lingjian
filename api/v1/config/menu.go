package config

import (
	"net/http"

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
// @Param        request body model.CreateMenuReq true "创建菜单配置请求参数"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menus [post]
func (api *ConfigAPI) CreateMenu(c *gin.Context) {
	var req model.CreateMenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if utils.IsSystemMenu(req.TableName) {
		utils.Error(c, http.StatusBadRequest, "表名不允许使用，请更换")
		return
	}

	userID := c.GetUint("user_id")
	appID := c.GetUint("app_id")
	menuID, err := api.configService.CreateMenu(userID, appID, &req)
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
	dimID := utils.ParseUint(c.Param("id"))
	if dimID == 0 {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	var menu model.UpdateMenuReq
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if utils.IsSystemMenu(menu.TableName) {
		c.JSON(http.StatusBadRequest, Response{Error: "不允许修改system菜单编码"})
		return
	}

	userID := c.GetUint("user_id")
	if err := api.configService.UpdateMenu(&menu, userID, dimID); err != nil {
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

// @Summary      获取系统菜单id
// @Description  获取系统菜单id
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Success      200  {object}  Response
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

// @Summary      获取菜单配置列表
// @Description  获取菜单配置列表
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menus [get]
func (api *ConfigAPI) GetMenuList(c *gin.Context) {
	appID := c.GetUint("app_id")
	userID := c.GetUint("user_id")
	menus, err := api.configService.GetMenuList(userID, appID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, menus)
}

// @Summary      获取菜单配置
// @Description  获取指定的菜单配置
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
	id := utils.ParseUint(c.Param("id"))
	if id == 0 {
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
