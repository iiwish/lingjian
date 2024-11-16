package config

import (
	"fmt"
	"strings"

	"github.com/iiwish/lingjian/internal/model"
)

// CreateTableItem 创建数据表记录
func (s *TableService) CreateTableItem(tableItem map[string]interface{}, creatorID uint, tableID uint) (uint, error) {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 从配置表读取表配置
	var tableName string
	err = tx.Get(&tableName, "SELECT table_name FROM sys_config_tables WHERE id = ?", tableID)
	if err != nil {
		return 0, fmt.Errorf("get table name failed: %v", err)
	}

	// 构建插入SQL
	columns := make([]string, 0)
	values := make([]string, 0)
	args := make([]interface{}, 0)

	// 添加基础字段
	if val, ok := tableItem["creator_id"]; ok {
		columns = append(columns, "creator_id")
		values = append(values, "?")
		args = append(args, val)
	}
	if val, ok := tableItem["updater_id"]; ok {
		columns = append(columns, "updater_id")
		values = append(values, "?")
		args = append(args, val)
	}
	if val, ok := tableItem["created_at"]; ok {
		columns = append(columns, "created_at")
		values = append(values, "?")
		args = append(args, val)
	}
	if val, ok := tableItem["updated_at"]; ok {
		columns = append(columns, "updated_at")
		values = append(values, "?")
		args = append(args, val)
	}

	// 添加动态字段
	for k, v := range tableItem {
		if k != "creator_id" && k != "updater_id" && k != "created_at" && k != "updated_at" {
			columns = append(columns, k)
			values = append(values, "?")
			args = append(args, v)
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ","),
		strings.Join(values, ","))

	// 执行插入
	result, err := tx.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("insert table item failed: %v", err)
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("commit transaction failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id failed: %v", err)
	}

	return uint(id), nil
}

// BatchCreateTableItems 批量创建数据表记录
func (s *TableService) BatchCreateTableItems(tableItems []map[string]interface{}, creatorID uint, tableID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 从配置表读取表配置
	var tableName string
	err = tx.Get(&tableName, "SELECT table_name FROM sys_config_tables WHERE id = ?", tableID)
	if err != nil {
		return fmt.Errorf("get table name failed: %v", err)
	}

	// 构建批量插入SQL
	for _, item := range tableItems {
		columns := make([]string, 0)
		values := make([]string, 0)
		args := make([]interface{}, 0)

		// 添加基础字段
		if val, ok := item["creator_id"]; ok {
			columns = append(columns, "creator_id")
			values = append(values, "?")
			args = append(args, val)
		}
		if val, ok := item["updater_id"]; ok {
			columns = append(columns, "updater_id")
			values = append(values, "?")
			args = append(args, val)
		}
		if val, ok := item["created_at"]; ok {
			columns = append(columns, "created_at")
			values = append(values, "?")
			args = append(args, val)
		}
		if val, ok := item["updated_at"]; ok {
			columns = append(columns, "updated_at")
			values = append(values, "?")
			args = append(args, val)
		}

		// 添加动态字段
		for k, v := range item {
			if k != "creator_id" && k != "updater_id" && k != "created_at" && k != "updated_at" {
				columns = append(columns, k)
				values = append(values, "?")
				args = append(args, v)
			}
		}

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			tableName,
			strings.Join(columns, ","),
			strings.Join(values, ","))

		// 执行插入
		_, err := tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("insert table item failed: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// UpdateTableItem 更新数据表记录
func (s *TableService) UpdateTableItem(tableItem map[string]interface{}, updaterID uint, tableID uint, itemID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 从配置表读取表配置
	var tableName string
	err = tx.Get(&tableName, "SELECT table_name FROM sys_config_tables WHERE id = ?", tableID)
	if err != nil {
		return fmt.Errorf("get table name failed: %v", err)
	}

	// 构建更新SQL
	sets := make([]string, 0)
	args := make([]interface{}, 0)

	// 添加基础字段
	if val, ok := tableItem["updater_id"]; ok {
		sets = append(sets, "updater_id = ?")
		args = append(args, val)
	}
	if val, ok := tableItem["updated_at"]; ok {
		sets = append(sets, "updated_at = ?")
		args = append(args, val)
	}

	// 添加动态字段
	for k, v := range tableItem {
		if k != "updater_id" && k != "updated_at" {
			sets = append(sets, fmt.Sprintf("%s = ?", k))
			args = append(args, v)
		}
	}

	// 添加WHERE条件
	args = append(args, itemID)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?",
		tableName,
		strings.Join(sets, ","))

	// 执行更新
	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("update table item failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// GetTableItem 获取数据表记录
func (s *TableService) GetTableItem(tableID uint, itemID uint) (map[string]interface{}, error) {
	// 从配置表读取表配置
	var tableName string
	err := s.db.Get(&tableName, "SELECT table_name FROM sys_config_tables WHERE id = ?", tableID)
	if err != nil {
		return nil, fmt.Errorf("get table name failed: %v", err)
	}

	// 查询记录
	var result map[string]interface{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", tableName)
	err = s.db.Get(&result, query, itemID)
	if err != nil {
		return nil, fmt.Errorf("get table item failed: %v", err)
	}

	return result, nil
}

// ListTableItems 获取数据表记录列表
func (s *TableService) ListTableItems(tableID uint, page int, pageSize int, query *model.QueryCondition) ([]map[string]interface{}, int, error) {
	// 从配置表读取表配置
	var tableName string
	err := s.db.Get(&tableName, "SELECT table_name FROM sys_config_tables WHERE id = ?", tableID)
	if err != nil {
		return nil, 0, fmt.Errorf("get table name failed: %v", err)
	}

	// 构建基础查询SQL
	baseQuery, args := query.BuildQuery(tableName)

	// 查询总数
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS t", baseQuery)
	err = s.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("count table items failed: %v", err)
	}

	// 添加分页
	offset := (page - 1) * pageSize
	baseQuery += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	// 查询记录
	var results []map[string]interface{}
	err = s.db.Select(&results, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list table items failed: %v", err)
	}

	return results, total, nil
}

// DeleteTableItem 删除数据表记录
func (s *TableService) DeleteTableItem(operatorID uint, tableID uint, itemID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 从配置表读取表配置
	var tableName string
	err = tx.Get(&tableName, "SELECT table_name FROM sys_config_tables WHERE id = ?", tableID)
	if err != nil {
		return fmt.Errorf("get table name failed: %v", err)
	}

	// 执行删除
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	_, err = tx.Exec(query, itemID)
	if err != nil {
		return fmt.Errorf("delete table item failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// BatchDeleteTableItems 批量删除数据表记录
func (s *TableService) BatchDeleteTableItems(operatorID uint, tableID uint, itemIDs []uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 从配置表读取表配置
	var tableName string
	err = tx.Get(&tableName, "SELECT table_name FROM sys_config_tables WHERE id = ?", tableID)
	if err != nil {
		return fmt.Errorf("get table name failed: %v", err)
	}

	// 构建IN条件
	placeholders := make([]string, len(itemIDs))
	args := make([]interface{}, len(itemIDs))
	for i, id := range itemIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	// 执行批量删除
	query := fmt.Sprintf("DELETE FROM %s WHERE id IN (%s)", tableName, strings.Join(placeholders, ","))
	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("batch delete table items failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}
