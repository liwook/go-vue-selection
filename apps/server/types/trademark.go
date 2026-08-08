package types

type Trademark struct {
	BaseModel

	TmID    string `json:"tmId"`
	TmName  string `json:"tmName"`
	LogoUrl string `json:"logoUrl"`
}

// ParamTmSave 品牌创建参数
type ParamTmSave struct {
	TmName  string `json:"tmName" binding:"required"`
	LogoUrl string `json:"logoUrl" binding:"required"`
}

// ParamTmUpdate 品牌更新参数
type ParamTmUpdate struct {
	TmID    string `json:"tmId"`
	TmName  string `json:"tmName" binding:"required"`
	LogoUrl string `json:"logoUrl" binding:"required"`
}

// ResponseTmList 返回品牌列表
type ResponseTmList struct {
	Records []*Trademark `json:"records"`

	Total       int64 `json:"total"`
	Size        int64 `json:"size"`
	Current     int64 `json:"current"`
	SearchCount bool  `json:"searchCount"`
	Pages       int64 `json:"pages"`
}
