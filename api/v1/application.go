package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/service"
	"github.com/iiwish/lingjian/pkg/utils"
)

// RegisterAppRoutes 注册应用相关路由
func RegisterAppRoutes(r *gin.RouterGroup) {
	app := r.Group("/apps")
	{
		app.POST("", CreateApp)
		app.POST("/:app_id/users/:user_id", AssignAppToUser)
		app.GET("/users/:user_id", GetUserApps)
		app.GET("/users/:user_id/default", GetDefaultApp)

		// 模板相关路由
		template := app.Group("/templates")
		{
			template.POST("", CreateTemplate)
			template.GET("", ListTemplates)
			template.POST("/:template_id/publish", PublishTemplate)
			template.POST("/:template_id/create", CreateFromTemplate)
		}
	}
}

type CreateAppRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
}

// CreateApp 创建应用
func CreateApp(c *gin.Context) {
	var req CreateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	appService := &service.AppService{}
	if err := appService.CreateApp(req.Name, req.Code, req.Description); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

type AssignAppRequest struct {
	IsDefault bool `json:"is_default"`
}

// AssignAppToUser 为用户分配应用
func AssignAppToUser(c *gin.Context) {
	var req AssignAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	userID := utils.ParseUint(c.Param("user_id"))
	appID := utils.ParseUint(c.Param("app_id"))

	appService := &service.AppService{}
	if err := appService.AssignAppToUser(userID, appID, req.IsDefault); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// GetUserApps 获取用户的所有应用
func GetUserApps(c *gin.Context) {
	userID := utils.ParseUint(c.Param("user_id"))

	appService := &service.AppService{}
	apps, err := appService.GetUserApps(userID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, apps)
}

// GetDefaultApp 获取用户的默认应用
func GetDefaultApp(c *gin.Context) {
	userID := utils.ParseUint(c.Param("user_id"))

	appService := &service.AppService{}
	app, err := appService.GetDefaultApp(userID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, app)
}

type CreateTemplateRequest struct {
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	Configuration string  `json:"configuration" binding:"required"`
	Price         float64 `json:"price"`
}

// CreateTemplate 创建应用模板
func CreateTemplate(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	creatorID, _ := c.Get("user_id")
	appService := &service.AppService{}
	if err := appService.CreateAppTemplate(
		req.Name,
		req.Description,
		req.Configuration,
		req.Price,
		creatorID.(uint),
	); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// ListTemplates 列出应用模板
func ListTemplates(c *gin.Context) {
	appService := &service.AppService{}
	templates, err := appService.ListAppTemplates(1) // 只列出已上架的模板
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, templates)
}

type CreateFromTemplateRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

// CreateFromTemplate 从模板创建应用
func CreateFromTemplate(c *gin.Context) {
	var req CreateFromTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	templateID := utils.ParseUint(c.Param("template_id"))
	userID, _ := c.Get("user_id")

	appService := &service.AppService{}
	if err := appService.CreateAppFromTemplate(templateID, userID.(uint), req.Name, req.Code); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// PublishTemplate 发布模板
func PublishTemplate(c *gin.Context) {
	templateID := utils.ParseUint(c.Param("template_id"))

	appService := &service.AppService{}
	if err := appService.PublishTemplate(templateID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}
