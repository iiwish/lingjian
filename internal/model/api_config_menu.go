package model

type CreateMenuReq struct {
	ParentID    uint   `json:"parent_id"`   // 父节点ID
	MenuName    string `json:"menu_name"`   // 名称
	TableName   string `json:"table_name"`  // 数据表名
	Description string `json:"description"` // 描述
}

type UpdateMenuReq struct {
	MenuName    string `json:"menu_name"`   // 名称
	TableName   string `json:"table_name"`  // 数据表名
	Description string `json:"description"` // 描述
}
