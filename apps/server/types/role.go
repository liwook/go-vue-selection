package types

type Role struct {
	BaseModel

	RoleID   string `json:"roleId"`
	RoleName string `json:"roleName"`
	Remark   string `json:"remark"`
}

type ParamRoleSave struct {
	RoleName string `json:"roleName"`         // 角色名称，不能为空
	Remark   string `json:"remark,omitempty"` // 角色备注
}

type ParamRoleUpdate struct {
	RoleID   string `json:"roleId"`           // 角色ID，
	RoleName string `json:"roleName"`         // 角色名称，不能为空
	Remark   string `json:"remark,omitempty"` // 角色备注
}

// ResponseRoleList 角色列表返回数据
type ResponseRoleList struct {
	Records     []*Role `json:"records"`
	Total       int64   `json:"total"`
	Size        int64   `json:"size"`
	Current     int64   `json:"current"`
	Pages       int64   `json:"pages"`
	SearchCount bool    `json:"searchCount"`
}

type ResponseToAssignRole struct {
	AssignRoles  []*Role `json:"assignRoles"`
	AllRolesList []*Role `json:"allRolesList"`
}
