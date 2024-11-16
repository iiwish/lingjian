package config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
)

// @Summary      创建维度配置
// @Description  创建新的维度配置
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        dimension body model.ConfigDimension true "创建维度配置请求参数"
// @Success      201  {object}  gin.H{"ID": uint}
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions [post]
func (api *ConfigAPI) CreateDimension(c *gin.Context) {
	var dimension model.ConfigDimension
	if err := c.ShouldBindJSON(&dimension); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	// 校验请求参数
	if dimension.AppID != uint(c.GetInt64("app_id")) {
		c.JSON(http.StatusBadRequest, Response{Error: "app_id与请求路径中的ID不一致"})
		return
	}
	if dimension.TableName == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "table_name不能为空"})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	id, err := api.configService.CreateDimension(&dimension, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ID": id})
}

// @Summary      更新维度配置
// @Description  更新已存在的维度配置
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Param        dim_id path int true "配置ID"
// @Param        dimension body model.ConfigDimension true "更新维度配置请求参数"
// @Success      200  {object}  model.ConfigDimension
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions/{dim_id} [put]
func (api *ConfigAPI) UpdateDimension(c *gin.Context) {
	// 获取请求参数
	id, err := strconv.ParseUint(c.Param("dim_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	// 绑定请求参数
	var dimension model.ConfigDimension
	if err := c.ShouldBindJSON(&dimension); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}
	dimension.ID = uint(id)

	// 校验请求参数
	if dimension.ID != uint(c.GetInt64("app_id")) {
		c.JSON(http.StatusBadRequest, Response{Error: "app_id与请求路径中的ID不一致"})
		return
	}

	if dimension.TableName == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "table_name不能为空"})
		return
	}

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
// @Success      200  {array}   model.ConfigDimension
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions [get]
func (api *ConfigAPI) ListDimensions(c *gin.Context) {
	appID := uint(c.GetInt64("app_id"))

	dimensions, err := api.configService.ListDimensions(appID)
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
// @Param        dim_id path int true "配置ID"
// @Success      200  {object}  model.ConfigDimension
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions/{dim_id} [get]
func (api *ConfigAPI) GetDimension(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("dim_id"), 10, 64)
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
// @Param        dim_id path int true "配置ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions/{dim_id} [delete]
func (api *ConfigAPI) DeleteDimension(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("dim_id"), 10, 64)
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
