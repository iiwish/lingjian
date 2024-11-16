package model

import (
	"strings"
)

// Operator 查询操作符
type Operator string

const (
	OpEq         Operator = "eq"          // 等于
	OpNe         Operator = "ne"          // 不等于
	OpGt         Operator = "gt"          // 大于
	OpGte        Operator = "gte"         // 大于等于
	OpLt         Operator = "lt"          // 小于
	OpLte        Operator = "lte"         // 小于等于
	OpLike       Operator = "like"        // 模糊查询
	OpNotLike    Operator = "not_like"    // 不包含
	OpIn         Operator = "in"          // IN查询
	OpNotIn      Operator = "not_in"      // NOT IN查询
	OpBetween    Operator = "between"     // 区间查询
	OpNotBetween Operator = "not_between" // 不在区间
)

// Condition 查询条件
type Condition struct {
	Field    string      `json:"field"`    // 字段名
	Operator Operator    `json:"operator"` // 操作符
	Value    interface{} `json:"value"`    // 查询值
}

// QueryCondition 查询条件组合
type QueryCondition struct {
	Conditions []Condition `json:"conditions"` // 条件列表
	OrderBy    []OrderBy   `json:"order_by"`   // 排序
	GroupBy    []string    `json:"group_by"`   // 分组
}

// OrderBy 排序
type OrderBy struct {
	Field string `json:"field"` // 排序字段
	Desc  bool   `json:"desc"`  // 是否降序
}

// BuildQuery 构建查询SQL
func (q *QueryCondition) BuildQuery(tableName string) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	// 构建WHERE条件
	for _, condition := range q.Conditions {
		switch condition.Operator {
		case OpEq:
			conditions = append(conditions, condition.Field+" = ?")
			args = append(args, condition.Value)
		case OpNe:
			conditions = append(conditions, condition.Field+" != ?")
			args = append(args, condition.Value)
		case OpGt:
			conditions = append(conditions, condition.Field+" > ?")
			args = append(args, condition.Value)
		case OpGte:
			conditions = append(conditions, condition.Field+" >= ?")
			args = append(args, condition.Value)
		case OpLt:
			conditions = append(conditions, condition.Field+" < ?")
			args = append(args, condition.Value)
		case OpLte:
			conditions = append(conditions, condition.Field+" <= ?")
			args = append(args, condition.Value)
		case OpLike:
			conditions = append(conditions, condition.Field+" LIKE ?")
			args = append(args, "%"+condition.Value.(string)+"%")
		case OpNotLike:
			conditions = append(conditions, condition.Field+" NOT LIKE ?")
			args = append(args, "%"+condition.Value.(string)+"%")
		case OpIn:
			values := condition.Value.([]interface{})
			placeholders := make([]string, len(values))
			for i := range values {
				placeholders[i] = "?"
				args = append(args, values[i])
			}
			conditions = append(conditions, condition.Field+" IN ("+strings.Join(placeholders, ",")+")")
		case OpNotIn:
			values := condition.Value.([]interface{})
			placeholders := make([]string, len(values))
			for i := range values {
				placeholders[i] = "?"
				args = append(args, values[i])
			}
			conditions = append(conditions, condition.Field+" NOT IN ("+strings.Join(placeholders, ",")+")")
		case OpBetween:
			values := condition.Value.([]interface{})
			conditions = append(conditions, condition.Field+" BETWEEN ? AND ?")
			args = append(args, values[0], values[1])
		case OpNotBetween:
			values := condition.Value.([]interface{})
			conditions = append(conditions, condition.Field+" NOT BETWEEN ? AND ?")
			args = append(args, values[0], values[1])
		}
	}

	// 构建查询SQL
	query := "SELECT * FROM " + tableName
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// 添加GROUP BY
	if len(q.GroupBy) > 0 {
		query += " GROUP BY " + strings.Join(q.GroupBy, ",")
	}

	// 添加ORDER BY
	if len(q.OrderBy) > 0 {
		var orders []string
		for _, order := range q.OrderBy {
			if order.Desc {
				orders = append(orders, order.Field+" DESC")
			} else {
				orders = append(orders, order.Field+" ASC")
			}
		}
		query += " ORDER BY " + strings.Join(orders, ",")
	}

	return query, args
}
