package element

import (
	"fmt"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/jmoiron/sqlx"
)

// DimensionService 维度配置服务
type DimensionService struct {
	db *sqlx.DB
}

// NewDimensionService 创建维度配置服务实例
func NewDimensionService(db *sqlx.DB) *DimensionService {
	return &DimensionService{db: db}
}

// CreateDimensionItem 创建维度明细配置
func (s *DimensionService) CreateDimensionItem(dimension *model.ConfigDimensionItem, creatorID uint, dim_id uint) (uint, error) {
	dimension.Status = 1
	dimension.CreatorID = creatorID
	dimension.UpdaterID = creatorID

	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 从配置表读取维度配置
	var table_name string
	err = tx.Get(&table_name, "SELECT table_name FROM sys_config_dimensions WHERE id = ?", dim_id)
	if err != nil {
		return 0, fmt.Errorf("get table name failed: %v", err)
	}

	// 插入维度配置
	result, err := tx.NamedExec(fmt.Sprintf(`
        INSERT INTO %s (node_id, parent_id, name, code, description, level, sort, status, created_at, creator_id, updated_at, updater_id)
        VALUES (:node_id, :parent_id, :name, :code, :description, :level, :sort, :status, Now(), :creator_id, Now(), :updater_id)
    `, table_name), dimension)
	if err != nil {
		return 0, fmt.Errorf("insert sys_config_dimensions failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction failed: %v", err)
	}

	return uint(id), nil
}

// BatchCreateDimensionItems 批量创建维度明细配置
func (s *DimensionService) BatchCreateDimensionItems(dimensionItems []*model.ConfigDimensionItem, creatorID uint, dim_id uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 从配置表读取维度配置
	var table_name string
	err = tx.Get(&table_name, "SELECT table_name FROM sys_config_dimensions WHERE id = ?", dim_id)
	if err != nil {
		return fmt.Errorf("get table name failed: %v", err)
	}

	// 插入维度配置
	for _, dimension := range dimensionItems {
		dimension.Status = 1
		dimension.CreatorID = creatorID
		dimension.UpdaterID = creatorID

		_, err := tx.NamedExec(fmt.Sprintf(`
			INSERT INTO %s (node_id, parent_id, name, code, description, level, sort, status, created_at, creator_id, updated_at, updater_id)
			VALUES (:node_id, :parent_id, :name, :code, :description, :level, :sort, :status, Now(), :creator_id, Now(), :updater_id)
		`, table_name), dimension)
		if err != nil {
			return fmt.Errorf("insert sys_config_dimensions failed: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// UpdateDimensionItem 更新维度明细配置
func (s *DimensionService) UpdateDimensionItem(dimension *model.ConfigDimensionItem, updaterID uint, dim_id uint) error {
	dimension.UpdaterID = updaterID

	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 从配置表读取维度配置
	var table_name string
	err = tx.Get(&table_name, "SELECT table_name FROM sys_config_dimensions WHERE id = ?", dim_id)
	if err != nil {
		return fmt.Errorf("get table name failed: %v", err)
	}

	// 更新维度配置
	_, err = tx.NamedExec(fmt.Sprintf(`
		UPDATE %s SET node_id = :node_id, parent_id = :parent_id, name = :name, code = :code, description = :description, level = :level, sort = :sort, status = :status, updated_at = Now(), updater_id = :updater_id
		WHERE id = :id
	`, table_name), dimension)
	if err != nil {
		return fmt.Errorf("update sys_config_dimensions failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// TreeDimensionItems 获取维度明细配置树形结构
func (s *DimensionService) TreeDimensionItems(dim_id uint) ([]*model.TreeConfigDimensionItem, error) {
	// 从配置表读取维度配置
	var tableName string
	err := s.db.Get(&tableName, "SELECT table_name FROM sys_config_dimensions WHERE id = ?", dim_id)
	if err != nil {
		return nil, fmt.Errorf("获取表名失败: %v", err)
	}

	// 校验 tableName 是否有效，防止SQL注入
	allowedTableNames := map[string]bool{
		"dimension_table1": true,
		"dimension_table2": true,
		// 添加其他允许的表名
	}
	if !allowedTableNames[tableName] {
		return nil, fmt.Errorf("无效的表名: %s", tableName)
	}

	// 查询维度配置
	var dimensionItems []*model.TreeConfigDimensionItem
	query := fmt.Sprintf(`
        SELECT * FROM %s ORDER BY sort
    `, tableName)
	err = s.db.Select(&dimensionItems, query)
	if err != nil {
		return nil, fmt.Errorf("查询维度配置失败: %v", err)
	}

	// 构建ID到节点的映射
	itemMap := make(map[uint]*model.TreeConfigDimensionItem)
	for _, item := range dimensionItems {
		item.Children = []*model.TreeConfigDimensionItem{} // 初始化Children切片
		itemMap[item.ID] = item
	}

	// 构建树形结构
	var roots []*model.TreeConfigDimensionItem
	for _, item := range dimensionItems {
		if item.ParentID == 0 {
			roots = append(roots, item) // 根节点
		} else {
			if parent, exists := itemMap[item.ParentID]; exists {
				parent.Children = append(parent.Children, item)
			}
		}
	}

	return roots, nil
}

// DeleteDimensionItem 删除维度明细配置
func (s *DimensionService) DeleteDimensionItem(operatorID uint, dim_id uint, itemID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 从配置表读取维度配置
	var tableName string
	err = tx.Get(&tableName, "SELECT table_name FROM sys_config_dimensions WHERE id = ?", dim_id)
	if err != nil {
		return fmt.Errorf("get table name failed: %v", err)
	}

	// 删除维度配置
	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName), itemID)
	if err != nil {
		return fmt.Errorf("delete dimension failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// BatchDeleteDimensionItems 批量删除维度明细配置
func (s *DimensionService) BatchDeleteDimensionItems(operatorID uint, dim_id uint, itemIDs []uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 从配置表读取维度配置
	var tableName string
	err = tx.Get(&tableName, "SELECT table_name FROM sys_config_dimensions WHERE id = ?", dim_id)
	if err != nil {
		return fmt.Errorf("get table name failed: %v", err)
	}

	// 删除维度配置
	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE id IN (?)", tableName), itemIDs)
	if err != nil {
		return fmt.Errorf("delete dimension failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}
