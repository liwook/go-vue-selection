//go:build integration

package tests

import (
	"testing"

	"github.com/liwook/go-vue-selection/types"
)

// TestRoleFlow 覆盖角色写链路：save → update → list 验证 → delete。
func TestRoleFlow(t *testing.T) {
	const roleName = "it_role_tmp"

	// 1. save
	apiClient.Post(t, "/api/v1/acl/role", types.ParamRoleSave{
		RoleName: roleName,
		Remark:   "integration-test",
	})

	// 2. 从 list 反查新角色的 roleId
	roleID := findRoleIDByName(t, roleName)
	if roleID == "" {
		t.Fatalf("created role %q not found in list", roleName)
	}

	// 3. update
	apiClient.Put(t, "/api/v1/acl/role/"+roleID, types.ParamRoleUpdate{
		RoleID:   roleID,
		RoleName: roleName,
		Remark:   "integration-test-updated",
	})

	// 4. list 验证（更新后仍在列表）
	if got := findRoleIDByName(t, roleName); got != roleID {
		t.Fatalf("role %q missing after update, got=%q want=%q", roleName, got, roleID)
	}

	// 5. delete（同时作为清理）
	apiClient.Delete(t, "/api/v1/acl/role/"+roleID)

	// 6. 删除后不再出现
	if got := findRoleIDByName(t, roleName); got != "" {
		t.Fatalf("role %q should be deleted, but still found", roleName)
	}
}

// findRoleIDByName 从角色列表里按名称查找 roleId（字符串）。
func findRoleIDByName(t *testing.T, name string) string {
	t.Helper()
	resp := apiClient.Get(t, "/api/v1/acl/role?page=1&limit=50")
	var list types.ResponseRoleList
	resp.decodeData(&list)
	for _, r := range list.Records {
		if r.RoleName == name {
			return r.RoleID
		}
	}
	return ""
}
