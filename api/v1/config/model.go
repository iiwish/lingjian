package config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      创建数据模型配置
// @Description  创建新的数据模型配置
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        request body model.ConfigModel true "创建数据模型配置请求参数"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/models [post]
func (api *ConfigAPI) CreateModel(c *gin.Context) {
	var req model.ConfigModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	id, err := api.configService.CreateModel(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ID": id})
}

// @Summary      更新数据模型配置
// @Description  更新已存在的数据模型配置
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        id path int true "配置ID"
// @Param        request body model.ConfigModel true "更新数据模型配置请求参数"
// @Success      200  {object}  model.ConfigModel
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/models/{id} [put]
func (api *ConfigAPI) UpdateModel(c *gin.Context) {
	id := utils.ParseUint(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	var dataModel model.ConfigModel
	if err := c.ShouldBindJSON(&dataModel); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}
	dataModel.ID = uint(id)

	userID := c.GetUint("user_id")
	if err := api.configService.UpdateModel(&dataModel, userID); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dataModel)
}

// @Summary      获取数据模型配置列表
// @Description  获取指定应用的数据模型配置列表
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Success      200  {array}   model.ConfigModel
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/models [get]
func (api *ConfigAPI) ListModels(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Query("app_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid app_id"})
		return
	}

	models, err := api.configService.ListModels(uint(appID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models)
}

// @Summary      获取数据模型配置详情
// @Description  获取指定数据模型配置的详细信息
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        id path int true "配置ID"
// @Success      200  {object}  model.ConfigModel
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/models/{id} [get]
func (api *ConfigAPI) GetModel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	dataModel, err := api.configService.GetModel(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dataModel)
}

// @Summary      删除数据模型配置
// @Description  删除指定的数据模型配置
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        id path int true "配置ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/models/{id} [delete]
func (api *ConfigAPI) DeleteModel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	if err := api.configService.DeleteModel(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
