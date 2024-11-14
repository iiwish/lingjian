package config

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service/config"
)

// ConfigAPI 配置相关API处理器
type ConfigAPI struct {
	configService *config.ConfigService
}

// NewConfigAPI 创建配置API处理器
func NewConfigAPI(configService *config.ConfigService) *ConfigAPI {
	return &ConfigAPI{
		configService: configService,
	}
}

// Response 通用响应结构
// @Description 通用HTTP响应结构
type Response struct {
	Error string `json:"error,omitempty"`
}

// RegisterConfigRoutes 注册配置相关路由
func RegisterConfigRoutes(router *gin.RouterGroup) {
	configService := config.NewConfigService(model.DB)
	configAPI := NewConfigAPI(configService)
	configAPI.RegisterRoutes(router)
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

		// 维度配置
		config.POST("/dimensions", api.CreateDimension)
		config.PUT("/dimensions/:id", api.UpdateDimension)
		config.GET("/dimensions", api.ListDimensions)
		config.GET("/dimensions/:id", api.GetDimension)
		config.DELETE("/dimensions/:id", api.DeleteDimension)

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
		config.GET("/menus/tree", api.MenusTree)
		config.GET("/menus/:id", api.GetMenu)
		config.DELETE("/menus/:id", api.DeleteMenu)
	}
}
