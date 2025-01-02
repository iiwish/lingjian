package element

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service/element"
)

// ElementAPI 元素相关API处理器
type ElementAPI struct {
	elementService *element.ElementService
}

// NewElementAPI 创建元素API处理器
func NewElementAPI(elementService *element.ElementService) *ElementAPI {
	return &ElementAPI{
		elementService: elementService,
	}
}

// RegisterElementRoutes 注册元素相关路由
func RegisterElementRoutes(router *gin.RouterGroup) {
	elementService := element.NewElementService(model.DB)
	elementAPI := NewElementAPI(elementService)
	elementAPI.RegisterRoutes(router)
}

// RegisterRoutes 注册路由
func (api *ElementAPI) RegisterRoutes(router *gin.RouterGroup) {
	// 维度明细配置
	router.GET("/dimension/:dim_id", api.GetDimensionItems)
	router.POST("/dimension/:dim_id", api.CreateDimensionItem)
	router.PUT("/dimension/:dim_id", api.UpdateDimensionItem)
	router.PUT("/dimension/:dim_id/:id", api.UpdateDimensionItemSort)
	router.DELETE("/dimension/:dim_id", api.DeleteDimensionItems)

	// 菜单明细配置
	router.GET("/menu/:menu_id", api.GetMenuItems)
	router.POST("/menu/:menu_id", api.CreateMenuItem)
	router.PUT("/menu/:menu_id/:id", api.UpdateMenuItem)
	router.DELETE("/menu/:menu_id/:id", api.DeleteMenuItem)

	// 数据表明细配置
	router.POST("/table/:table_id/query", api.QueryTableItems)
	// router.GET("/table/:table_id", api.GetTableItems)
	router.POST("/table/:table_id", api.CreateTableItems)
	router.PUT("/table/:table_id", api.UpdateTableItems)
	router.DELETE("/table/:table_id", api.DeleteTableItems)
}
