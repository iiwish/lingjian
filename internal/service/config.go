package service

import (
	"encoding/json"
	"fmt"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/jmoiron/sqlx"
)

// ConfigService 配置服务
type ConfigService struct {
	db *sqlx.DB
}

// NewConfigService 创建配置服务实例
func NewConfigService(db *sqlx.DB) *ConfigService {
	return &ConfigService{db: db}
}

// CreateTableRequest 创建数据表配置请求
type CreateTableRequest struct {
	ApplicationID  uint               `json:"application_id" binding:"required"`
	Name           string             `json:"name" binding:"required"`
	Code           string             `json:"code" binding:"required"`
	Description    string             `json:"description"`
	MySQLTableName string             `json:"mysql_table_name" binding:"required"`
	Fields         []model.TableField `json:"fields" binding:"required"`
	Indexes        []model.TableIndex `json:"indexes"`
}

// CreateTable 创建数据表配置
func (s *ConfigService) CreateTable(req *CreateTableRequest, creatorID uint) error {
	// 将字段和索引转换为JSON字符串
	fields, err := json.Marshal(req.Fields)
	if err != nil {
		return fmt.Errorf("marshal fields failed: %v", err)
	}

	indexes, err := json.Marshal(req.Indexes)
	if err != nil {
		return fmt.Errorf("marshal indexes failed: %v", err)
	}

	// 创建数据表配置
	table := &model.ConfigTable{
		ApplicationID:  req.ApplicationID,
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		MySQLTableName: req.MySQLTableName,
		Fields:         string(fields),
		Indexes:        string(indexes),
		Status:         1,
		Version:        1,
	}

	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 插入数据表配置
	result, err := tx.NamedExec(`
		INSERT INTO config_tables (
			application_id, name, code, description, mysql_table_name,
			fields, indexes, status, version, created_at, updated_at
		) VALUES (
			:application_id, :name, :code, :description, :mysql_table_name,
			:fields, :indexes, :status, :version, NOW(), NOW()
		)
	`, table)
	if err != nil {
		return fmt.Errorf("insert config_tables failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id failed: %v", err)
	}

	// 创建版本记录
	version := &model.ConfigVersion{
		ApplicationID: req.ApplicationID,
		ConfigType:    "table",
		ConfigID:      uint(id),
		Version:       1,
		Content:       string(fields), // 使用字段定义作为版本内容
		CreatorID:     creatorID,
	}

	_, err = tx.NamedExec(`
		INSERT INTO config_versions (
			application_id, config_type, config_id, version,
			content, creator_id, created_at
		) VALUES (
			:application_id, :config_type, :config_id, :version,
			:content, :creator_id, NOW()
		)
	`, version)
	if err != nil {
		return fmt.Errorf("insert config_versions failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// CreateDimensionRequest 创建维度配置请求
type CreateDimensionRequest struct {
	ApplicationID  uint                  `json:"application_id" binding:"required"`
	Name           string                `json:"name" binding:"required"`
	Code           string                `json:"code" binding:"required"`
	Type           string                `json:"type" binding:"required"`
	MySQLTableName string                `json:"mysql_table_name" binding:"required"`
	Configuration  model.DimensionConfig `json:"configuration" binding:"required"`
}

// CreateDimension 创建维度配置
func (s *ConfigService) CreateDimension(req *CreateDimensionRequest, creatorID uint) error {
	// 将配置转换为JSON字符串
	config, err := json.Marshal(req.Configuration)
	if err != nil {
		return fmt.Errorf("marshal configuration failed: %v", err)
	}

	// 创建维度配置
	dimension := &model.ConfigDimension{
		ApplicationID:  req.ApplicationID,
		Name:           req.Name,
		Code:           req.Code,
		Type:           req.Type,
		MySQLTableName: req.MySQLTableName,
		Configuration:  string(config),
		Status:         1,
		Version:        1,
	}

	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 插入维度配置
	result, err := tx.NamedExec(`
		INSERT INTO config_dimensions (
			application_id, name, code, type, mysql_table_name,
			configuration, status, version, created_at, updated_at
		) VALUES (
			:application_id, :name, :code, :type, :mysql_table_name,
			:configuration, :status, :version, NOW(), NOW()
		)
	`, dimension)
	if err != nil {
		return fmt.Errorf("insert config_dimensions failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id failed: %v", err)
	}

	// 创建版本记录
	version := &model.ConfigVersion{
		ApplicationID: req.ApplicationID,
		ConfigType:    "dimension",
		ConfigID:      uint(id),
		Version:       1,
		Content:       string(config),
		CreatorID:     creatorID,
	}

	_, err = tx.NamedExec(`
		INSERT INTO config_versions (
			application_id, config_type, config_id, version,
			content, creator_id, created_at
		) VALUES (
			:application_id, :config_type, :config_id, :version,
			:content, :creator_id, NOW()
		)
	`, version)
	if err != nil {
		return fmt.Errorf("insert config_versions failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// GetDimensionValues 获取维度值列表
func (s *ConfigService) GetDimensionValues(dimensionID uint, filter map[string]any) ([]map[string]any, error) {
	// 获取维度配置
	var dimension model.ConfigDimension
	err := s.db.Get(&dimension, "SELECT * FROM config_dimensions WHERE id = ?", dimensionID)
	if err != nil {
		return nil, fmt.Errorf("get dimension failed: %v", err)
	}

	// 解析维度配置
	var config model.DimensionConfig
	if err := json.Unmarshal([]byte(dimension.Configuration), &config); err != nil {
		return nil, fmt.Errorf("unmarshal configuration failed: %v", err)
	}

	// 构建查询SQL
	query := fmt.Sprintf("SELECT * FROM %s", dimension.MySQLTableName)
	var params []any

	// 添加过滤条件
	if len(filter) > 0 {
		query += " WHERE "
		var conditions []string
		for k, v := range filter {
			conditions = append(conditions, fmt.Sprintf("%s = ?", k))
			params = append(params, v)
		}
		query += fmt.Sprint(conditions[0])
		for i := 1; i < len(conditions); i++ {
			query += fmt.Sprintf(" AND %s", conditions[i])
		}
	}

	// 添加排序
	if config.OrderField != "" {
		query += fmt.Sprintf(" ORDER BY %s", config.OrderField)
	}

	// 执行查询
	var values []map[string]any
	err = s.db.Select(&values, query, params...)
	if err != nil {
		return nil, fmt.Errorf("query dimension values failed: %v", err)
	}

	return values, nil
}

// UpdateTable 更新数据表配置
func (s *ConfigService) UpdateTable(table *model.ConfigTable) error {
	// TODO: 实现更新数据表配置的逻辑
	return nil
}

// GetTable 获取数据表配置
func (s *ConfigService) GetTable(id uint) (*model.ConfigTable, error) {
	// TODO: 实现获取数据表配置的逻辑
	return nil, nil
}

// DeleteTable 删除数据表配置
func (s *ConfigService) DeleteTable(id uint) error {
	// TODO: 实现删除数据表配置的逻辑
	return nil
}

// ListTables 获取数据表配置列表
func (s *ConfigService) ListTables(appID uint) ([]model.ConfigTable, error) {
	// TODO: 实现获取数据表配置列表的逻辑
	return nil, nil
}

// GetTableVersions 获取数据表配置版本历史
func (s *ConfigService) GetTableVersions(id uint) ([]model.ConfigVersion, error) {
	// TODO: 实现获取数据表配置版本历史的逻辑
	return nil, nil
}

// RollbackTable 回滚数据表配置到指定版本
func (s *ConfigService) RollbackTable(id uint, version int) error {
	// TODO: 实现回滚数据表配置的逻辑
	return nil
}

// UpdateDimension 更新维度配置
func (s *ConfigService) UpdateDimension(dimension *model.ConfigDimension) error {
	// TODO: 实现更新维度配置的逻辑
	return nil
}

// GetDimension 获取维度配置
func (s *ConfigService) GetDimension(id uint) (*model.ConfigDimension, error) {
	// TODO: 实现获取维度配置的逻辑
	return nil, nil
}

// DeleteDimension 删除维度配置
func (s *ConfigService) DeleteDimension(id uint) error {
	// TODO: 实现删除维度配置的逻辑
	return nil
}

// ListDimensions 获取维度配置列表
func (s *ConfigService) ListDimensions(appID uint) ([]model.ConfigDimension, error) {
	// TODO: 实现获取维度配置列表的逻辑
	return nil, nil
}

// GetDimensionVersions 获取维度配置版本历史
func (s *ConfigService) GetDimensionVersions(id uint) ([]model.ConfigVersion, error) {
	// TODO: 实现获取维度配置版本历史的逻辑
	return nil, nil
}

// RollbackDimension 回滚维度配置到指定版本
func (s *ConfigService) RollbackDimension(id uint, version int) error {
	// TODO: 实现回滚维度配置的逻辑
	return nil
}

// 以下是数据模型配置的服务方法...
func (s *ConfigService) CreateModel(model *model.ConfigDataModel) error {
	// TODO: 实现创建数据模型配置的逻辑
	return nil
}

func (s *ConfigService) UpdateModel(model *model.ConfigDataModel) error {
	// TODO: 实现更新数据模型配置的逻辑
	return nil
}

func (s *ConfigService) GetModel(id uint) (*model.ConfigDataModel, error) {
	// TODO: 实现获取数据模型配置的逻辑
	return nil, nil
}

func (s *ConfigService) DeleteModel(id uint) error {
	// TODO: 实现删除数据模型配置的逻辑
	return nil
}

func (s *ConfigService) ListModels(appID uint) ([]model.ConfigDataModel, error) {
	// TODO: 实现获取数据模型配置列表的逻辑
	return nil, nil
}

func (s *ConfigService) GetModelVersions(id uint) ([]model.ConfigVersion, error) {
	// TODO: 实现获取数据模型配置版本历史的逻辑
	return nil, nil
}

func (s *ConfigService) RollbackModel(id uint, version int) error {
	// TODO: 实现回滚数据模型配置的逻辑
	return nil
}

// 以下是表单配置的服务方法...
func (s *ConfigService) CreateForm(form *model.ConfigForm) error {
	// TODO: 实现创建表单配置的逻辑
	return nil
}

func (s *ConfigService) UpdateForm(form *model.ConfigForm) error {
	// TODO: 实现更新表单配置的逻辑
	return nil
}

func (s *ConfigService) GetForm(id uint) (*model.ConfigForm, error) {
	// TODO: 实现获取表单配置的逻辑
	return nil, nil
}

func (s *ConfigService) DeleteForm(id uint) error {
	// TODO: 实现删除表单配置的逻辑
	return nil
}

func (s *ConfigService) ListForms(appID uint) ([]model.ConfigForm, error) {
	// TODO: 实现获取表单配置列表的逻辑
	return nil, nil
}

func (s *ConfigService) GetFormVersions(id uint) ([]model.ConfigVersion, error) {
	// TODO: 实现获取表单配置版本历史的逻辑
	return nil, nil
}

func (s *ConfigService) RollbackForm(id uint, version int) error {
	// TODO: 实现回滚表单配置的逻辑
	return nil
}

// 以下是菜单配置的服务方法...
func (s *ConfigService) CreateMenu(menu *model.ConfigMenu) error {
	// TODO: 实现创建菜单配置的逻辑
	return nil
}

func (s *ConfigService) UpdateMenu(menu *model.ConfigMenu) error {
	// TODO: 实现更新菜单配置的逻辑
	return nil
}

func (s *ConfigService) GetMenu(id uint) (*model.ConfigMenu, error) {
	// TODO: 实现获取菜单配置的逻辑
	return nil, nil
}

func (s *ConfigService) DeleteMenu(id uint) error {
	// TODO: 实现删除菜单配置的逻辑
	return nil
}

func (s *ConfigService) ListMenus(appID uint) ([]model.ConfigMenu, error) {
	// TODO: 实现获取菜单配置列表的逻辑
	return nil, nil
}

func (s *ConfigService) GetMenuVersions(id uint) ([]model.ConfigVersion, error) {
	// TODO: 实现获取菜单配置版本历史的逻辑
	return nil, nil
}

func (s *ConfigService) RollbackMenu(id uint, version int) error {
	// TODO: 实现回滚菜单配置的逻辑
	return nil
}
