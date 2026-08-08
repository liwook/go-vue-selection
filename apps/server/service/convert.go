package service

import (
	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/pkg/idconv"
	"github.com/liwook/go-vue-selection/types"
)

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func pageCount(total, size int64) int64 {
	if size <= 0 {
		return 0
	}
	pages := total / size
	if total%size != 0 {
		pages++
	}
	return pages
}

// modelToCategory1 将数据库实体转换为 API 层分类类型
func modelToCategory1(r *model.Category1) *types.Category1 {
	return &types.Category1{
		CategoryID: idconv.ToStr(r.Category1ID),
		Name:       r.Name,
		BaseModel: types.BaseModel{
			CreateTime: r.CreateTime,
			UpdateTime: r.UpdateTime,
		},
	}
}

func modelToCategory2(r *model.Category2) *types.Category2 {
	return &types.Category2{
		Category2ID: idconv.ToStr(r.Category2ID),
		Category1ID: idconv.ToStr(r.Category1ID),
		Name:        r.Name,
		BaseModel: types.BaseModel{
			CreateTime: r.CreateTime,
			UpdateTime: r.UpdateTime,
		},
	}
}

func modelToCategory3(r *model.Category3) *types.Category3 {
	return &types.Category3{
		Category3ID: idconv.ToStr(r.Category3ID),
		Category2ID: idconv.ToStr(r.Category2ID),
		Name:        r.Name,
		BaseModel: types.BaseModel{
			CreateTime: r.CreateTime,
			UpdateTime: r.UpdateTime,
		},
	}
}

// modelToMenu 将数据库菜单实体转换为 API 层菜单类型（含树形字段初始化）
func modelToMenu(m *model.Menu) *types.Menu {
	return &types.Menu{
		MenuID:   idconv.ToStr(m.MenuID),
		ParentID: idconv.ToStrPtr(m.Pid),
		Name:     m.Name,
		CODE:     derefStr(m.Code),
		TOCODE:   derefStr(m.ToCode),
		STATUS:   m.Status,
		LEVEL:    int(m.Level),
		TYPE:     int(m.Type),
		SELECT:   false,
		BaseModel: types.BaseModel{
			CreateTime: m.CreateTime,
			UpdateTime: m.UpdateTime,
		},
	}
}

func modelToMenuList(list []*model.Menu) []*types.Menu {
	res := make([]*types.Menu, 0, len(list))
	for _, m := range list {
		res = append(res, modelToMenu(m))
	}
	return res
}

func modelToRole(r *model.Role) *types.Role {
	return &types.Role{
		RoleID:   idconv.ToStr(r.RoleID),
		RoleName: r.RoleName,
		Remark:   derefStr(r.Remark),
		BaseModel: types.BaseModel{
			CreateTime: r.CreateTime,
			UpdateTime: r.UpdateTime,
		},
	}
}

func modelToUser(u *model.User) *types.User {
	m := &types.User{
		UserID:   idconv.ToStr(u.UserID),
		Username: u.Username,
		Password: u.Password,
		Name:     derefStr(u.Name),
		BaseModel: types.BaseModel{
			CreateTime: u.CreateTime,
			UpdateTime: u.UpdateTime,
		},
	}
	if u.Phone != nil {
		m.Phone = *u.Phone
	}
	if u.Avatar != nil {
		m.Avatar = *u.Avatar
	}
	m.Status = u.Status
	return m
}

func modelToResponseUser(u *model.User) *types.ResponseUser {
	m := modelToUser(u)
	return &types.ResponseUser{
		UserID:   m.UserID,
		Username: m.Username,
		Name:     m.Name,
		Phone:    m.Phone,
		Status:   m.Status,
	}
}

func modelToTrademark(t *model.Trademark) *types.Trademark {
	return &types.Trademark{
		TmID:    idconv.ToStr(t.TmID),
		TmName:  t.TmName,
		LogoUrl: derefStr(t.LogoURL),
		BaseModel: types.BaseModel{
			CreateTime: t.CreateTime,
			UpdateTime: t.UpdateTime,
		},
	}
}

func modelToAttr(a *model.Attr) *types.Attr {
	return &types.Attr{
		AttrID:     idconv.ToStr(a.AttrID),
		AttrName:   a.AttrName,
		CategoryID: idconv.ToStr(a.CategoryID),
		BaseModel: types.BaseModel{
			CreateTime: a.CreateTime,
			UpdateTime: a.UpdateTime,
		},
	}
}

func modelToAttrValue(v *model.AttrValue) *types.AttrValue {
	return &types.AttrValue{
		AttrValueID: idconv.ToStr(v.AttrValueID),
		ValueName:   v.ValueName,
		AttrID:      idconv.ToStr(v.AttrID),
	}
}
