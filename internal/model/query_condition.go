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

// LogicOperator 逻辑运算符
type LogicOperator string

const (
	LogicAnd LogicOperator = "AND"
	LogicOr  LogicOperator = "OR"
)

// Condition 查询条件
type Condition struct {
	Field    string      `json:"field"`    // 字段名
	Operator Operator    `json:"operator"` // 操作符
	Value    interface{} `json:"value"`    // 查询值
}

// ConditionGroup 条件组
type ConditionGroup struct {
	Logic      LogicOperator `json:"logic"`      // 逻辑运算符
	Conditions []interface{} `json:"conditions"` // 可以是Condition或ConditionGroup
}

// QueryCondition 查询条件组合
type QueryCondition struct {
	Root    ConditionGroup `json:"root"`     // 根条件组
	OrderBy []OrderBy      `json:"order_by"` // 排序
	GroupBy []string       `json:"group_by"` // 分组
}

// OrderBy 排序
type OrderBy struct {
	Field string `json:"field"` // 排序字段
	Desc  bool   `json:"desc"`  // 是否降序
}

// BuildQuery 构建查询SQL
func (q *QueryCondition) BuildQuery(tableName string) (string, []interface{}) {
	query := "SELECT * FROM " + tableName
	where, args := buildWhereClause(&q.Root)
	if where != "" {
		query += " WHERE " + where
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

func buildWhereClause(group *ConditionGroup) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	for _, item := range group.Conditions {
		switch v := item.(type) {
		case Condition:
			condition, arg := buildCondition(v)
			conditions = append(conditions, condition)
			args = append(args, arg...)
		case ConditionGroup:
			condition, arg := buildWhereClause(&v)
			conditions = append(conditions, "("+condition+")")
			args = append(args, arg...)
		case map[string]interface{}:
			// 将map转换为Condition结构体
			c := Condition{
				Field:    v["field"].(string),
				Operator: Operator(v["operator"].(string)),
				Value:    v["value"],
			}
			condition, arg := buildCondition(c)
			conditions = append(conditions, condition)
			args = append(args, arg...)
		}
	}

	return strings.Join(conditions, " "+string(group.Logic)+" "), args
}

func buildCondition(condition Condition) (string, []interface{}) {
	var query string
	var args []interface{}

	switch condition.Operator {
	case OpEq:
		query = condition.Field + " = ?"
		args = append(args, condition.Value)
	case OpNe:
		query = condition.Field + " != ?"
		args = append(args, condition.Value)
	case OpGt:
		query = condition.Field + " > ?"
		args = append(args, condition.Value)
	case OpGte:
		query = condition.Field + " >= ?"
		args = append(args, condition.Value)
	case OpLt:
		query = condition.Field + " < ?"
		args = append(args, condition.Value)
	case OpLte:
		query = condition.Field + " <= ?"
		args = append(args, condition.Value)
	case OpLike:
		query = condition.Field + " LIKE ?"
		args = append(args, "%"+condition.Value.(string)+"%")
	case OpNotLike:
		query = condition.Field + " NOT LIKE ?"
		args = append(args, "%"+condition.Value.(string)+"%")
	case OpIn:
		values := condition.Value.([]interface{})
		placeholders := make([]string, len(values))
		for i := range values {
			placeholders[i] = "?"
			args = append(args, values[i])
		}
		query = condition.Field + " IN (" + strings.Join(placeholders, ",") + ")"
	case OpNotIn:
		values := condition.Value.([]interface{})
		placeholders := make([]string, len(values))
		for i := range values {
			placeholders[i] = "?"
			args = append(args, values[i])
		}
		query = condition.Field + " NOT IN (" + strings.Join(placeholders, ",") + ")"
	case OpBetween:
		values := condition.Value.([]interface{})
		query = condition.Field + " BETWEEN ? AND ?"
		args = append(args, values[0], values[1])
	case OpNotBetween:
		values := condition.Value.([]interface{})
		query = condition.Field + " NOT BETWEEN ? AND ?"
		args = append(args, values[0], values[1])
	}

	return query, args
}
