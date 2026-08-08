package types

type User struct {
	BaseModel

	UserID   string `json:"userId"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar,omitempty"`
	Status   bool   `json:"status"`
}

// ParamUserLogin 用户登录参数
type ParamUserLogin struct {
	Username string `json:"username" binding:"required"` // 用户名，不能为空
	Password string `json:"password" binding:"required"` // 密码，不能为空
}

// ParamUserSignUp 注册用户参数
type ParamUserSignUp struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamUserUpdate 更新用户参数（部分更新：只更新前端显式发送的非空字段，
// 未发送（nil）的字段保持数据库原值不变）。user_id 由 JWT 令牌解析得到，无需前端传入。
type ParamUserUpdate struct {
	Name   *string `json:"name,omitempty"`   // 用户昵称
	Avatar *string `json:"avatar,omitempty"` // 头像 URL
}

// ResponseUserInfo 用户信息返回数据
type ResponseUserInfo struct {
	Routes  []string `json:"routes"`
	Buttons []string `json:"buttons"`
	Roles   []string `json:"roles"`
	Name    string   `json:"name"`
	Avatar  string   `json:"avatar"`
}

// ResponseUser 用户返回数据
type ResponseUser struct {
	BaseModel
	UserID   string `json:"userId"`
	Username string `json:"username"`
	RoleName string `json:"roleName"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Status   bool   `json:"status"`
}

// ResponseUserList 用户列表返回数据
type ResponseUserList struct {
	Records []*ResponseUser `json:"records"`
	Total   int64           `json:"total"`
	Size    int64           `json:"size"`
	Current int64           `json:"current"`
	Pages   int64           `json:"pages"`
}

type ParamDoAssignRole struct {
	UserID     string   `json:"userId"`
	RoleIDList []string `json:"roleIdList"`
}

// ParamUserLock 锁定/解锁用户参数
type ParamUserLock struct {
	UserID string `json:"userId" binding:"required"`
	Status bool   `json:"status"` // true=解锁（恢复正常），false=锁定（禁用）
}
