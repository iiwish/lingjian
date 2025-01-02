package element

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
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

// BatchCreateDimensionItems 批量创建维度明细配置
func (s *DimensionService) CreateDimensionItem(dimensionItem *model.DimensionItem, creatorID uint, dim_id uint) (uint, error) {
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

	dimensionItem.Status = 1
	dimensionItem.CreatorID = creatorID
	dimensionItem.UpdaterID = creatorID
	// 获取父节点的node_id
	if dimensionItem.ParentID != 0 {
		var parent struct {
			NodeID string `db:"node_id"`
			Level  int    `db:"level"`
		}
		err = tx.Get(&parent, fmt.Sprintf("SELECT node_id, level FROM %s WHERE id = ?", table_name), dimensionItem.ParentID)
		if err != nil {
			return 0, fmt.Errorf("get parent node_id and level failed: %v", err)
		}
		dimensionItem.NodeID = parent.NodeID
		dimensionItem.Level = parent.Level + 1
	} else {
		dimensionItem.NodeID = ""
		dimensionItem.Level = 1
	}

	// 获取sort数值
	var sort int
	err = tx.Get(&sort, fmt.Sprintf("SELECT IFNULL(MAX(sort), 0) FROM %s WHERE parent_id = ?", table_name), dimensionItem.ParentID)
	if err != nil {
		return 0, fmt.Errorf("get max sort failed: %v", err)
	}
	dimensionItem.Sort = sort + 1

	// 获取表的列信息
	var columns []struct {
		ColumnName string `db:"COLUMN_NAME"`
	}
	err = tx.Select(&columns, `
			SELECT COLUMN_NAME 
			FROM INFORMATION_SCHEMA.COLUMNS 
			WHERE TABLE_SCHEMA = DATABASE() 
			AND TABLE_NAME = ? 
			AND COLUMN_NAME NOT IN (
				'id', 'node_id', 'parent_id', 'name', 'code', 'description',
				'level', 'sort', 'status', 'created_at', 'creator_id', 
				'updated_at', 'updater_id'
			)
		`, table_name)
	if err != nil {
		return 0, fmt.Errorf("get columns failed: %v", err)
	}

	// 构建INSERT语句
	var insertSQL strings.Builder
	insertSQL.WriteString(fmt.Sprintf(`
			INSERT INTO %s (
				node_id, parent_id, name, code, description, 
				level, sort, status, created_at, creator_id, 
				updated_at, updater_id
		`, table_name))

	// 添加自定义列
	for _, col := range columns {
		insertSQL.WriteString(", " + col.ColumnName)
	}

	insertSQL.WriteString(") VALUES (")
	insertSQL.WriteString(`
			:node_id, :parent_id, :name, :code, :description,
			:level, :sort, :status, Now(), :creator_id,
			Now(), :updater_id
		`)

	// 添加自定义列的值
	for _, col := range columns {
		insertSQL.WriteString(", :" + col.ColumnName)
	}
	insertSQL.WriteString(")")

	// 准备参数
	params := map[string]interface{}{
		"node_id":     dimensionItem.NodeID,
		"parent_id":   dimensionItem.ParentID,
		"name":        dimensionItem.Name,
		"code":        dimensionItem.Code,
		"description": dimensionItem.Description,
		"level":       dimensionItem.Level,
		"sort":        dimensionItem.Sort,
		"status":      dimensionItem.Status,
		"creator_id":  dimensionItem.CreatorID,
		"updater_id":  dimensionItem.UpdaterID,
	}

	// 添加自定义列的值
	for _, col := range columns {
		if val, ok := dimensionItem.CustomData[col.ColumnName]; ok {
			params[col.ColumnName] = val
		} else {
			params[col.ColumnName] = "" // 默认空字符串
		}
	}

	result, err := tx.NamedExec(insertSQL.String(), params)
	if err != nil {
		return 0, fmt.Errorf("insert sys_config_dimensions failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id failed: %v", err)
	}

	// 更新node_id为node_id拼接下划线和id
	_, err = tx.Exec(fmt.Sprintf(`
			UPDATE %s SET node_id = CASE 
				WHEN node_id = '' OR node_id IS NULL THEN CAST(? AS CHAR)
				ELSE CONCAT(node_id, '_', ?)
			END
			WHERE id = ?
		`, table_name), id, id, id)
	if err != nil {
		return 0, fmt.Errorf("update node_id failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction failed: %v", err)
	}

	return uint(id), nil
}

// UpdateDimensionItem 更新维度明细配置
func (s *DimensionService) UpdateDimensionItems(dimensionItems []*model.DimensionItem, updaterID uint, dim_id uint) error {
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

	for _, dimension := range dimensionItems {
		dimension.UpdaterID = updaterID

		// 获取表的列信息
		var columns []struct {
			ColumnName string `db:"COLUMN_NAME"`
		}
		err = tx.Select(&columns, `
			SELECT COLUMN_NAME 
			FROM INFORMATION_SCHEMA.COLUMNS 
			WHERE TABLE_SCHEMA = DATABASE() 
			AND TABLE_NAME = ? 
			AND COLUMN_NAME NOT IN (
				'id', 'node_id', 'parent_id', 'name', 'code', 'description',
				'level', 'sort', 'status', 'created_at', 'creator_id', 
				'updated_at', 'updater_id'
			)
		`, table_name)
		if err != nil {
			return fmt.Errorf("get columns failed: %v", err)
		}

		// 构建UPDATE语句
		var updateSQL strings.Builder
		updateSQL.WriteString(fmt.Sprintf(`
			UPDATE %s SET 
				name = :name, 
				code = :code, 
				description = :description, 
				status = :status, 
				updated_at = Now(), 
				updater_id = :updater_id
		`, table_name))

		// 添加自定义列
		for _, col := range columns {
			updateSQL.WriteString(fmt.Sprintf(", %s = :%s", col.ColumnName, col.ColumnName))
		}

		updateSQL.WriteString(" WHERE id = :id")

		// 准备参数
		params := map[string]interface{}{
			"id":          dimension.ID,
			"name":        dimension.Name,
			"code":        dimension.Code,
			"description": dimension.Description,
			"status":      dimension.Status,
			"updater_id":  dimension.UpdaterID,
		}

		// 添加自定义列的值
		for _, col := range columns {
			if val, ok := dimension.CustomData[col.ColumnName]; ok {
				params[col.ColumnName] = val
			} else {
				params[col.ColumnName] = "" // 默认空字符串
			}
		}

		// 更新维度配置
		_, err = tx.NamedExec(updateSQL.String(), params)
		if err != nil {
			return fmt.Errorf("update sys_config_dimensions failed: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// TreeDimensionItems 获取维度明细配置树形结构
func (s *DimensionService) TreeDimensionItems(userID uint, dim_id uint, id uint, query_type string, query_level uint) ([]model.TreeDimensionItem, error) {
	// 首先检查用户是否有该维度的权限
	permissionQuery := `
		SELECT DISTINCT p.item_id FROM sys_permissions p
		INNER JOIN sys_role_permissions rp ON p.id = rp.permission_id
		INNER JOIN sys_user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = ? 
		AND p.dim_id = ?
		AND p.status = 1
	`
	var permItemIDs []uint
	err := s.db.Select(&permItemIDs, permissionQuery, userID, dim_id)
	if err != nil {
		return nil, fmt.Errorf("get user permissions failed: %v", err)
	}

	// 如果没有该维度的权限记录,直接返回空值
	if len(permItemIDs) == 0 {
		return []model.TreeDimensionItem{}, nil
	}

	// 检查是否有全部权限(item_id为0)
	hasFullAccess := false
	for _, itemID := range permItemIDs {
		if itemID == 0 {
			hasFullAccess = true
			break
		}
	}

	// 从sys_config_dimensions表读取维度表名
	var tableName string
	err = s.db.Get(&tableName, "SELECT table_name FROM sys_config_dimensions WHERE id = ?", dim_id)
	if err != nil {
		return nil, fmt.Errorf("get table name failed: %v", err)
	}

	// 获取表的列信息
	var columns []struct {
		ColumnName string `db:"COLUMN_NAME"`
	}
	err = s.db.Select(&columns, `
		SELECT COLUMN_NAME 
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE() 
		AND TABLE_NAME = ? 
		AND COLUMN_NAME NOT IN (
			'id', 'node_id', 'parent_id', 'name', 'code', 'description',
			'level', 'sort', 'status', 'created_at', 'creator_id', 
			'updated_at', 'updater_id'
		)
	`, tableName)
	if err != nil {
		return nil, fmt.Errorf("get columns failed: %v", err)
	}

	// 如果没有全部权限,需要获取有权限的节点的所有上下级节点
	var allowedIDs []uint
	if !hasFullAccess && len(permItemIDs) > 0 {
		// 获取所有有权限的item_id对应的node_id
		nodeIDQuery := fmt.Sprintf(`
			SELECT id, node_id FROM %s WHERE id IN (?)
		`, tableName)
		query, args, err := sqlx.In(nodeIDQuery, permItemIDs)
		if err != nil {
			return nil, fmt.Errorf("prepare node_id query failed: %v", err)
		}
		query = s.db.Rebind(query)

		type nodeInfo struct {
			ID     uint   `db:"id"`
			NodeID string `db:"node_id"`
		}
		var nodes []nodeInfo
		err = s.db.Select(&nodes, query, args...)
		if err != nil {
			return nil, fmt.Errorf("get node_ids failed: %v", err)
		}

		// 使用map去重
		idMap := make(map[uint]struct{})

		// 处理每个有权限的节点
		for _, node := range nodes {
			// 添加节点本身
			idMap[node.ID] = struct{}{}

			// 添加所有父节点
			idStrs := strings.Split(node.NodeID, "_")
			for _, idStr := range idStrs {
				idMap[utils.ParseUint(idStr)] = struct{}{}
			}

			// 获取所有子节点
			var childIDs []uint
			childQuery := fmt.Sprintf(`
				SELECT id FROM %s WHERE node_id LIKE ?
			`, tableName)
			err = s.db.Select(&childIDs, childQuery, node.NodeID+"_%")
			if err != nil {
				return nil, fmt.Errorf("get child ids failed: %v", err)
			}
			for _, id := range childIDs {
				idMap[id] = struct{}{}
			}
		}

		// 转换为切片
		for id := range idMap {
			allowedIDs = append(allowedIDs, id)
		}
	}

	// 构建查询
	var query strings.Builder
	query.WriteString(`
		SELECT DISTINCT id, node_id, parent_id, name, code, description, 
		level, sort, status, created_at, creator_id, updated_at, updater_id
	`)

	// 添加自定义列
	for _, col := range columns {
		query.WriteString(fmt.Sprintf(", %s", col.ColumnName))
	}

	query.WriteString(fmt.Sprintf(" FROM %s WHERE 1 = 1", tableName))
	args := []interface{}{}

	// 添加权限过滤条件
	if !hasFullAccess {
		if len(allowedIDs) == 0 {
			return []model.TreeDimensionItem{}, nil
		}
		query.WriteString(" AND id IN (?)")
		args = append(args, allowedIDs)
	}

	// 根据查询类型构建不同的查询条件
	switch query_type {
	case "children":
		query.WriteString(" AND parent_id = ?")
		args = append(args, id)
	case "leaves":
		if id != 0 {
			query.WriteString(" AND (node_id LIKE CONCAT(?,'_%') OR node_id LIKE CONCAT('%_' , ?,'_%'))")
			args = append(args, id)
		}
		query.WriteString(" AND NOT EXISTS (SELECT 1 FROM " + tableName + " b WHERE b.parent_id = " + tableName + ".id)")
	case "descendants":
		if id != 0 {
			// 获取节点的node_id
			var nodeID string
			err := s.db.Get(&nodeID, "SELECT node_id FROM sys_config_dimensions WHERE id = ?", id)
			if err != nil {
				return nil, fmt.Errorf("get node_id failed: %v", err)
			}
			query.WriteString(" AND node_id LIKE ?")
			args = append(args, nodeID+"_%")
		}
	}

	// 添加层级过滤
	if query_level != 0 {
		query.WriteString(" AND level = ?")
		args = append(args, query_level)
	}

	// 添加排序
	query.WriteString(" ORDER BY sort ASC, id ASC")

	// 执行查询
	rows, err := s.db.Queryx(query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query dimensions failed: %v", err)
	}
	defer rows.Close()

	var items []model.TreeDimensionItem
	for rows.Next() {
		var item model.TreeDimensionItem
		item.CustomData = make(map[string]string)

		// 创建map存储所有列的值
		data := make(map[string]interface{})
		err := rows.MapScan(data)
		if err != nil {
			return nil, fmt.Errorf("scan row failed: %v", err)
		}

		// 设置基础字段
		item.ID = utils.ParseUint(string(data["id"].([]byte)))
		item.NodeID = string(data["node_id"].([]byte))
		item.ParentID = utils.ParseUint(string(data["parent_id"].([]byte)))
		item.Name = string(data["name"].([]byte))
		item.Code = string(data["code"].([]byte))
		item.Description = string(data["description"].([]byte))
		item.Level = utils.ParseInt(string(data["level"].([]byte)))
		item.Sort = utils.ParseInt(string(data["sort"].([]byte)))
		item.Status = utils.ParseInt(string(data["status"].([]byte)))
		item.CreatedAt = utils.NewCustomTime(data["created_at"].(time.Time))
		item.CreatorID = utils.ParseUint(string(data["creator_id"].([]byte)))
		item.UpdatedAt = utils.NewCustomTime(data["updated_at"].(time.Time))
		item.UpdaterID = utils.ParseUint(string(data["updater_id"].([]byte)))

		// 设置自定义列的值
		for _, col := range columns {
			if val, ok := data[col.ColumnName]; ok {
				switch v := val.(type) {
				case []byte:
					item.CustomData[col.ColumnName] = string(v)
				case int64:
					item.CustomData[col.ColumnName] = fmt.Sprintf("%d", v)
				case string:
					item.CustomData[col.ColumnName] = v
				default:
					item.CustomData[col.ColumnName] = fmt.Sprintf("%v", v)
				}
			}
		}

		items = append(items, item)
	}
	if len(items) == 0 {
		return []model.TreeDimensionItem{}, nil
	}

	// 如果是children类型,直接返回结果
	if query_type == "children" {
		return items, nil
	}

	// 构建树形结构
	itemMap := make(map[uint]*model.TreeDimensionItem)
	for i := range items {
		item := &items[i]
		item.Children = []*model.TreeDimensionItem{}
		itemMap[item.ID] = item
	}

	var treeItems []*model.TreeDimensionItem
	for _, item := range itemMap {
		if item.ParentID == 0 || itemMap[item.ParentID] == nil {
			treeItems = append(treeItems, item)
		} else {
			parent := itemMap[item.ParentID]
			parent.Children = append(parent.Children, item)
		}
	}

	// 递归排序
	var sortTree func(items []*model.TreeDimensionItem)
	sortTree = func(items []*model.TreeDimensionItem) {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Sort < items[j].Sort
		})
		for _, item := range items {
			sortTree(item.Children)
		}
	}
	sortTree(treeItems)

	result := make([]model.TreeDimensionItem, len(treeItems))
	for i, item := range treeItems {
		result[i] = *item
	}

	return result, nil
}

// BatchDeleteDimensionItems 批量删除维度明细配置
func (s *DimensionService) DeleteDimensionItems(operatorID uint, dim_id uint, itemIDs []uint) error {
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
	query, args, err := sqlx.In(fmt.Sprintf("DELETE FROM %s WHERE id IN (?)", tableName), itemIDs)
	if err != nil {
		return fmt.Errorf("prepare query failed: %v", err)
	}

	query = tx.Rebind(query)
	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("delete dimension failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// UpdateDimensionItemSort 更新维度明细配置排序
func (s *DimensionService) UpdateDimensionItemSort(updaterID uint, dim_id uint, itemID uint, parentID uint, sort int) error {
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

	// 获取当前节点的node_id、parent_id和sort
	var oldNode struct {
		NodeID   string `db:"node_id"`
		ParentID uint   `db:"parent_id"`
		Sort     int    `db:"sort"`
	}
	err = tx.Get(&oldNode, fmt.Sprintf("SELECT node_id, parent_id, sort FROM %s WHERE id = ?", tableName), itemID)
	if err != nil {
		return fmt.Errorf("get node_id and sort failed: %v", err)
	}

	// 检查父节点是否变化
	if parentID != oldNode.ParentID {
		// 获取新的node_id
		var newNodeID string
		if parentID != 0 {
			var parent struct {
				NodeID string `db:"node_id"`
			}
			err = tx.Get(&parent, fmt.Sprintf("SELECT node_id FROM %s WHERE id = ?", tableName), parentID)
			if err != nil {
				return fmt.Errorf("get parent node_id failed: %v", err)
			}
			newNodeID = parent.NodeID + "_" + fmt.Sprint(itemID)
		} else {
			newNodeID = fmt.Sprint(itemID)
		}
		// 更新node_id
		_, err = tx.Exec(fmt.Sprintf("UPDATE %s SET node_id = ? WHERE id = ?", tableName), newNodeID, itemID)
		if err != nil {
			return fmt.Errorf("update node_id failed: %v", err)
		}
	}

	// 如果父节点变更，需要先从旧父节点移除，再插入到新父节点
	if parentID != oldNode.ParentID {
		// 获取新的node_id
		var newNodeID string
		if parentID != 0 {
			var parent struct {
				NodeID string `db:"node_id"`
			}
			err = tx.Get(&parent, fmt.Sprintf("SELECT node_id FROM %s WHERE id = ?", tableName), parentID)
			if err != nil {
				return fmt.Errorf("get parent node_id failed: %v", err)
			}
			newNodeID = parent.NodeID + "_" + fmt.Sprint(itemID)
		} else {
			newNodeID = fmt.Sprint(itemID)
		}

		// 1. 在旧父节点下移除该节点
		if oldNode.Sort != -1 {
			_, err = tx.Exec(fmt.Sprintf(`
                UPDATE %s 
                SET sort = sort - 1 
                WHERE parent_id = ? 
                AND sort > ?
            `, tableName), oldNode.ParentID, oldNode.Sort)
			if err != nil {
				return err
			}
		}

		// 2. 为新父节点的sort腾位置
		_, err = tx.Exec(fmt.Sprintf(`
            UPDATE %s
            SET sort = sort + 1
            WHERE parent_id = ?
            AND sort >= ?
        `, tableName), parentID, sort)
		if err != nil {
			return err
		}

		// 3. 更新该节点的父节点以及sort
		_, err = tx.Exec(fmt.Sprintf(`
            UPDATE %s
            SET parent_id = ?,
                sort = ?,
                node_id = ?
            WHERE id = ?
        `, tableName), parentID, sort, newNodeID, itemID)
		if err != nil {
			return err
		}
	} else {
		// 如果只是同一父节点内sort值变动，按原逻辑处理
		if sort != oldNode.Sort {
			_, err = tx.Exec(fmt.Sprintf("UPDATE %s SET sort = -1 WHERE id = ?", tableName), itemID)
			if err != nil {
				return err
			}
			if sort < oldNode.Sort {
				_, err = tx.Exec(fmt.Sprintf(`
                    UPDATE %s 
                    SET sort = sort + 1 
                    WHERE parent_id = ? 
                    AND sort >= ? 
                    AND sort < ?
                `, tableName), parentID, sort, oldNode.Sort)
			} else {
				_, err = tx.Exec(fmt.Sprintf(`
                    UPDATE %s 
                    SET sort = sort - 1 
                    WHERE parent_id = ? 
                    AND sort > ? 
                    AND sort <= ?
                `, tableName), parentID, oldNode.Sort, sort)
			}
			if err != nil {
				return err
			}
			_, err = tx.Exec(fmt.Sprintf("UPDATE %s SET sort = ? WHERE id = ?", tableName), sort, itemID)
			if err != nil {
				return err
			}
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}
