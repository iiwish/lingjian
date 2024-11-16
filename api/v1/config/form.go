package config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
)

// @Summary      创建表单配置
// @Description  创建新的表单配置
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        request body config.CreateFormRequest true "创建表单配置请求参数"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/forms [post]
func (api *ConfigAPI) CreateForm(c *gin.Context) {
	var req model.ConfigForm
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	userID := uint(c.GetInt64("user_id"))
	id, err := api.configService.CreateForm(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ID": id})
}

// @Summary      更新表单配置
// @Description  更新已存在的表单配置
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        request body model.ConfigForm true "更新表单配置请求参数"
// @Success      200  {object}  model.ConfigForm
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/forms/{id} [put]
func (api *ConfigAPI) UpdateForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	var form model.ConfigForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}
	form.ID = uint(id)

	userID := uint(c.GetInt64("user_id"))
	if err := api.configService.UpdateForm(&form, userID); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, form)
}

// @Summary      获取表单配置列表
// @Description  获取表单配置列表
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Success      200  {array}   model.ConfigForm
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/forms [get]
func (api *ConfigAPI) ListForms(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Query("app_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid app_id"})
		return
	}

	forms, err := api.configService.ListForms(uint(appID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, forms)
}

// @Summary      获取表单配置详情
// @Description  获取指定表单配置的详细信息
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {object}  model.ConfigForm
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/forms/{id} [get]
func (api *ConfigAPI) GetForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	form, err := api.configService.GetForm(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, form)
}

// @Summary      删除表单配置
// @Description  删除指定的表单配置
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/forms/{id} [delete]
func (api *ConfigAPI) DeleteForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "invalid id"})
		return
	}

	if err := api.configService.DeleteForm(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
