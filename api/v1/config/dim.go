package config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      创建维度配置
// @Description  创建新的维度配置
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        dimension body model.CreateDimReq true "创建维度配置请求参数"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions [post]
func (api *ConfigAPI) CreateDimension(c *gin.Context) {
	var dimension model.CreateDimReq
	if err := c.ShouldBindJSON(&dimension); err != nil {
		utils.ServerError(c, err)
		return
	}

	// 校验请求参数
	if dimension.TableName == "" {
		utils.Error(c, http.StatusBadRequest, "table_name不能为空")
		return
	}

	AppID := c.GetUint("app_id")
	userID := c.GetUint("user_id")
	id, err := api.configService.CreateDimension(&dimension, userID, AppID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{"ID": id})
}

// @Summary      更新维度配置
// @Description  更新已存在的维度配置
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        dim_id path int true "配置ID"
// @Param        dimension body model.UpdateDimensionReq true "更新维度配置请求参数"
// @Success      200  {object}  model.UpdateDimensionReq
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions/{dim_id} [put]
func (api *ConfigAPI) UpdateDimension(c *gin.Context) {
	// 获取请求参数
	id, err := strconv.ParseUint(c.Param("dim_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	// 绑定请求参数
	var dimension model.UpdateDimensionReq
	if err := c.ShouldBindJSON(&dimension); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	dimension.ID = uint(id)

	userID := c.GetUint("user_id")
	if err := api.configService.UpdateDimension(&dimension, userID); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary      获取维度配置详情
// @Description  获取指定维度配置的详细信息
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        dim_id path int true "配置ID"
// @Success      200  {object}  model.GetDimResp
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions/{dim_id} [get]
func (api *ConfigAPI) GetDimensionByID(c *gin.Context) {
	id := utils.ParseUint(c.Param("dim_id"))
	if id == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	dimension, err := api.configService.GetDimension(uint(id))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, dimension)
}

// @Summary      删除维度配置
// @Description  删除指定的维度配置
// @Tags         ConfigDimension
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer token"
// @Param        App-ID header string true "应用ID"
// @Param        dim_id path int true "配置ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/dimensions/{dim_id} [delete]
func (api *ConfigAPI) DeleteDimension(c *gin.Context) {
	id := utils.ParseUint(c.Param("dim_id"))
	if id == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := api.configService.DeleteDimension(uint(id)); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}
