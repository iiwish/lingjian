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
		app.GET("", ListApps) // 获取应用列表
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

// @Summary      获取应用列表
// @Description  获取所有应用列表
// @Tags         Application
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /apps [get]
func ListApps(c *gin.Context) {
	appService := &service.AppService{}
	result, err := appService.ListApps()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, result)
}

type CreateAppRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
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
	var req CreateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	appService := &service.AppService{}
	result, err := appService.CreateApp(req.Name, req.Code, req.Description)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, result)
}

type AssignAppRequest struct {
	IsDefault bool `json:"is_default"`
}

// @Summary      为用户分配应用
// @Description  将应用分配给指定用户
// @Tags         Application
// @Accept       json
// @Produce      json
// @Param        app_id  path     int  true  "应用ID"
// @Param        user_id path     int  true  "用户ID"
// @Param        request body     AssignAppRequest true "分配应用请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /apps/{app_id}/users/{user_id} [post]
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

// @Summary      获取用户应用列表
// @Description  获取指定用户的所有应用
// @Tags         Application
// @Accept       json
// @Produce      json
// @Param        user_id path     int  true  "用户ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /apps/users/{user_id} [get]
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

// @Summary      获取用户默认应用
// @Description  获取指定用户的默认应用
// @Tags         Application
// @Accept       json
// @Produce      json
// @Param        user_id path     int  true  "用户ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /apps/users/{user_id}/default [get]
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

// @Summary      创建应用模板
// @Description  创建新的应用模板
// @Tags         ApplicationTemplate
// @Accept       json
// @Produce      json
// @Param        request body CreateTemplateRequest true "创建模板请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /apps/templates [post]
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

// @Summary      获取模板列表
// @Description  获取所有已上架的应用模板
// @Tags         ApplicationTemplate
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /apps/templates [get]
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

// @Summary      从模板创建应用
// @Description  基于指定模板创建新应用
// @Tags         ApplicationTemplate
// @Accept       json
// @Produce      json
// @Param        template_id path     int  true  "模板ID"
// @Param        request     body     CreateFromTemplateRequest true "创建应用请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /apps/templates/{template_id}/create [post]
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

// @Summary      发布模板
// @Description  将应用模板发布上架
// @Tags         ApplicationTemplate
// @Accept       json
// @Produce      json
// @Param        template_id path     int  true  "模板ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /apps/templates/{template_id}/publish [post]
func PublishTemplate(c *gin.Context) {
	templateID := utils.ParseUint(c.Param("template_id"))

	appService := &service.AppService{}
	if err := appService.PublishTemplate(templateID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}
