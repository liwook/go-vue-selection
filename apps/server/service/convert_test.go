package service

import (
	"testing"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/types"
)

func strPtr(s string) *string { return &s }
func int64Ptr(i int64) *int64 { return &i }

func TestDerefStr(t *testing.T) {
	if got := derefStr(nil); got != "" {
		t.Errorf("derefStr(nil) = %q, want empty", got)
	}
	if got := derefStr(strPtr("hi")); got != "hi" {
		t.Errorf("derefStr = %q, want hi", got)
	}
}

func TestDerefInt64(t *testing.T) {
	if got := derefInt64(nil); got != 0 {
		t.Errorf("derefInt64(nil) = %d, want 0", got)
	}
	if got := derefInt64(int64Ptr(7)); got != 7 {
		t.Errorf("derefInt64 = %d, want 7", got)
	}
}

func TestPageCount(t *testing.T) {
	tests := []struct {
		total, size, want int64
	}{
		{0, 10, 0},
		{100, 10, 10},
		{101, 10, 11},
		{95, 10, 10},
		{10, 0, 0},  // size<=0 视为 0 页
		{10, -1, 0}, // 负 size 视为 0 页
	}
	for _, tt := range tests {
		if got := pageCount(tt.total, tt.size); got != tt.want {
			t.Errorf("pageCount(%d,%d) = %d, want %d", tt.total, tt.size, got, tt.want)
		}
	}
}

func TestModelToCategory1(t *testing.T) {
	r := &model.Category1{Category1ID: 11, Name: "电子"}
	got := modelToCategory1(r)
	if got.CategoryID != "11" || got.Name != "电子" {
		t.Errorf("got = %+v", got)
	}
}

func TestModelToCategory2(t *testing.T) {
	r := &model.Category2{Category2ID: 22, Category1ID: 11, Name: "手机"}
	got := modelToCategory2(r)
	if got.Category2ID != "22" || got.Category1ID != "11" || got.Name != "手机" {
		t.Errorf("got = %+v", got)
	}
}

func TestModelToCategory3(t *testing.T) {
	r := &model.Category3{Category3ID: 33, Category2ID: 22, Name: "智能手机"}
	got := modelToCategory3(r)
	if got.Category3ID != "33" || got.Category2ID != "22" || got.Name != "智能手机" {
		t.Errorf("got = %+v", got)
	}
}

func TestModelToMenu(t *testing.T) {
	pid := int64(5)
	r := &model.Menu{
		MenuID: 10, Pid: &pid, Name: "系统管理",
		Code: strPtr("sys"), ToCode: strPtr("tc"),
		Type: 1, Status: true, Level: 2,
	}
	got := modelToMenu(r)
	if got.MenuID != "10" {
		t.Errorf("MenuID = %q, want 10", got.MenuID)
	}
	if got.ParentID != "5" {
		t.Errorf("ParentID = %q, want 5 (from *int64)", got.ParentID)
	}
	if got.CODE != "sys" || got.TOCODE != "tc" {
		t.Errorf("CODE/TOCODE = %q/%q", got.CODE, got.TOCODE)
	}
	if got.TYPE != 1 || !got.STATUS || got.LEVEL != 2 {
		t.Errorf("TYPE/STATUS/LEVEL = %d/%v/%d", got.TYPE, got.STATUS, got.LEVEL)
	}
	if got.SELECT {
		t.Error("SELECT should default false")
	}
}

func TestModelToMenu_NilPid(t *testing.T) {
	r := &model.Menu{MenuID: 1, Pid: nil, Name: "root"}
	got := modelToMenu(r)
	if got.ParentID != "" {
		t.Errorf("ParentID = %q, want empty string for nil pid", got.ParentID)
	}
}

func TestModelToMenuList(t *testing.T) {
	in := []*model.Menu{
		{MenuID: 1, Name: "a"},
		{MenuID: 2, Name: "b"},
	}
	got := modelToMenuList(in)
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].MenuID != "1" || got[1].MenuID != "2" {
		t.Errorf("got = %+v", got)
	}
}

func TestModelToRole(t *testing.T) {
	r := &model.Role{RoleID: 7, RoleName: "admin", Remark: strPtr("超级管理员")}
	got := modelToRole(r)
	if got.RoleID != "7" || got.RoleName != "admin" || got.Remark != "超级管理员" {
		t.Errorf("got = %+v", got)
	}
}

func TestModelToRole_NilRemark(t *testing.T) {
	r := &model.Role{RoleID: 8, RoleName: "guest", Remark: nil}
	got := modelToRole(r)
	if got.Remark != "" {
		t.Errorf("Remark = %q, want empty", got.Remark)
	}
}

func TestModelToUser(t *testing.T) {
	phone := "13800000000"
	avatar := "http://x/a.png"
	r := &model.User{
		UserID: 100, Username: "u1", Password: "pwd", Name: strPtr("用户一"),
		Phone: &phone, Avatar: &avatar, Status: true,
	}
	got := modelToUser(r)
	if got.UserID != "100" || got.Username != "u1" || got.Password != "pwd" {
		t.Errorf("base fields got = %+v", got)
	}
	if got.Name != "用户一" {
		t.Errorf("Name = %q, want 用户一", got.Name)
	}
	if got.Phone != phone {
		t.Errorf("Phone = %q, want %q", got.Phone, phone)
	}
	if got.Avatar != avatar {
		t.Errorf("Avatar = %q, want %q", got.Avatar, avatar)
	}
	if !got.Status {
		t.Error("Status should be true")
	}
}

func TestModelToUser_NilOptionals(t *testing.T) {
	r := &model.User{UserID: 101, Username: "u2", Name: nil}
	got := modelToUser(r)
	if got.Name != "" {
		t.Errorf("Name = %q, want empty for nil", got.Name)
	}
	if got.Phone != "" || got.Avatar != "" {
		t.Errorf("Phone/Avatar = %q/%q, want empty for nil", got.Phone, got.Avatar)
	}
}

func TestModelToResponseUser(t *testing.T) {
	phone := "13900000000"
	r := &model.User{
		UserID: 200, Username: "u3", Name: strPtr("用户三"), Phone: &phone, Status: false,
	}
	got := modelToResponseUser(r)
	want := types.ResponseUser{
		UserID: "200", Username: "u3", Name: "用户三", Phone: phone, Status: false,
	}
	if *got != want {
		t.Errorf("got = %+v, want %+v", *got, want)
	}
}

func TestModelToTrademark(t *testing.T) {
	r := &model.Trademark{TmID: 50, TmName: "Nike", LogoURL: strPtr("http://logo/n.png")}
	got := modelToTrademark(r)
	if got.TmID != "50" || got.TmName != "Nike" || got.LogoUrl != "http://logo/n.png" {
		t.Errorf("got = %+v", got)
	}
}

func TestModelToTrademark_NilLogo(t *testing.T) {
	r := &model.Trademark{TmID: 51, TmName: "Adidas", LogoURL: nil}
	got := modelToTrademark(r)
	if got.LogoUrl != "" {
		t.Errorf("LogoUrl = %q, want empty for nil", got.LogoUrl)
	}
}

func TestModelToAttr(t *testing.T) {
	r := &model.Attr{AttrID: 60, AttrName: "颜色", CategoryID: 33}
	got := modelToAttr(r)
	if got.AttrID != "60" || got.AttrName != "颜色" || got.CategoryID != "33" {
		t.Errorf("got = %+v", got)
	}
}

func TestModelToAttrValue(t *testing.T) {
	r := &model.AttrValue{AttrValueID: 70, ValueName: "红", AttrID: 60}
	got := modelToAttrValue(r)
	if got.AttrValueID != "70" || got.ValueName != "红" || got.AttrID != "60" {
		t.Errorf("got = %+v", got)
	}
}
