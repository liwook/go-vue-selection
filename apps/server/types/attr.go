package types

type Attr struct {
	BaseModel

	AttrID        string       `json:"attrId"`                        // 属性ID
	AttrName      string       `json:"attrName" binding:"required"`   // 属性名称
	CategoryID    string       `json:"categoryId" binding:"required"` // 三级分类ID
	AttrValueList []*AttrValue `json:"attrValueList"`                 // 属性值列表
}

type AttrValue struct {
	AttrValueID string `json:"attrValueId"`
	ValueName   string `json:"valueName" binding:"required"`
	AttrID      string `json:"attrId"`
}

type ParamAttrCreate struct {
	AttrID        string                  `json:"attrId"`
	AttrName      string                  `json:"attrName" binding:"required"`
	CategoryID    string                  `json:"categoryId" binding:"required"`
	AttrValueList []*ParamAttrValueCreate `json:"attrValueList"` // 属性值列表
}

type ParamAttrValueCreate struct {
	ValueName string `json:"valueName" binding:"required"`
}

//func GetCommunityList() (communityList []*models.Community, err error) {
//	sqlStr := `SELECT community_id, community_name FROM community`
//	if err = db.Select(&communityList, sqlStr); err != nil {
//		if err == sql.ErrNoRows {
//			slog.Warn("there is no community in db")
//			return nil, err
//		}
//	}
//
//	return
//}
