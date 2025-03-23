package config

import (
	"fmt"

	"encoding/json"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service/element"
	"github.com/jmoiron/sqlx"
)

// ModelService 数据模型配置服务
type ModelService struct {
	db *sqlx.DB
}

// NewModelService 创建数据模型配置服务实例
func NewModelService(db *sqlx.DB) *ModelService {
	return &ModelService{db: db}
}

func toJSONString2(v interface{}) string {
	if v == nil {
		return "{}"
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

// CreateModel 创建数据模型配置
func (s *ModelService) CreateModel(appID uint, userID uint, Req *model.CreateModelReq) (uint, error) {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	dataModel := &model.ConfigModel{
		AppID:         appID,
		ModelCode:     Req.ModelCode,
		DisplayName:   Req.DisplayName,
		Description:   Req.Description,
		Configuration: toJSONString2(Req.Configuration),
		Status:        1,
		CreatorID:     userID,
		UpdaterID:     userID,
	}

	// 插入数据模型配置
	result, err := tx.NamedExec(`
    INSERT INTO sys_config_models (
    app_id, model_code, display_name, description, configuration, status, created_at, creator_id, updated_at, updater_id
    ) VALUES (
    :app_id, :model_code, :display_name, :description, :configuration, :status, NOW(), :creator_id, NOW(), :creator_id
    )
    `, dataModel)
	if err != nil {
		return 0, fmt.Errorf("insert sys_config_models failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id failed: %v", err)
	}

	// 创建对应的系统菜单的menu
	menuService := element.NewMenuService(s.db)
	menu := &model.CreateMenuItemReq{
		ParentID:    Req.ParentID,
		MenuName:    dataModel.DisplayName,
		MenuCode:    dataModel.ModelCode,
		Description: dataModel.Description,
		MenuType:    5, // 表示model类型
		Status:      1,
		IconPath:    "model",
		SourceID:    uint(id),
	}

	err = menuService.CreateSysMenu(appID, userID, menu)
	if err != nil {
		return 0, fmt.Errorf("create menu failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction failed: %v", err)
	}

	return uint(id), nil
}

// UpdateModel 更新数据模型配置
func (s *ModelService) UpdateModel(appID uint, userID uint, Req *model.UpdateModelReq) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	dataModel := &model.ConfigModel{
		ID:            Req.ID,
		AppID:         appID,
		ModelCode:     Req.ModelCode,
		DisplayName:   Req.DisplayName,
		Description:   Req.Description,
		Configuration: toJSONString2(Req.Configuration),
		Status:        Req.Status,
		UpdaterID:     userID,
	}

	// 更新数据模型配置·
	_, err = tx.NamedExec(`
		UPDATE sys_config_models SET 
			app_id = :app_id,
			display_name = :display_name,
			description = :description,
			configuration = :configuration,
			status = :status,
			updated_at = NOW(),
			updater_id = :updater_id
		WHERE id = :id
	`, dataModel)
	if err != nil {
		return fmt.Errorf("update sys_config_models failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// GetModel 获取数据模型配置
func (s *ModelService) GetModel(id uint) (*model.ModelResp, error) {
	fmt.Printf("开始获取模型配置，ID: %d\n", id)

	var dataModel model.ConfigModel
	err := s.db.Get(&dataModel, "SELECT * FROM sys_config_models WHERE id = ?", id)
	if err != nil {
		fmt.Printf("数据库查询失败，错误: %v\n", err)
		return nil, fmt.Errorf("get model failed: %v", err)
	}
	fmt.Printf("数据库查询成功，获取到的数据: %+v\n", dataModel)

	var configItem model.ModelConfigItem
	fmt.Printf("开始解析configuration字段，原始数据: %s\n", dataModel.Configuration)
	err = json.Unmarshal([]byte(dataModel.Configuration), &configItem)
	if err != nil {
		fmt.Printf("configuration解析失败，错误: %v\n", err)
		return nil, fmt.Errorf("unmarshal configuration failed: %v", err)
	}
	fmt.Printf("configuration解析成功，解析后数据: %+v\n", configItem)

	resp := model.ModelResp{
		ID:            dataModel.ID,
		ModelCode:     dataModel.ModelCode,
		DisplayName:   dataModel.DisplayName,
		Description:   dataModel.Description,
		Configuration: configItem,
		Status:        dataModel.Status,
	}
	fmt.Printf("响应数据组装完成: %+v\n", resp)

	return &resp, nil
}

// DeleteModel 删除数据模型配置
func (s *ModelService) DeleteModel(id uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 删除数据模型配置
	_, err = tx.Exec("DELETE FROM sys_config_models WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete model failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}
