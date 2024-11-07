package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service"
)

// RegisterConfigRoutes 注册配置相关路由
func RegisterConfigRoutes(router *gin.RouterGroup) {
	configService := service.NewConfigService(model.DB)
	configAPI := NewConfigAPI(configService)
	configAPI.RegisterRoutes(router)
}

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

// @Summary      创建数据表配置
// @Description  创建新的数据表配置
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        request body service.CreateTableRequest true "创建数据表配置请求参数"
// @Success      201  {object}  gin.H
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/tables [post]
func (api *ConfigAPI) CreateTable(c *gin.Context) {
	var req service.CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文中获取当前用户ID
	userID := uint(c.GetInt64("user_id"))

	if err := api.configService.CreateTable(&req, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "数据表配置创建成功"})
}

// @Summary      更新数据表配置
// @Description  更新已存在的数据表配置
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        request body model.ConfigTable true "更新数据表配置请求参数"
// @Success      200  {object}  model.ConfigTable
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/tables/{id} [put]
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

// @Summary      获取数据表配置列表
// @Description  获取指定应用的所有数据表配置
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        app_id query int true "应用ID"
// @Success      200  {array}   model.ConfigTable
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/tables [get]
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

// @Summary      获取数据表配置详情
// @Description  获取指定数据表配置的详细信息
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {object}  model.ConfigTable
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/tables/{id} [get]
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

// @Summary      删除数据表配置
// @Description  删除指定的数据表配置
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/tables/{id} [delete]
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

// @Summary      获取数据表配置版本历史
// @Description  获取指定数据表配置的版本历史记录
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {array}   model.ConfigTableVersion
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/tables/{id}/versions [get]
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

// @Summary      回滚数据表配置
// @Description  将数据表配置回滚到指定版本
// @Tags         ConfigTable
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        version query int true "目标版本号"
// @Success      200  {object}  nil
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/tables/{id}/rollback [post]
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

// @Summary      创建维度配置
// @Description  创建新的维度配置
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        request body service.CreateDimensionRequest true "创建维度配置请求参数"
// @Success      201  {object}  gin.H
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/dimensions [post]
func (api *ConfigAPI) CreateDimension(c *gin.Context) {
	var req service.CreateDimensionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文中获取当前用户ID
	userID := uint(c.GetInt64("user_id"))

	if err := api.configService.CreateDimension(&req, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "维度配置创建成功"})
}

// @Summary      更新维度配置
// @Description  更新已存在的维度配置
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        request body model.ConfigDimension true "更新维度配置请求参数"
// @Success      200  {object}  model.ConfigDimension
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/dimensions/{id} [put]
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

// @Summary      获取维度配置列表
// @Description  获取指定应用的维度配置列表
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        app_id query int true "应用ID"
// @Success      200  {array}   model.ConfigDimension
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/dimensions [get]
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

// @Summary      获取维度配置详情
// @Description  获取指定维度配置的详细信息
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {object}  model.ConfigDimension
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/dimensions/{id} [get]
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

// @Summary      删除维度配置
// @Description  删除指定的维度配置
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/dimensions/{id} [delete]
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

// @Summary      获取维度配置版本历史
// @Description  获取指定维度配置的版本历史记录
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {array}   model.ConfigDimensionVersion
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/dimensions/{id}/versions [get]
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

// @Summary      回滚维度配置
// @Description  将维度配置回滚到指定版本
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        version query int true "目标版本号"
// @Success      200  {object}  nil
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/dimensions/{id}/rollback [post]
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

// @Summary      创建数据模型配置
// @Description  创建新的数据模型配置
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Param        request body model.ConfigDataModel true "创建数据模型配置请求参数"
// @Success      201  {object}  model.ConfigDataModel
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/models [post]
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

// @Summary      更新数据模型配置
// @Description  更新已存在的数据模型配置
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        request body model.ConfigDataModel true "更新数据模型配置请求参数"
// @Success      200  {object}  model.ConfigDataModel
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/models/{id} [put]
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

// @Summary      获取数据模型配置列表
// @Description  获取指定应用的数据模型配置列表
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Param        app_id query int true "应用ID"
// @Success      200  {array}   model.ConfigDataModel
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/models [get]
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

// @Summary      获取数据模型配置详情
// @Description  获取指定数据模型配置的详细信息
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {object}  model.ConfigDataModel
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/models/{id} [get]
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

// @Summary      删除数据模型配置
// @Description  删除指定的数据模型配置
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/models/{id} [delete]
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

// @Summary      获取数据模型配置版本历史
// @Description  获取指定数据模型配置的版本历史记录
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {array}   model.ConfigDataModelVersion
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/models/{id}/versions [get]
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

// @Summary      回滚数据模型配置
// @Description  将数据模型配置回滚到指定版本
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        version query int true "目标版本号"
// @Success      200  {object}  nil
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/models/{id}/rollback [post]
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

// @Summary      创建表单配置
// @Description  创建新的表单配置
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        request body model.ConfigForm true "创建表单配置请求参数"
// @Success      201  {object}  model.ConfigForm
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/forms [post]
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

// @Summary      更新表单配置
// @Description  更新已存在的表单配置
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        request body model.ConfigForm true "更新表单配置请求参数"
// @Success      200  {object}  model.ConfigForm
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/forms/{id} [put]
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

// @Summary      获取表单配置列表
// @Description  获取指定应用的表单配置列表
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        app_id query int true "应用ID"
// @Success      200  {array}   model.ConfigForm
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/forms [get]
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

// @Summary      获取表单配置详情
// @Description  获取指定表单配置的详细信息
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {object}  model.ConfigForm
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/forms/{id} [get]
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

// @Summary      删除表单配置
// @Description  删除指定的表单配置
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/forms/{id} [delete]
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

// @Summary      获取表单配置版本历史
// @Description  获取指定表单配置的版本历史记录
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {array}   model.ConfigFormVersion
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/forms/{id}/versions [get]
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

// @Summary      回滚表单配置
// @Description  将表单配置回滚到指定版本
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        version query int true "目标版本号"
// @Success      200  {object}  nil
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/forms/{id}/rollback [post]
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

// @Summary      创建菜单配置
// @Description  创建新的菜单配置
// @Tags         ConfigMenu
// @Accept       json
// @Produce      json
// @Param        request body model.ConfigMenu true "创建菜单配置请求参数"
// @Success      201  {object}  model.ConfigMenu
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
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
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
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
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
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
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
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
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
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
// @Success      200  {array}   model.ConfigMenuVersion
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
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
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
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
