package model

// CreateModelReq 创建模型请求参数
type CreateModelReq struct {
	ParentID      uint            `json:"parent_id"`
	ModelName     string          `json:"model_name"`
	DisplayName   string          `json:"display_name"`
	Description   string          `json:"description"`
	Configuration ModelConfigItem `json:"configuration"`
	Status        int             `json:"status"` // 0:禁用 1:启用
}

// UpdateModelReq 更新模型请求参数
type UpdateModelReq struct {
	ID            uint            `json:"id"`
	ModelName     string          `json:"model_name"`
	DisplayName   string          `json:"display_name"`
	Description   string          `json:"description"`
	Configuration ModelConfigItem `json:"configuration"`
	Status        int             `json:"status"` // 0:禁用 1:启用
}

// ModelResp 数据模型
type ModelResp struct {
	ID            uint            `json:"id"`
	ModelName     string          `json:"model_name"`
	DisplayName   string          `json:"display_name"`
	Description   string          `json:"description"`
	Configuration ModelConfigItem `json:"configuration"`
	Status        int             `json:"status"`
}

// ModelConfigItemDim 数据模型配置项维度
type ModelConfigItemDim struct {
	TableField string `json:"table_field"`
	DimField   string `json:"dim_field"`
	DimID      uint   `json:"dim_id"`
	ItemID     uint   `json:"item_id"`
	Type       string `json:"type"`
}

// ModelConfigItemRelField 定义关系字段映射
type ModelConfigItemRelField struct {
	FromField string `json:"fromField"`
	ToField   string `json:"toField"`
}

// ModelConfigItemRel 数据模型配置项关联
type ModelConfigItemRel struct {
	Type   string                    `json:"type"`
	Fields []ModelConfigItemRelField `json:"fields"`
}

// ModelConfigItem 数据模型配置项
type ModelConfigItem struct {
	TableID       uint                 `json:"table_id"`
	Dimensions    []ModelConfigItemDim `json:"dimensions"`
	Relationships ModelConfigItemRel   `json:"relationships"`
	Childrens     []ModelConfigItem    `json:"childrens"`
}
