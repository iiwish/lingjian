package v1

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service"
	"github.com/iiwish/lingjian/pkg/utils"
)

// RegisterAppRoutes 注册应用相关路由
func RegisterAppRoutes(r *gin.RouterGroup) {
	app := r.Group("/apps")
	{
		app.GET("", ListApps)             // 获取应用列表
		app.GET("/:app_id", GetApp)       // 获取应用详情
		app.POST("", CreateApp)           // 创建应用
		app.PUT("/:app_id", UpdateApp)    // 更新应用
		app.DELETE("/:app_id", DeleteApp) // 删除应用
	}
}

// CreateAppRequest 创建应用请求结构
type CreateAppRequest struct {
	Name        string `json:"name" example:"测试应用" binding:"required"`
	Description string `json:"description" example:"这是一个测试应用"`
	Status      int    `json:"status" example:"1"`
}

// @Summary      获取应用列表
// @Description  获取所有应用列表
// @Tags         Application
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Success      200  {object}  utils.Response{data=map[string]interface{}}
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

// @Summary      获取应用详情
// @Description  获取应用详情
// @Tags         Application
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        app_id path int true "应用ID"
// @Success      200  {object}  utils.Response{data=model.App}
// @Failure      500  {object}  utils.Response
// @Router       /apps/{app_id} [get]
func GetApp(c *gin.Context) {
	appID := utils.ParseUint(c.Param("app_id"))
	if appID == 0 {
		utils.Error(c, 400, "无效的 app_id 参数")
		return
	}

	userID := c.GetUint("user_id")
	appService := &service.AppService{}
	result, err := appService.GetAppByID(appID, userID)
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
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        request body CreateAppRequest true "创建应用请求参数"
// @Success      200  {object}  utils.Response{data=model.App}
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

// @Summary      更新应用
// @Description  更新已存在的应用
// @Tags         Application
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        app_id path int true "应用ID"
// @Param        request body CreateAppRequest true "更新应用请求参数"
// @Success      200  {object}  utils.Response{data=model.App}
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /apps/{app_id} [put]
func UpdateApp(c *gin.Context) {
	log.Printf("UpdateApp: %v", c)
	appID := utils.ParseUint(c.Param("app_id"))
	if appID == 0 {
		utils.Error(c, 400, "无效的 app_id 参数")
		return
	}

	var req model.App
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}
	req.ID = appID

	// 获取当前用户ID
	userID := c.GetUint("user_id")

	log.Printf("UpdateApp: appID=%d, userID=%d, req=%+v", appID, userID, req)

	appService := &service.AppService{}
	err := appService.UpdateApp(&req, userID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      删除应用
// @Description  删除指定的应用
// @Tags         Application
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        app_id path int true "应用ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /apps/{app_id} [delete]
func DeleteApp(c *gin.Context) {
	appID := utils.ParseUint(c.Param("app_id"))
	if appID == 0 {
		utils.Error(c, 400, "无效的 app_id 参数")
		return
	}

	// 获取当前用户ID
	userID := c.GetUint("user_id")

	appService := &service.AppService{}
	err := appService.DeleteApp(appID, userID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}
