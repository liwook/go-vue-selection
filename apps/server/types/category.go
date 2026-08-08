package types

type Category1 struct {
	BaseModel

	CategoryID string `json:"category1Id"`
	Name       string `json:"name"`
}

type Category2 struct {
	BaseModel

	Category2ID string `json:"category2Id"`
	Name        string `json:"name"`
	Category1ID string `json:"category1Id"`
}

type Category3 struct {
	BaseModel

	Category3ID string `json:"category3Id"`
	Name        string `json:"name"`
	Category2ID string `json:"category2Id"`
}

// ParamC2Create 新增二级分类入参：前端只传所属一级分类 ID 与名称，
// 二级分类 ID 由 DB 自增生成，不应由前端指定。
type ParamC2Create struct {
	Name        string `json:"name" binding:"required"`
	Category1ID string `json:"category1Id" binding:"required"`
}

// ParamC3Create 新增三级分类入参：前端只传所属二级分类 ID 与名称，
// 三级分类 ID 由 DB 自增生成，不应由前端指定。
type ParamC3Create struct {
	Name        string `json:"name" binding:"required"`
	Category2ID string `json:"category2Id" binding:"required"`
}
