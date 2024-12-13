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
		// 维度主体配置
		config.POST("/dimensions/", api.CreateDimension)
		config.PUT("/dimensions/:dim_id", api.UpdateDimension)
		// config.GET("/dimensions", api.ListDimensions)
		config.GET("/dimensions/:dim_id", api.GetDimension)
		config.DELETE("/dimensions/:dim_id", api.DeleteDimension)

		// 表单配置
		config.POST("/forms", api.CreateForm)
		config.PUT("/forms/:id", api.UpdateForm)
		// config.GET("/forms", api.ListForms)
		config.GET("/forms/:id", api.GetForm)
		config.DELETE("/forms/:id", api.DeleteForm)

		// 菜单配置
		config.POST("/menus", api.CreateMenu)
		config.PUT("/menus/:id", api.UpdateMenu)
		config.GET("/menus", api.GetMenus)
		config.GET("/menus/:id", api.GetMenuByID)
		config.DELETE("/menus/:id", api.DeleteMenu)

		// 数据模型配置
		config.POST("/models", api.CreateModel)
		config.PUT("/models/:id", api.UpdateModel)
		// config.GET("/models", api.ListModels)
		config.GET("/models/:id", api.GetModel)
		config.DELETE("/models/:id", api.DeleteModel)

		// 数据表主体配置
		config.GET("/tables/:table_id", api.GetTable)
		config.POST("/tables", api.CreateTable)
		config.PUT("/tables/:table_id", api.UpdateTable)
		// config.GET("/tables", api.ListTables)
		config.DELETE("/tables/:table_id", api.DeleteTable)
	}
}
