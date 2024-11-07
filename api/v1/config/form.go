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
// @Param        request body model.ConfigForm true "创建表单配置请求参数"
// @Success      201  {object}  model.ConfigForm
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/forms [post]
func (api *ConfigAPI) CreateForm(c *gin.Context) {
	var form model.ConfigForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := api.configService.CreateForm(&form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, form)
}

// @Summary      更新表单配置
// @Description  更新已存在的表单配置
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        request body model.ConfigForm true "更新表单配置请求参数"
// @Success      200  {object}  model.ConfigForm
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/forms/{id} [put]
func (api *ConfigAPI) UpdateForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var form model.ConfigForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	form.ID = uint(id)

	if err := api.configService.UpdateForm(&form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, form)
}

// @Summary      获取表单配置列表
// @Description  获取指定应用的表单配置列表
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        app_id query int true "应用ID"
// @Success      200  {array}   model.ConfigForm
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/forms [get]
func (api *ConfigAPI) ListForms(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Query("app_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid app_id"})
		return
	}

	forms, err := api.configService.ListForms(uint(appID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/forms/{id} [get]
func (api *ConfigAPI) GetForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	form, err := api.configService.GetForm(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/forms/{id} [delete]
func (api *ConfigAPI) DeleteForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := api.configService.DeleteForm(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary      获取表单配置版本历史
// @Description  获取指定表单配置的版本历史记录
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Success      200  {array}   model.ConfigFormVersion
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/forms/{id}/versions [get]
func (api *ConfigAPI) GetFormVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	versions, err := api.configService.GetFormVersions(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// @Summary      回滚表单配置
// @Description  将表单配置回滚到指定版本
// @Tags         ConfigForm
// @Accept       json
// @Produce      json
// @Param        id path int true "配置ID"
// @Param        version query int true "目标版本号"
// @Success      200  {object}  nil
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /config/forms/{id}/rollback [post]
func (api *ConfigAPI) RollbackForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	version, err := strconv.Atoi(c.Query("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version"})
		return
	}

	if err := api.configService.RollbackForm(uint(id), version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
