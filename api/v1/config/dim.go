package config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service/config"
)

// @Summary      创建维度配置
// @Description  创建新的维度配置
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        request body config.CreateDimensionRequest true "创建维度配置请求参数"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions [post]
func (api *ConfigAPI) CreateDimension(c *gin.Context) {
	var req config.CreateDimensionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	if err := api.configService.CreateDimension(&req, userID); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, Response{})
}

// @Summary      更新维度配置
// @Description  更新已存在的维度配置
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        request body model.ConfigDimension true "更新维度配置请求参数"
// @Success      200  {object}  model.ConfigDimension
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions/{id} [put]
func (api *ConfigAPI) UpdateDimension(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	var dimension model.ConfigDimension
	if err := c.ShouldBindJSON(&dimension); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}
	dimension.ID = uint(id)

	userID := uint(c.GetInt64("user_id"))
	if err := api.configService.UpdateDimension(&dimension, userID); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dimension)
}

// @Summary      获取维度配置列表
// @Description  获取指定应用的维度配置列表
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        app_id query int true "应用ID"
// @Success      200  {array}   model.ConfigDimension
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions [get]
func (api *ConfigAPI) ListDimensions(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Query("app_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid app_id"})
		return
	}

	dimensions, err := api.configService.ListDimensions(uint(appID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dimensions)
}

// @Summary      获取维度配置详情
// @Description  获取指定维度配置的详细信息
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {object}  model.ConfigDimension
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions/{id} [get]
func (api *ConfigAPI) GetDimension(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	dimension, err := api.configService.GetDimension(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dimension)
}

// @Summary      删除维度配置
// @Description  删除指定的维度配置
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions/{id} [delete]
func (api *ConfigAPI) DeleteDimension(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	if err := api.configService.DeleteDimension(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
