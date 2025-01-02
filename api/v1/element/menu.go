package element

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/pkg/utils"
)

// @Summary      获取菜单明细树
// @Description  获取菜单明细树
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization header string true "Bearer Token"
// @Param        App-ID header string true "应用ID"
// @Param        menu_id  path int true "菜单ID"
// @Param        id 	query int     false "节点ID，不指定则返回整个维度配置项树"
// @Param 	 	 type      query    string  false  "菜单类型，可选值为 'children'、'descendants'、'leaves' , 默认为 'descendants'"
// @Param 	  	 level     query    int     false  "树的层级，可选值为 0、1、2、3， 默认为 0不指定层级"
// @Success      200 {array}   []model.TreeDimensionItem
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /menu/{menu_id} [get]
func (api *ElementAPI) GetMenuItems(c *gin.Context) {
	// 获取id参数
	menuID := utils.ParseUint(c.Param("menu_id"))
	if menuID == 0 {
		utils.Error(c, http.StatusBadRequest, "menu_id不能为空")
		return
	}

	// 获取请求参数
	nodeID := utils.ParseUint(c.Query("id"))

	// 获取query参数
	queryLevel := utils.ParseUint(c.Query("level"))

	// 获取type参数,默认为descendants
	queryType := c.Query("type")
	if queryType == "" {
		queryType = "descendants"
	}

	userID := c.GetUint("user_id")
	menus, err := api.elementService.GetMenuItems(userID, menuID, nodeID, queryType, queryLevel)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, menus)
}
