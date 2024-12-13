package element

import (
	"fmt"
	"strings"
	"time"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
	"github.com/jmoiron/sqlx"
)

// TableService 数据表元素服务
type TableService struct {
	db *sqlx.DB
}

// NewTableService 创建数据表元素服务实例
func NewTableService(db *sqlx.DB) *TableService {
	return &TableService{db: db}
}

// GetTableItems 获取数据表记录列表
func (s *TableService) GetTableItems(tableID uint, page int, pageSize int, query *model.QueryCondition) ([]map[string]interface{}, int, error) {
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

// BatchCreateTableItems 批量创建数据表记录
func (s *TableService) CreateTableItems(tableItems []map[string]interface{}, creatorID uint, tableID uint) error {
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

// UpdateTableItems 更新数据表记录
func (s *TableService) UpdateTableItems(req model.UpdateTableItemsRequest, updaterID uint, tableID uint) error {
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

	for _, item := range req.Items {
		// 构建更新SQL
		sets := make([]string, 0)
		args := make([]interface{}, 0)

		// 添加基础字段
		if val, ok := item["updater_id"]; ok {
			sets = append(sets, "updater_id = ?")
			args = append(args, val)
		} else {
			sets = append(sets, "updater_id = ?")
			args = append(args, updaterID)
		}
		if val, ok := item["updated_at"]; ok {
			sets = append(sets, "updated_at = ?")
			args = append(args, val)
		} else {
			sets = append(sets, "updated_at = ?")
			args = append(args, time.Now())
		}

		// 添加动态字段
		primaryKeyValues := make([]interface{}, 0)
		for k, v := range item {
			if k != "updater_id" && k != "updated_at" && !utils.Contains(req.PrimaryKeyColumns, k) {
				sets = append(sets, fmt.Sprintf("%s = ?", k))
				args = append(args, v)
			} else if utils.Contains(req.PrimaryKeyColumns, k) {
				primaryKeyValues = append(primaryKeyValues, v)
			}
		}

		// 添加WHERE条件
		whereClauses := make([]string, 0)
		for _, pk := range req.PrimaryKeyColumns {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", pk))
		}
		args = append(args, primaryKeyValues...)

		query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
			tableName,
			strings.Join(sets, ","),
			strings.Join(whereClauses, " AND "))

		// 执行更新
		_, err = tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("update table item failed: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// DeleteTableItems 批量删除数据表记录
func (s *TableService) DeleteTableItems(operatorID uint, tableID uint, req []map[string]interface{}) error {
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

	for _, condition := range req {
		// 构建WHERE条件
		whereClauses := make([]string, 0)
		args := make([]interface{}, 0)
		for col, val := range condition {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", col))
			args = append(args, val)
		}

		query := fmt.Sprintf("DELETE FROM %s WHERE %s",
			tableName,
			strings.Join(whereClauses, " AND "))

		// 使用预编译语句执行删除
		stmt, err := tx.Preparex(query)
		if err != nil {
			return fmt.Errorf("prepare statement failed: %v", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(args...)
		if err != nil {
			return fmt.Errorf("delete table item failed: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}
