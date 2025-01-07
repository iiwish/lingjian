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
// @Param        request body model.CreateModelReq true "创建数据模型配置请求参数"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/models [post]
func (api *ConfigAPI) CreateModel(c *gin.Context) {
	var req model.CreateModelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	appID := c.GetUint("app_id")
	id, err := api.configService.CreateModel(appID, userID, &req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{"id": id})
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
// @Param        request body model.UpdateModelReq true "更新数据模型配置请求参数"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/models/{id} [put]
func (api *ConfigAPI) UpdateModel(c *gin.Context) {
	id := utils.ParseUint(c.Param("id"))
	if id == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var dataModel model.UpdateModelReq
	if err := c.ShouldBindJSON(&dataModel); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	dataModel.ID = uint(id)

	userID := c.GetUint("user_id")
	appID := c.GetUint("app_id")
	if err := api.configService.UpdateModel(appID, userID, &dataModel); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
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
// @Success      200  {object}  model.ModelResp
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /config/models/{id} [get]
func (api *ConfigAPI) GetModel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	dataModel, err := api.configService.GetModel(uint(id))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, dataModel)
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
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := api.configService.DeleteModel(uint(id)); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}
