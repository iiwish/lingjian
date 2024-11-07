package config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
)

// @Summary      创建菜单配置
// @Description  创建新的菜单配置
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Param        request body model.ConfigMenu true "创建菜单配置请求参数"
// @Success      201  {object}  model.ConfigMenu
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menus [post]
func (api *ConfigAPI) CreateMenu(c *gin.Context) {
	var menu model.ConfigMenu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := api.configService.CreateMenu(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, menu)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var menu model.ConfigMenu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	menu.ID = uint(id)

	if err := api.configService.UpdateMenu(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid app_id"})
		return
	}

	menus, err := api.configService.ListMenus(uint(appID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	menu, err := api.configService.GetMenu(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := api.configService.DeleteMenu(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary      获取菜单配置版本历史
// @Description  获取指定菜单配置的版本历史记录
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {array}   model.ConfigVersion
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menus/{id}/versions [get]
func (api *ConfigAPI) GetMenuVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	versions, err := api.configService.GetMenuVersions(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// @Summary      回滚菜单配置
// @Description  将菜单配置回滚到指定版本
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        version query int true "目标版本号"
// @Success      200  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/menus/{id}/rollback [post]
func (api *ConfigAPI) RollbackMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	version, err := strconv.Atoi(c.Query("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version"})
		return
	}

	if err := api.configService.RollbackMenu(uint(id), version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
