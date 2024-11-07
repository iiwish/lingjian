package config

import (
	"encoding/json"

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
func (s *ConfigService) CreateTable(req *CreateTableRequest, creatorID uint) error {
	// 将字段和索引转换为JSON字符串
	fields, err := json.Marshal(req.Fields)
	if err != nil {
		return err
	}

	indexes, err := json.Marshal(req.Indexes)
	if err != nil {
		return err
	}

	// 创建数据表配置
	table := &model.ConfigTable{
		AppID:          req.AppID,
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		MySQLTableName: req.MySQLTableName,
		Fields:         string(fields),
		Indexes:        string(indexes),
		Status:         1,
		Version:        1,
	}

	return s.tableService.CreateTable(table, creatorID)
}

func (s *ConfigService) UpdateTable(table *model.ConfigTable, updaterID uint) error {
	return s.tableService.UpdateTable(table, updaterID)
}

func (s *ConfigService) GetTable(id uint) (*model.ConfigTable, error) {
	return s.tableService.GetTable(id)
}

func (s *ConfigService) DeleteTable(id uint) error {
	return s.tableService.DeleteTable(id)
}

func (s *ConfigService) ListTables(appID uint) ([]model.ConfigTable, error) {
	return s.tableService.ListTables(appID)
}

func (s *ConfigService) GetTableVersions(id uint) ([]model.ConfigVersion, error) {
	return s.tableService.GetTableVersions(id)
}

func (s *ConfigService) RollbackTable(id uint, version int, updaterID uint) error {
	return s.tableService.RollbackTable(id, version, updaterID)
}

// 维度配置相关方法
func (s *ConfigService) CreateDimension(req *CreateDimensionRequest, creatorID uint) error {
	// 将配置转换为JSON字符串
	config, err := json.Marshal(req.Configuration)
	if err != nil {
		return err
	}

	// 创建维度配置
	dimension := &model.ConfigDimension{
		AppID:          req.AppID,
		Name:           req.Name,
		Code:           req.Code,
		Type:           req.Type,
		MySQLTableName: req.MySQLTableName,
		Configuration:  string(config),
		Status:         1,
		Version:        1,
	}

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

func (s *ConfigService) GetDimensionVersions(id uint) ([]model.ConfigVersion, error) {
	return s.dimensionService.GetDimensionVersions(id)
}

func (s *ConfigService) RollbackDimension(id uint, version int, updaterID uint) error {
	return s.dimensionService.RollbackDimension(id, version, updaterID)
}

func (s *ConfigService) GetDimensionValues(dimensionID uint, filter map[string]any) ([]map[string]any, error) {
	return s.dimensionService.GetDimensionValues(dimensionID, filter)
}

// 数据模型配置相关方法
func (s *ConfigService) CreateModel(req *CreateModelRequest, creatorID uint) error {
	// 将字段、维度和指标转换为JSON字符串
	fields, err := json.Marshal(req.Fields)
	if err != nil {
		return err
	}

	dimensions, err := json.Marshal(req.Dimensions)
	if err != nil {
		return err
	}

	metrics, err := json.Marshal(req.Metrics)
	if err != nil {
		return err
	}

	// 创建数据模型配置
	dataModel := &model.ConfigDataModel{
		AppID:      req.AppID,
		Name:       req.Name,
		Code:       req.Code,
		TableID:    req.TableID,
		Fields:     string(fields),
		Dimensions: string(dimensions),
		Metrics:    string(metrics),
		Status:     1,
		Version:    1,
	}

	return s.modelService.CreateModel(dataModel, creatorID)
}

func (s *ConfigService) UpdateModel(dataModel *model.ConfigDataModel, updaterID uint) error {
	return s.modelService.UpdateModel(dataModel, updaterID)
}

func (s *ConfigService) GetModel(id uint) (*model.ConfigDataModel, error) {
	return s.modelService.GetModel(id)
}

func (s *ConfigService) DeleteModel(id uint) error {
	return s.modelService.DeleteModel(id)
}

func (s *ConfigService) ListModels(appID uint) ([]model.ConfigDataModel, error) {
	return s.modelService.ListModels(appID)
}

func (s *ConfigService) GetModelVersions(id uint) ([]model.ConfigVersion, error) {
	return s.modelService.GetModelVersions(id)
}

func (s *ConfigService) RollbackModel(id uint, version int, updaterID uint) error {
	return s.modelService.RollbackModel(id, version, updaterID)
}

// 表单配置相关方法
func (s *ConfigService) CreateForm(req *CreateFormRequest, creatorID uint) error {
	// 将布局、字段、规则和事件转换为JSON字符串
	layout, err := json.Marshal(req.Layout)
	if err != nil {
		return err
	}

	fields, err := json.Marshal(req.Fields)
	if err != nil {
		return err
	}

	rules, err := json.Marshal(req.Rules)
	if err != nil {
		return err
	}

	events, err := json.Marshal(req.Events)
	if err != nil {
		return err
	}

	// 创建表单配置
	form := &model.ConfigForm{
		AppID:   req.AppID,
		Name:    req.Name,
		Code:    req.Code,
		Type:    req.Type,
		TableID: req.TableID,
		Layout:  string(layout),
		Fields:  string(fields),
		Rules:   string(rules),
		Events:  string(events),
		Status:  1,
		Version: 1,
	}

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

func (s *ConfigService) GetFormVersions(id uint) ([]model.ConfigVersion, error) {
	return s.formService.GetFormVersions(id)
}

func (s *ConfigService) RollbackForm(id uint, version int, updaterID uint) error {
	return s.formService.RollbackForm(id, version, updaterID)
}

// 菜单配置相关方法
func (s *ConfigService) CreateMenu(req *CreateMenuRequest, creatorID uint) error {
	// 创建菜单配置
	menu := &model.ConfigMenu{
		AppID:     req.AppID,
		ParentID:  req.ParentID,
		Name:      req.Name,
		Code:      req.Code,
		Icon:      req.Icon,
		Path:      req.Path,
		Component: req.Component,
		Sort:      req.Sort,
		Status:    1,
		Version:   1,
	}

	return s.menuService.CreateMenu(menu, creatorID)
}

func (s *ConfigService) UpdateMenu(menu *model.ConfigMenu, updaterID uint) error {
	return s.menuService.UpdateMenu(menu, updaterID)
}

func (s *ConfigService) GetMenu(id uint) (*model.ConfigMenu, error) {
	return s.menuService.GetMenu(id)
}

func (s *ConfigService) DeleteMenu(id uint) error {
	return s.menuService.DeleteMenu(id)
}

func (s *ConfigService) ListMenus(appID uint) ([]model.ConfigMenu, error) {
	return s.menuService.ListMenus(appID)
}

func (s *ConfigService) GetMenuVersions(id uint) ([]model.ConfigVersion, error) {
	return s.menuService.GetMenuVersions(id)
}

func (s *ConfigService) RollbackMenu(id uint, version int, updaterID uint) error {
	return s.menuService.RollbackMenu(id, version, updaterID)
}
