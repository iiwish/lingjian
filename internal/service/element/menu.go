package element

import (
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type MenuService struct {
	db *sqlx.DB
}

func NewMenuService(db *sqlx.DB) *MenuService {
	return &MenuService{db: db}
}

// GetMenuItems 获取菜单明细配置
func (s *MenuService) GetMenuItems(userID uint, menuID uint, itemID uint, queryType string, queryLevel uint) ([]*model.TreeMenuItem, error) {
	dimService := NewDimensionService(s.db)
	dimItems, err := dimService.TreeDimensionItems(userID, menuID, itemID, queryType, queryLevel)
	if err != nil {
		return nil, err
	}

	// 转换函数
	var convert func(items []model.TreeDimensionItem) []*model.TreeMenuItem
	convert = func(items []model.TreeDimensionItem) []*model.TreeMenuItem {
		var menuItems []*model.TreeMenuItem
		for _, item := range items {
			menuItem := &model.TreeMenuItem{
				ID:        item.ID,
				NodeID:    item.NodeID,
				ParentID:  item.ParentID,
				MenuName:  item.Name,
				MenuCode:  item.Code,
				Level:     item.Level,
				Sort:      item.Sort,
				Status:    item.Status,
				CreatedAt: item.CreatedAt,
				CreatorID: item.CreatorID,
				UpdatedAt: item.UpdatedAt,
				UpdaterID: item.UpdaterID,
				// 从CustomData中获取自定义列
				SourceID: utils.ParseUint(item.CustomData["source_id"]),
				MenuType: utils.ParseInt(item.CustomData["menu_type"]),
				IconPath: item.CustomData["icon_path"],
			}

			// 递归转换子节点
			if len(item.Children) > 0 {
				// 将[]*TreeDimensionItem转换为[]TreeDimensionItem
				children := make([]model.TreeDimensionItem, len(item.Children))
				for i, child := range item.Children {
					children[i] = *child
				}
				menuItem.Children = convert(children)
			}

			menuItems = append(menuItems, menuItem)
		}
		return menuItems
	}

	return convert(dimItems), nil
}
