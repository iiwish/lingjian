package config

import (
	"fmt"
	"strings"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/jmoiron/sqlx"
)

// TableService 数据表配置服务
type TableService struct {
	db *sqlx.DB
}

// NewTableService 创建数据表配置服务实例
func NewTableService(db *sqlx.DB) *TableService {
	return &TableService{db: db}
}

// CreateTable 创建数据表配置
func (s *TableService) CreateTable(tableinfo *model.CreateTableReq, creatorID uint) (uint, error) {
	var table model.ConfigTable
	table.AppID = tableinfo.AppID
	table.TableName = tableinfo.TableName
	table.DisplayName = tableinfo.DisplayName
	table.Description = tableinfo.Description
	table.Status = 1
	table.CreatorID = creatorID
	table.UpdaterID = creatorID

	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 插入数据表配置
	result, err := tx.NamedExec(`
        INSERT INTO sys_config_tables (
            app_id, table_name, display_name, description, func, status, created_at, creator_id, updated_at, updater_id
        ) VALUES (
            :app_id, :table_name, :display_name, :description, :func, :status, NOW(), :creator_id, NOW(), :creator_id
        )
    `, table)
	if err != nil {
		return 0, fmt.Errorf("insert sys_config_tables failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id failed: %v", err)
	}

	// 构建创建数据表的SQL语句
	createTableSQL := fmt.Sprintf(`
        CREATE TABLE %s (
    `, table.TableName)

	for _, field := range tableinfo.Fields {
		fieldSQL := fmt.Sprintf("%s %s", field.Name, field.Type)
		if field.NotNull {
			fieldSQL += " NOT NULL"
		}
		if field.AutoIncrement {
			fieldSQL += " AUTO_INCREMENT"
		}
		if field.Default != "" {
			fieldSQL += fmt.Sprintf(" DEFAULT '%s'", field.Default)
		}
		if field.Comment != "" {
			fieldSQL += fmt.Sprintf(" COMMENT '%s'", field.Comment)
		}
		createTableSQL += ", " + fieldSQL
	}

	for _, index := range tableinfo.Indexes {
		indexSQL := fmt.Sprintf("INDEX %s (%s)", index.Name, strings.Join(index.Fields, ", "))
		createTableSQL += ", " + indexSQL
	}

	createTableSQL = strings.TrimSuffix(createTableSQL, ", ") + ")"

	// 创建数据表
	_, err = tx.Exec(createTableSQL)
	if err != nil {
		return 0, fmt.Errorf("create table failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction failed: %v", err)
	}

	return uint(id), nil
}

// UpdateTable 更新数据表配置
func (s *TableService) UpdateTable(table *model.ConfigTable, updaterID uint) error {
	table.UpdaterID = updaterID

	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 获取数据表名称
	var tableName string
	err = tx.Get(&tableName, "SELECT table_name FROM sys_config_tables WHERE id = ?", table.ID)
	if err != nil {
		return fmt.Errorf("get table name failed: %v", err)
	}

	// 对比数据表名称是否有变化
	if tableName != table.TableName {
		// 检查新表名是否已存在
		var count int
		err = tx.Get(&count, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", table.TableName)
		if err != nil {
			return fmt.Errorf("check table name failed: %v", err)
		}
		if count > 0 {
			return fmt.Errorf("table name already exists")
		}

		// 修改数据表名称
		_, err = tx.Exec("RENAME TABLE " + tableName + " TO " + table.TableName)
		if err != nil {
			return fmt.Errorf("rename table failed: %v", err)
		}
	}

	// 更新数据表配置
	_, err = tx.NamedExec(`
		UPDATE sys_config_tables SET 
			table_name = :table_name,
			display_name = :display_name,
			description = :description,
			func = :func,
			status = :status,
			updater_id = :updater_id,
			updated_at = NOW()
		WHERE id = :id
	`, table)
	if err != nil {
		return fmt.Errorf("update sys_config_tables failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// GetTable 获取数据表配置
func (s *TableService) GetTable(id uint) (*model.CreateTableReq, error) {
	var table model.ConfigTable
	query := `
        SELECT 
            id, app_id, table_name, display_name, description, 
            IFNULL(func, '') AS func, status, created_at, creator_id, updated_at, updater_id 
        FROM sys_config_tables 
        WHERE id = ?
    `
	err := s.db.Get(&table, query, id)
	if err != nil {
		return nil, fmt.Errorf("get table failed: %v", err)
	}

	var tableInfo model.CreateTableReq
	tableInfo.AppID = table.AppID
	tableInfo.TableName = table.TableName
	tableInfo.DisplayName = table.DisplayName
	tableInfo.Description = table.Description
	tableInfo.Func = table.Func

	// 根据数据表名称获取字段信息
	var fields []model.MySQLField
	query = "SELECT " +
		"`COLUMN_NAME` AS `Field`, " +
		"`COLUMN_TYPE` AS `Type`, " +
		"IFNULL(`COLLATION_NAME`, '') AS `Collation`, " +
		"ORDINAL_POSITION AS `Sort`, " +
		"`IS_NULLABLE` AS `Null`, " +
		"`COLUMN_KEY` AS `Key`, " +
		"IFNULL(`COLUMN_DEFAULT`, '') AS `Default`, " +
		"`EXTRA` AS `Extra`, " +
		"`PRIVILEGES` AS `Privileges`, " +
		"IFNULL(`COLUMN_COMMENT`, '') AS `Comment` " +
		"FROM `information_schema`.`columns` " +
		"WHERE `TABLE_SCHEMA` = DATABASE() AND `TABLE_NAME` = ?"
	err = s.db.Select(&fields, query, table.TableName)
	if err != nil {
		return nil, fmt.Errorf("get fields failed: %v", err)
	}

	for _, field := range fields {
		var f model.Field
		f.Name = field.Field
		f.Comment = field.Comment
		f.Type = field.Type
		f.Sort = field.Sort
		f.PrimaryKey = field.Key == "PRI"
		f.AutoIncrement = field.Extra == "auto_increment"
		f.NotNull = field.Null == "NO"
		f.Default = field.Default
		tableInfo.Fields = append(tableInfo.Fields, f)
	}

	// 根据数据表名称获取索引信息
	var indexes []model.MySQLIndex
	err = s.db.Select(&indexes, "SHOW INDEX FROM "+table.TableName)
	if err != nil {
		return nil, fmt.Errorf("get indexes failed: %v", err)
	}

	indexMap := make(map[string]*model.Index)
	for _, index := range indexes {
		if idx, ok := indexMap[index.KeyName]; ok {
			idx.Fields = append(idx.Fields, index.ColumnName)
		} else {
			indexMap[index.KeyName] = &model.Index{
				Name:   index.KeyName,
				Type:   index.IndexType,
				Fields: []string{index.ColumnName},
			}
		}
	}

	for _, idx := range indexMap {
		tableInfo.Indexes = append(tableInfo.Indexes, *idx)
	}

	return &tableInfo, nil
}

// DeleteTable 删除数据表配置
func (s *TableService) DeleteTable(id uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 获取数据表名称
	var tableName string
	err = tx.Get(&tableName, "SELECT table_name FROM sys_config_tables WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("get table name failed: %v", err)
	}

	// 直接删除数据表
	_, err = tx.Exec("DELETE FROM sys_config_tables WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete table failed: %v", err)
	}

	// 删除数据表
	_, err = tx.Exec("DROP TABLE IF EXISTS " + tableName)
	if err != nil {
		return fmt.Errorf("drop table failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// ListTables 获取数据表配置列表
func (s *TableService) ListTables(appID uint) ([]model.ConfigTable, error) {
	var tables []model.ConfigTable
	err := s.db.Select(&tables, "SELECT * FROM sys_config_tables WHERE app_id = ? AND status = 1 ORDER BY id DESC", appID)
	if err != nil {
		return nil, fmt.Errorf("list tables failed: %v", err)
	}
	return tables, nil
}
