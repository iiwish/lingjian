package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/service"
	"github.com/iiwish/lingjian/pkg/utils"
)

// RegisterConfigRoutes 注册配置相关路由
func RegisterConfigRoutes(r *gin.RouterGroup) {
	config := r.Group("/config")
	{
		// 数据表配置
		tables := config.Group("/tables")
		{
			tables.POST("", CreateTable)
			tables.PUT("/:id", UpdateTable)
			tables.GET("/:id/versions", GetTableVersions)
			tables.POST("/:id/rollback/:version", RollbackTable)
		}

		// 维度配置
		dimensions := config.Group("/dimensions")
		{
			dimensions.POST("", CreateDimension)
			dimensions.GET("/:id/values", GetDimensionValues)
		}

		// TODO: 添加其他配置类型的路由
		// - 数据模型配置
		// - 表单配置
		// - 菜单配置
	}
}

// CreateTable 创建数据表配置
func CreateTable(c *gin.Context) {
	var req service.CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	creatorID, _ := c.Get("user_id")
	configService := &service.ConfigService{}
	if err := configService.CreateTable(&req, creatorID.(uint)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// CreateDimension 创建维度配置
func CreateDimension(c *gin.Context) {
	var req service.CreateDimensionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	creatorID, _ := c.Get("user_id")
	configService := &service.ConfigService{}
	if err := configService.CreateDimension(&req, creatorID.(uint)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// GetDimensionValues 获取维度值列表
func GetDimensionValues(c *gin.Context) {
	dimensionID := utils.ParseUint(c.Param("id"))
	var filter map[string]any
	if err := c.ShouldBindJSON(&filter); err != nil {
		filter = nil
	}

	configService := &service.ConfigService{}
	values, err := configService.GetDimensionValues(dimensionID, filter)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, values)
}

// ... 其他代码保持不变 ...
