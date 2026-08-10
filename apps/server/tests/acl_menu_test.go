//go:build integration

package tests

import (
	"testing"

	"github.com/liwook/go-vue-selection/pkg/idconv"
	"github.com/liwook/go-vue-selection/types"
)

// TestMenuFlow 覆盖菜单写链路：save → doAssign(给普通用户角色) → update → list 验证 → delete。
func TestMenuFlow(t *testing.T) {
	const (
		menuName = "it_menu_tmp"
		menuCode = "it:menu:tmp"
	)

	// 1. save（一级目录，挂在根下 parentId=0）
	apiClient.Post(t, "/api/v1/acl/permission", types.ParamMenuSave{
		Name:     menuName,
		ParentID: "0",
		CODE:     menuCode,
		TYPE:     types.MenuTypeDir,
		LEVEL:    types.MenuLevel1,
	})

	// 2. 从菜单树反查 menuId
	menuID := findMenuIDByCode(t, menuCode)
	if menuID == "" {
		t.Fatalf("created menu %q not found in tree", menuCode)
	}

	// 3. doAssign：把该菜单分配给普通用户角色(role_id=2)
	apiClient.Post(t, "/api/v1/acl/permission/role/2", []int64{idconv.ToInt64Safe(menuID)})

	// 4. update
	apiClient.Put(t, "/api/v1/acl/permission/"+menuID, types.ParamMenuUpdate{
		MenuID:   menuID,
		Name:     menuName,
		ParentID: "0",
		CODE:     menuCode,
		LEVEL:    types.MenuLevel1,
	})

	// 5. list 验证（更新后仍在树中）
	if got := findMenuIDByCode(t, menuCode); got != menuID {
		t.Fatalf("menu %q missing after update, got=%q want=%q", menuCode, got, menuID)
	}

	// 6. delete（同时作为清理）
	apiClient.Delete(t, "/api/v1/acl/permission/"+menuID)

	// 7. 删除后不再出现
	if got := findMenuIDByCode(t, menuCode); got != "" {
		t.Fatalf("menu %q should be deleted, but still found", menuCode)
	}
}

// findMenuIDByCode 在菜单树里按 code 递归查找 menuId。
func findMenuIDByCode(t *testing.T, code string) string {
	t.Helper()
	resp := apiClient.Get(t, "/api/v1/acl/permission")
	var menus []types.Menu
	resp.decodeData(&menus)
	return walkMenu(menus, code)
}

func walkMenu(menus []types.Menu, code string) string {
	for i := range menus {
		if menus[i].CODE == code {
			return menus[i].MenuID
		}
		if id := walkMenu(menus[i].CHILDREN, code); id != "" {
			return id
		}
	}
	return ""
}
