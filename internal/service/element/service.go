package element

import (
	"github.com/iiwish/lingjian/internal/model"
	"github.com/jmoiron/sqlx"
)

// ElementService 元素服务
type ElementService struct {
	tableService     *TableService
	dimensionService *DimensionService
	menuService      *MenuService
	db               *sqlx.DB
}

// NewElementService 创建元素服务实例
func NewElementService(db *sqlx.DB) *ElementService {
	return &ElementService{
		tableService:     NewTableService(db),
		dimensionService: NewDimensionService(db),
		menuService:      NewMenuService(db),
		db:               db,
	}
}

// Table
func (s *ElementService) GetTableItems(tableID uint, page int, pageSize int, query *model.QueryCondition) ([]map[string]interface{}, int, error) {
	return s.tableService.GetTableItems(tableID, page, pageSize, query)
}

func (s *ElementService) CreateTableItems(tableItems []map[string]interface{}, creatorID uint, tableID uint) error {
	return s.tableService.CreateTableItems(tableItems, creatorID, tableID)
}

func (s *ElementService) UpdateTableItems(reqItems model.UpdateTableItemsRequest, updaterID uint, tableID uint) error {
	return s.tableService.UpdateTableItems(reqItems, updaterID, tableID)
}

func (s *ElementService) DeleteTableItems(operatorID uint, tableID uint, reqItems []map[string]interface{}) error {
	return s.tableService.DeleteTableItems(operatorID, tableID, reqItems)
}

// Dimension
func (s *ElementService) CreateDimensionItem(item *model.CreateDimensionItemReq, creatorID uint, dimID uint) (uint, error) {
	return s.dimensionService.CreateDimensionItem(item, creatorID, dimID)
}

func (s *ElementService) UpdateDimensionItem(item *model.UpdateDimensionItemReq, updaterID uint, dimID uint) error {
	return s.dimensionService.UpdateDimensionItem(item, updaterID, dimID)
}

func (s *ElementService) DeleteDimensionItems(operatorID uint, dimID uint, itemIDs []uint) error {
	return s.dimensionService.DeleteDimensionItems(operatorID, dimID, itemIDs)
}

func (s *ElementService) TreeDimensionItems(userID uint, itemID uint, nodeID uint, query_type string, query_level uint) ([]model.TreeDimensionItem, error) {
	return s.dimensionService.TreeDimensionItems(userID, itemID, nodeID, query_type, query_level)
}

func (s *ElementService) UpdateDimensionItemSort(updaterID uint, dimID uint, itemID uint, parentID uint, sort int) error {
	return s.dimensionService.UpdateDimensionItemSort(updaterID, dimID, itemID, parentID, sort)
}

// Menu
func (s *ElementService) GetMenuItems(userID uint, menuID uint, nodeID uint, queryType string, queryLevel uint) ([]*model.TreeMenuItem, error) {
	return s.menuService.GetMenuItems(userID, menuID, nodeID, queryType, queryLevel)
}

func (s *ElementService) CreateMenuItem(userID uint, menu *model.CreateMenuItemReq, menuID uint) (uint, error) {
	return s.menuService.CreateMenuItem(userID, menu, menuID)
}

func (s *ElementService) UpdateMenuItem(menu *model.UpdateMenuItemReq, updaterID uint, menuID uint) error {
	return s.menuService.UpdateMenuItem(menu, updaterID, menuID)
}

func (s *ElementService) DeleteMenuItem(operatorID uint, menuID uint, itemIDs []uint) error {
	return s.menuService.DeleteMenuItem(operatorID, menuID, itemIDs)
}

func (s *ElementService) UpdateMenuItemSort(userID uint, menuID uint, itemID uint, parentID uint, sort int) error {
	return s.menuService.UpdateMenuItemSort(userID, menuID, itemID, parentID, sort)
}
