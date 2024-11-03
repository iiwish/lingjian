package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service"
)

// ConfigAPI 配置相关API处理器
type ConfigAPI struct {
	configService *service.ConfigService
}

// NewConfigAPI 创建配置API处理器
func NewConfigAPI(configService *service.ConfigService) *ConfigAPI {
	return &ConfigAPI{
		configService: configService,
	}
}

// RegisterRoutes 注册路由
func (api *ConfigAPI) RegisterRoutes(router *gin.RouterGroup) {
	config := router.Group("/config")
	{
		// 数据表配置
		config.POST("/tables", api.CreateTable)
		config.PUT("/tables/:id", api.UpdateTable)
		config.GET("/tables", api.ListTables)
		config.GET("/tables/:id", api.GetTable)
		config.DELETE("/tables/:id", api.DeleteTable)
		config.GET("/tables/:id/versions", api.GetTableVersions)
		config.POST("/tables/:id/rollback", api.RollbackTable)

		// 维度配置
		config.POST("/dimensions", api.CreateDimension)
		config.PUT("/dimensions/:id", api.UpdateDimension)
		config.GET("/dimensions", api.ListDimensions)
		config.GET("/dimensions/:id", api.GetDimension)
		config.DELETE("/dimensions/:id", api.DeleteDimension)
		config.GET("/dimensions/:id/versions", api.GetDimensionVersions)
		config.POST("/dimensions/:id/rollback", api.RollbackDimension)

		// 数据模型配置
		config.POST("/models", api.CreateModel)
		config.PUT("/models/:id", api.UpdateModel)
		config.GET("/models", api.ListModels)
		config.GET("/models/:id", api.GetModel)
		config.DELETE("/models/:id", api.DeleteModel)
		config.GET("/models/:id/versions", api.GetModelVersions)
		config.POST("/models/:id/rollback", api.RollbackModel)

		// 表单配置
		config.POST("/forms", api.CreateForm)
		config.PUT("/forms/:id", api.UpdateForm)
		config.GET("/forms", api.ListForms)
		config.GET("/forms/:id", api.GetForm)
		config.DELETE("/forms/:id", api.DeleteForm)
		config.GET("/forms/:id/versions", api.GetFormVersions)
		config.POST("/forms/:id/rollback", api.RollbackForm)

		// 菜单配置
		config.POST("/menus", api.CreateMenu)
		config.PUT("/menus/:id", api.UpdateMenu)
		config.GET("/menus", api.ListMenus)
		config.GET("/menus/:id", api.GetMenu)
		config.DELETE("/menus/:id", api.DeleteMenu)
		config.GET("/menus/:id/versions", api.GetMenuVersions)
		config.POST("/menus/:id/rollback", api.RollbackMenu)
	}
}

// CreateTable 创建数据表配置
func (api *ConfigAPI) CreateTable(c *gin.Context) {
	var table model.ConfigTable
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := api.configService.CreateTable(&table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, table)
}

// UpdateTable 更新数据表配置
func (api *ConfigAPI) UpdateTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var table model.ConfigTable
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	table.ID = uint(id)

	if err := api.configService.UpdateTable(&table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

// ListTables 获取数据表配置列表
func (api *ConfigAPI) ListTables(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Query("app_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid app_id"})
		return
	}

	tables, err := api.configService.ListTables(uint(appID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tables)
}

// GetTable 获取数据表配置详情
func (api *ConfigAPI) GetTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	table, err := api.configService.GetTable(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

// DeleteTable 删除数据表配置
func (api *ConfigAPI) DeleteTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := api.configService.DeleteTable(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetTableVersions 获取数据表配置版本历史
func (api *ConfigAPI) GetTableVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	versions, err := api.configService.GetTableVersions(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// RollbackTable 回滚数据表配置到指定版本
func (api *ConfigAPI) RollbackTable(c *gin.Context) {
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

	if err := api.configService.RollbackTable(uint(id), version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// 以下是维度配置的API处理函数
// CreateDimension 创建维度配置
func (api *ConfigAPI) CreateDimension(c *gin.Context) {
	var dimension model.ConfigDimension
	if err := c.ShouldBindJSON(&dimension); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := api.configService.CreateDimension(&dimension); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dimension)
}

// UpdateDimension 更新维度配置
func (api *ConfigAPI) UpdateDimension(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var dimension model.ConfigDimension
	if err := c.ShouldBindJSON(&dimension); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dimension.ID = uint(id)

	if err := api.configService.UpdateDimension(&dimension); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dimension)
}

// ListDimensions 获取维度配置列表
func (api *ConfigAPI) ListDimensions(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Query("app_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid app_id"})
		return
	}

	dimensions, err := api.configService.ListDimensions(uint(appID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dimensions)
}

// GetDimension 获取维度配置详情
func (api *ConfigAPI) GetDimension(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	dimension, err := api.configService.GetDimension(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dimension)
}

// DeleteDimension 删除维度配置
func (api *ConfigAPI) DeleteDimension(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := api.configService.DeleteDimension(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetDimensionVersions 获取维度配置版本历史
func (api *ConfigAPI) GetDimensionVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	versions, err := api.configService.GetDimensionVersions(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// RollbackDimension 回滚维度配置到指定版本
func (api *ConfigAPI) RollbackDimension(c *gin.Context) {
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

	if err := api.configService.RollbackDimension(uint(id), version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// 以下是数据模型配置的API处理函数...
// CreateModel 创建数据模型配置
func (api *ConfigAPI) CreateModel(c *gin.Context) {
	var model model.ConfigDataModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := api.configService.CreateModel(&model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model)
}

// UpdateModel 更新数据模型配置
func (api *ConfigAPI) UpdateModel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var model model.ConfigDataModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	model.ID = uint(id)

	if err := api.configService.UpdateModel(&model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model)
}

// ListModels 获取数据模型配置列表
func (api *ConfigAPI) ListModels(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Query("app_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid app_id"})
		return
	}

	models, err := api.configService.ListModels(uint(appID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models)
}

// GetModel 获取数据模型配置详情
func (api *ConfigAPI) GetModel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	model, err := api.configService.GetModel(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model)
}

// DeleteModel 删除数据模型配置
func (api *ConfigAPI) DeleteModel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := api.configService.DeleteModel(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetModelVersions 获取数据模型配置版本历史
func (api *ConfigAPI) GetModelVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	versions, err := api.configService.GetModelVersions(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// RollbackModel 回滚数据模型配置到指定版本
func (api *ConfigAPI) RollbackModel(c *gin.Context) {
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

	if err := api.configService.RollbackModel(uint(id), version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// 以下是表单配置的API处理函数...
// CreateForm 创建表单配置
func (api *ConfigAPI) CreateForm(c *gin.Context) {
	var form model.ConfigForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := api.configService.CreateForm(&form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, form)
}

// UpdateForm 更新表单配置
func (api *ConfigAPI) UpdateForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var form model.ConfigForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	form.ID = uint(id)

	if err := api.configService.UpdateForm(&form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, form)
}

// ListForms 获取表单配置列表
func (api *ConfigAPI) ListForms(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Query("app_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid app_id"})
		return
	}

	forms, err := api.configService.ListForms(uint(appID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, forms)
}

// GetForm 获取表单配置详情
func (api *ConfigAPI) GetForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	form, err := api.configService.GetForm(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, form)
}

// DeleteForm 删除表单配置
func (api *ConfigAPI) DeleteForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := api.configService.DeleteForm(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetFormVersions 获取表单配置版本历史
func (api *ConfigAPI) GetFormVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	versions, err := api.configService.GetFormVersions(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// RollbackForm 回滚表单配置到指定版本
func (api *ConfigAPI) RollbackForm(c *gin.Context) {
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

	if err := api.configService.RollbackForm(uint(id), version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// 以下是菜单配置的API处理函数...
// CreateMenu 创建菜单配置
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

// UpdateMenu 更新菜单配置
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

// ListMenus 获取菜单配置列表
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

// GetMenu 获取菜单配置详情
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

// DeleteMenu 删除菜单配置
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

// GetMenuVersions 获取菜单配置版本历史
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

// RollbackMenu 回滚菜单配置到指定版本
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
