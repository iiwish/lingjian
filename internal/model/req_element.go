package model

// UpdateTableItemsRequest 更新数据表记录请求参数
type UpdateTableItemsRequest struct {
	PrimaryKeyColumns []string                 `json:"primary_key_columns"` // 主键列名列表
	Items             []map[string]interface{} `json:"items"`               // 要更新的数据表记录
}

// DeleteTableItemsRequest 批量删除数据表记录请求参数
type DeleteTableItemsRequest struct {
	Conditions []map[string]interface{} `json:"conditions"` // 删除条件
}
