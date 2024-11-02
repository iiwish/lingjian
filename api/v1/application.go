package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/service"
	"github.com/iiwish/lingjian/pkg/utils"
)

// RegisterApplicationRoutes 注册应用相关路由
func RegisterApplicationRoutes(r *gin.RouterGroup) {
	app := r.Group("/applications")
	{
		app.POST("", CreateApplication)
		app.POST("/:app_id/users/:user_id", AssignApplicationToUser)
		app.GET("/users/:user_id", GetUserApplications)
		app.GET("/users/:user_id/default", GetDefaultApplication)

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

type CreateApplicationRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
}

// CreateApplication 创建应用
func CreateApplication(c *gin.Context) {
	var req CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	appService := &service.ApplicationService{}
	if err := appService.CreateApplication(req.Name, req.Code, req.Description); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

type AssignApplicationRequest struct {
	IsDefault bool `json:"is_default"`
}

// AssignApplicationToUser 为用户分配应用
func AssignApplicationToUser(c *gin.Context) {
	var req AssignApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	userID := utils.ParseUint(c.Param("user_id"))
	appID := utils.ParseUint(c.Param("app_id"))

	appService := &service.ApplicationService{}
	if err := appService.AssignApplicationToUser(userID, appID, req.IsDefault); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// GetUserApplications 获取用户的所有应用
func GetUserApplications(c *gin.Context) {
	userID := utils.ParseUint(c.Param("user_id"))

	appService := &service.ApplicationService{}
	apps, err := appService.GetUserApplications(userID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, apps)
}

// GetDefaultApplication 获取用户的默认应用
func GetDefaultApplication(c *gin.Context) {
	userID := utils.ParseUint(c.Param("user_id"))

	appService := &service.ApplicationService{}
	app, err := appService.GetDefaultApplication(userID)
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
	appService := &service.ApplicationService{}
	if err := appService.CreateApplicationTemplate(
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
	appService := &service.ApplicationService{}
	templates, err := appService.ListApplicationTemplates(1) // 只列出已上架的模板
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

	appService := &service.ApplicationService{}
	if err := appService.CreateApplicationFromTemplate(templateID, userID.(uint), req.Name, req.Code); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// PublishTemplate 发布模板
func PublishTemplate(c *gin.Context) {
	templateID := utils.ParseUint(c.Param("template_id"))

	appService := &service.ApplicationService{}
	if err := appService.PublishTemplate(templateID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}
