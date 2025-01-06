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
func (s *ConfigService) GetTable(id uint) (*model.CreateTableReq, error) {
	return s.tableService.GetTable(id)
}

func (s *ConfigService) CreateTable(table *model.CreateTableReq, creatorID uint, appID uint) (uint, error) {
	return s.tableService.CreateTable(table, creatorID, appID)
}

func (s *ConfigService) UpdateTable(tableID uint, req *model.TableUpdateReq, updaterID uint, appID uint) error {
	return s.tableService.UpdateTable(tableID, req, updaterID, appID)
}

func (s *ConfigService) DeleteTable(id uint) error {
	return s.tableService.DeleteTable(id)
}

// 维度配置相关方法
func (s *ConfigService) CreateDimension(dimension *model.CreateDimReq, creatorID uint, appID uint) (uint, error) {
	return s.dimensionService.CreateDimension(dimension, creatorID, appID)
}

func (s *ConfigService) UpdateDimension(req *model.UpdateDimensionReq, updaterID uint) error {
	return s.dimensionService.UpdateDimension(req, updaterID)
}

func (s *ConfigService) GetDimension(id uint) (*model.GetDimResp, error) {
	return s.dimensionService.GetDimension(id)
}

func (s *ConfigService) DeleteDimension(id uint) error {
	return s.dimensionService.DeleteDimension(id)
}

func (s *ConfigService) GetDimensions(userID uint, appID uint, dimType string) ([]model.GetDimResp, error) {
	return s.dimensionService.GetDimensions(userID, appID, dimType)
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
func (s *ConfigService) CreateMenu(userID uint, appID uint, menu *model.CreateMenuReq) (uint, error) {
	return s.menuService.CreateMenu(userID, appID, menu)
}

func (s *ConfigService) UpdateMenu(menu *model.UpdateMenuReq, updaterID uint, dimID uint) error {
	return s.menuService.UpdateMenu(menu, updaterID, dimID)
}

func (s *ConfigService) DeleteMenu(id uint) error {
	return s.menuService.DeleteMenu(id)
}

func (s *ConfigService) GetSystemMenuID(appID uint) (uint, error) {
	return s.menuService.GetSystemMenuID(appID)
}

func (s *ConfigService) GetMenuList(userID uint, appID uint) ([]model.GetDimResp, error) {
	return s.menuService.GetMenuList(userID, appID)
}

func (s *ConfigService) GetMenuByID(id uint) (*model.GetDimResp, error) {
	return s.menuService.GetMenuByID(id)
}
