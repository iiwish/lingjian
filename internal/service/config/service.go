package config

import (
	"github.com/iiwish/lingjian/internal/model"
	"github.com/jmoiron/sqlx"
)

// ConfigService 配置服务
type ConfigService struct {
	tableService     *TableService
	dimensionService *DimensionService
	modelService     *ModelService
	formService      *FormService
	menuService      *MenuService
}

// NewConfigService 创建配置服务实例
func NewConfigService(db *sqlx.DB) *ConfigService {
	return &ConfigService{
		tableService:     NewTableService(db),
		dimensionService: NewDimensionService(db),
		modelService:     NewModelService(db),
		formService:      NewFormService(db),
		menuService:      NewMenuService(db),
	}
}

// 表配置相关方法
func (s *ConfigService) CreateTable(table *model.CreateTableReq, creatorID uint) (uint, error) {
	return s.tableService.CreateTable(table, creatorID)
}

func (s *ConfigService) UpdateTable(table *model.ConfigTable, updaterID uint) error {
	return s.tableService.UpdateTable(table, updaterID)
}

func (s *ConfigService) GetTable(id uint) (*model.CreateTableReq, error) {
	return s.tableService.GetTable(id)
}

func (s *ConfigService) DeleteTable(id uint) error {
	return s.tableService.DeleteTable(id)
}

func (s *ConfigService) ListTables(appID uint) ([]model.ConfigTable, error) {
	return s.tableService.ListTables(appID)
}

// 维度配置相关方法
func (s *ConfigService) CreateDimension(dimension *model.ConfigDimension, creatorID uint) (uint, error) {
	return s.dimensionService.CreateDimension(dimension, creatorID)
}

func (s *ConfigService) UpdateDimension(dimension *model.ConfigDimension, updaterID uint) error {
	return s.dimensionService.UpdateDimension(dimension, updaterID)
}

func (s *ConfigService) GetDimension(id uint) (*model.ConfigDimension, error) {
	return s.dimensionService.GetDimension(id)
}

func (s *ConfigService) DeleteDimension(id uint) error {
	return s.dimensionService.DeleteDimension(id)
}

func (s *ConfigService) ListDimensions(appID uint) ([]model.ConfigDimension, error) {
	return s.dimensionService.ListDimensions(appID)
}

// 数据模型配置相关方法
func (s *ConfigService) CreateModel(dataModel *model.ConfigModel, creatorID uint) (uint, error) {
	return s.modelService.CreateModel(dataModel, creatorID)
}

func (s *ConfigService) UpdateModel(dataModel *model.ConfigModel, updaterID uint) error {
	return s.modelService.UpdateModel(dataModel, updaterID)
}

func (s *ConfigService) GetModel(id uint) (*model.ConfigModel, error) {
	return s.modelService.GetModel(id)
}

func (s *ConfigService) DeleteModel(id uint) error {
	return s.modelService.DeleteModel(id)
}

func (s *ConfigService) ListModels(appID uint) ([]model.ConfigModel, error) {
	return s.modelService.ListModels(appID)
}

// 表单配置相关方法
func (s *ConfigService) CreateForm(form *model.ConfigForm, creatorID uint) (uint, error) {
	return s.formService.CreateForm(form, creatorID)
}

func (s *ConfigService) UpdateForm(form *model.ConfigForm, updaterID uint) error {
	return s.formService.UpdateForm(form, updaterID)
}

func (s *ConfigService) GetForm(id uint) (*model.ConfigForm, error) {
	return s.formService.GetForm(id)
}

func (s *ConfigService) DeleteForm(id uint) error {
	return s.formService.DeleteForm(id)
}

func (s *ConfigService) ListForms(appID uint) ([]model.ConfigForm, error) {
	return s.formService.ListForms(appID)
}

// 菜单配置相关方法
func (s *ConfigService) CreateMenu(menu *model.ConfigMenu, creatorID uint) (uint, error) {
	return s.menuService.CreateMenu(menu, creatorID)
}

func (s *ConfigService) UpdateMenu(menu *model.ConfigMenu, updaterID uint) error {
	return s.menuService.UpdateMenu(menu, updaterID)
}

func (s *ConfigService) GetMenuByID(id uint) (*model.ConfigMenu, error) {
	return s.menuService.GetMenuByID(id)
}

func (s *ConfigService) DeleteMenu(id uint) error {
	return s.menuService.DeleteMenu(id)
}

func (s *ConfigService) GetMenus(appID uint, operatorID uint, level *int, parentID *uint, menuType string) ([]model.TreeConfigMenu, error) {
	return s.menuService.GetMenus(appID, operatorID, level, parentID, menuType)
}
