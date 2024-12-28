package element

import (
	"fmt"
	"sort"
	"strings"

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

// BatchCreateDimensionItems 批量创建维度明细配置
func (s *DimensionService) CreateDimensionItems(dimensionItems []*model.DimensionItem, creatorID uint, dim_id uint) error {
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
		// 获取父节点的node_id
		if dimension.ParentID != 0 {
			var parent struct {
				NodeID string `db:"node_id"`
				Level  int    `db:"level"`
			}
			err = tx.Get(&parent, fmt.Sprintf("SELECT node_id, level FROM %s WHERE id = ?", table_name), dimension.ParentID)
			if err != nil {
				return fmt.Errorf("get parent node_id and level failed: %v", err)
			}
			dimension.NodeID = parent.NodeID
			dimension.Level = parent.Level + 1
		} else {
			dimension.NodeID = ""
			dimension.Level = 1
		}

		// 获取sort数值
		var sort int
		err = tx.Get(&sort, fmt.Sprintf("SELECT IFNULL(MAX(sort), 0) FROM %s WHERE parent_id = ?", table_name), dimension.ParentID)
		if err != nil {
			return fmt.Errorf("get max sort failed: %v", err)
		}
		dimension.Sort = sort + 1

		result, err := tx.NamedExec(fmt.Sprintf(`
			INSERT INTO %s (node_id, parent_id, name, code, description, level, sort, status, created_at, creator_id, updated_at, updater_id)
			VALUES (:node_id, :parent_id, :name, :code, :description, :level, :sort, :status, Now(), :creator_id, Now(), :updater_id)
		`, table_name), dimension)
		if err != nil {
			return fmt.Errorf("insert sys_config_dimensions failed: %v", err)
		}

		// 获取插入的ID
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("get last insert id failed: %v", err)
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
			return fmt.Errorf("update node_id failed: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
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

		// 更新维度配置
		_, err = tx.NamedExec(fmt.Sprintf(`
			UPDATE %s SET name = :name, code = :code, description = :description, status = :status, updated_at = Now(), updater_id = :updater_id, custom1 = :custom1, custom2 = :custom2, custom3 = :custom3
			WHERE id = :id
		`, table_name), dimension)
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
func (s *DimensionService) TreeDimensionItems(dim_id uint, id uint, query_type string, query_level uint) ([]model.TreeDimensionItem, error) {
	// 从sys_config_dimensions表读取维度表名
	var tableName string
	err := s.db.Get(&tableName, "SELECT table_name FROM sys_config_dimensions WHERE id = ?", dim_id)
	if err != nil {
		return nil, fmt.Errorf("get table name failed: %v", err)
	}

	// 构建查询
	var query strings.Builder
	query.WriteString(fmt.Sprintf("SELECT * FROM %s WHERE 1 = 1", tableName))
	args := []interface{}{}

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
	var items []model.TreeDimensionItem
	err = s.db.Select(&items, query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("list dimensions failed: %v", err)
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
