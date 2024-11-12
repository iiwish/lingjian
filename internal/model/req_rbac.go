package model

// 修改角色权限请求详情
type PatchRolePermDetail struct {
	Op    string `json:"op" binding:"required,oneof=add remove"`
	Value []uint `json:"value" binding:"required"`
}

// 修改角色权限请求
type PatchRolePerms []PatchRolePermDetail
