//go:build integration

package tests

import (
	"net/http"
	"testing"

	"github.com/liwook/go-vue-selection/pkg/result"
	"github.com/liwook/go-vue-selection/types"
)

// TestUserFlow 覆盖用户写链路：signup → assignRole → lock/unlock → list 验证 → delete。
func TestUserFlow(t *testing.T) {
	const (
		username = "it_user_tmp"
		name     = "IT 临时用户"
	)

	// 1. signup
	apiClient.Post(t, "/admin/acl/user", types.ParamUserSignUp{
		Username: username,
		Name:     name,
		Password: "111111",
	})

	// 2. 从列表反查新用户 userId
	userID := findUserIDByUsername(t, username)
	if userID == "" {
		t.Fatalf("created user %q not found in list", username)
	}

	// 3. assignRole：分配普通用户角色(role_id=2)
	apiClient.Post(t, "/admin/acl/user/"+userID+"/role", types.ParamDoAssignRole{
		UserID:     userID,
		RoleIDList: []string{"2"},
	})

	// 4. lock（禁用）
	apiClient.Post(t, "/admin/acl/user/lock", types.ParamUserLock{
		UserID: userID,
		Status: false,
	})
	// 5. unlock（恢复）
	apiClient.Post(t, "/admin/acl/user/lock", types.ParamUserLock{
		UserID: userID,
		Status: true,
	})

	// 6. list 验证
	if got := findUserIDByUsername(t, username); got != userID {
		t.Fatalf("user %q missing after operations, got=%q want=%q", username, got, userID)
	}

	// 7. delete（同时作为清理）
	apiClient.Delete(t, "/admin/acl/user/"+userID)

	// 8. 删除后不再出现
	if got := findUserIDByUsername(t, username); got != "" {
		t.Fatalf("user %q should be deleted, but still found", username)
	}
}

// TestUserDeleteSelf 验证管理员不能删除自己（handler 层拦截 operatorID==targetUserID）。
func TestUserDeleteSelf(t *testing.T) {
	adminID := adminUserID(t)
	if adminID == "" {
		t.Fatalf("cannot resolve admin userId")
	}
	// 期望：HTTP 200（业务错误不触发 5xx）+ 业务 code=CodeNoPermission
	apiClient.Call(t, "DELETE", "/admin/acl/user/"+adminID, nil,
		http.StatusOK, int(result.CodeNoPermission))
}

// adminUserID 从用户列表按 username="admin" 反查当前登录管理员 userId。
func adminUserID(t *testing.T) string {
	t.Helper()
	return findUserIDByUsername(t, "admin")
}

// findUserIDByUsername 从用户列表按 username 反查 userId。
func findUserIDByUsername(t *testing.T, username string) string {
	t.Helper()
	resp := apiClient.Get(t, "/admin/acl/user?page=1&limit=50")
	var list types.ResponseUserList
	resp.decodeData(&list)
	for _, u := range list.Records {
		if u.Username == username {
			return u.UserID
		}
	}
	return ""
}
