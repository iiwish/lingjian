package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service"
	"github.com/iiwish/lingjian/pkg/utils"
)

// RegisterAppRoutes 注册应用相关路由
func RegisterAppRoutes(r *gin.RouterGroup) {
	app := r.Group("/apps")
	{
		app.GET("", ListApps)   // 获取应用列表
		app.POST("", CreateApp) // 创建应用

		// 模板相关路由
		// template := app.Group("/templates")
		// {
		// 	template.POST("", CreateTemplate)
		// 	template.GET("", ListTemplates)
		// 	template.POST("/:template_id/publish", PublishTemplate)
		// 	template.POST("/:template_id/create", CreateFromTemplate)
		// }
	}
}

// @Summary      获取应用列表
// @Description  获取所有应用列表
// @Tags         Application
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /apps [get]
func ListApps(c *gin.Context) {

	userID := c.GetUint("user_id")
	appService := &service.AppService{}
	result, err := appService.ListApps(userID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, result)
}

// @Summary      创建应用
// @Description  创建新的应用
// @Tags         Application
// @Accept       json
// @Produce      json
// @Param        request body CreateAppRequest true "创建应用请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /apps [post]
func CreateApp(c *gin.Context) {
	var req model.App
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	// 获取当前用户ID
	userID := c.GetUint("user_id")

	appService := &service.AppService{}
	result, err := appService.CreateApp(&req, userID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, result)
}
