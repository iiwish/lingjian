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

// 维度明细配置相关方法
func (s *ElementService) CreateDimensionItem(item *model.ConfigDimensionItem, creatorID uint, dimID uint) (uint, error) {
	return s.dimensionService.CreateDimensionItem(item, creatorID, dimID)
}

func (s *ElementService) BatchCreateDimensionItems(items []*model.ConfigDimensionItem, creatorID uint, dimID uint) error {
	return s.dimensionService.BatchCreateDimensionItems(items, creatorID, dimID)
}

func (s *ElementService) UpdateDimensionItem(item *model.ConfigDimensionItem, updaterID uint, dimID uint) error {
	return s.dimensionService.UpdateDimensionItem(item, updaterID, dimID)
}

func (s *ElementService) DeleteDimensionItem(operatorID uint, dimID uint, itemID uint) error {
	return s.dimensionService.DeleteDimensionItem(operatorID, dimID, itemID)
}

func (s *ElementService) BatchDeleteDimensionItems(operatorID uint, dimID uint, itemIDs []uint) error {
	return s.dimensionService.BatchDeleteDimensionItems(operatorID, dimID, itemIDs)
}

func (s *ElementService) TreeDimensionItems(itemID uint) ([]*model.TreeConfigDimensionItem, error) {
	return s.dimensionService.TreeDimensionItems(itemID)
}
