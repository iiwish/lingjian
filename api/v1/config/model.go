package config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service/config"
)

// @Summary      创建数据模型配置
// @Description  创建新的数据模型配置
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Param        request body config.CreateModelRequest true "创建数据模型配置请求参数"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/models [post]
func (api *ConfigAPI) CreateModel(c *gin.Context) {
	var req config.CreateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	if err := api.configService.CreateModel(&req, userID); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, Response{})
}

// @Summary      更新数据模型配置
// @Description  更新已存在的数据模型配置
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        request body model.ConfigDataModel true "更新数据模型配置请求参数"
// @Success      200  {object}  model.ConfigDataModel
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/models/{id} [put]
func (api *ConfigAPI) UpdateModel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	var dataModel model.ConfigModel
	if err := c.ShouldBindJSON(&dataModel); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}
	dataModel.ID = uint(id)

	userID := uint(c.GetInt64("user_id"))
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
// @Param        app_id query int true "应用ID"
// @Success      200  {array}   model.ConfigDataModel
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
// @Param        id path int true "配置ID"
// @Success      200  {object}  model.ConfigDataModel
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

// @Summary      获取数据模型配置版本历史
// @Description  获取指定数据模型配置的版本历史记录
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {array}   model.ConfigVersion
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/models/{id}/versions [get]
func (api *ConfigAPI) GetModelVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	versions, err := api.configService.GetModelVersions(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// @Summary      回滚数据模型配置
// @Description  将数据模型配置回滚到指定版本
// @Tags         ConfigModel
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        version query int true "目标版本号"
// @Success      200  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/models/{id}/rollback [post]
func (api *ConfigAPI) RollbackModel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	version, err := strconv.Atoi(c.Query("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid version"})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	if err := api.configService.RollbackModel(uint(id), version, userID); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
