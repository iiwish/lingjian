package element

import (
	"github.com/iiwish/lingjian/internal/model"
	"github.com/jmoiron/sqlx"
)

// ElementService 元素服务
type ElementService struct {
	tableService     *TableService
	dimensionService *DimensionService
	db               *sqlx.DB
}

// NewElementService 创建元素服务实例
func NewElementService(db *sqlx.DB) *ElementService {
	return &ElementService{
		tableService:     NewTableService(db),
		dimensionService: NewDimensionService(db),
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
func (s *ElementService) CreateDimensionItems(items []*model.DimensionItem, creatorID uint, dimID uint) ([]uint, error) {
	return s.dimensionService.CreateDimensionItems(items, creatorID, dimID)
}

func (s *ElementService) UpdateDimensionItems(item []*model.DimensionItem, updaterID uint, dimID uint) error {
	return s.dimensionService.UpdateDimensionItems(item, updaterID, dimID)
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
