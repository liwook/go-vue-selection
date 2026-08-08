//go:build integration

package tests

import (
	"encoding/json"
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

// TestUserLoginNegative 验证密码错误时登录被拒绝（统一返回 CodeInvalidPassword，不泄露账号状态）。
func TestUserLoginNegative(t *testing.T) {
	const (
		username = "it_login_neg"
		password = "111111"
	)
	// 预置一个已知用户
	userID := signupUser(t, username, password)
	t.Cleanup(func() { apiClient.Delete(t, "/admin/acl/user/"+userID) })

	// 错误密码：HTTP 200 + 业务 CodeInvalidPassword
	apiClient.Call(t, "POST", "/admin/acl/index/login", types.ParamUserLogin{
		Username: username,
		Password: "wrong-password",
	}, http.StatusOK, int(result.CodeInvalidPassword))

	// 正确密码：应成功并返回 token
	resp := apiClient.Call(t, "POST", "/admin/acl/index/login", types.ParamUserLogin{
		Username: username,
		Password: password,
	}, http.StatusOK, int(result.CodeSuccess))
	if token := decodeToken(t, resp); token == "" {
		t.Fatalf("login token should not be empty")
	}
}

// TestUserLogout 验证已登录用户可正常登出。
func TestUserLogout(t *testing.T) {
	const (
		username = "it_logout"
		password = "111111"
	)
	userID := signupUser(t, username, password)
	t.Cleanup(func() { apiClient.Delete(t, "/admin/acl/user/"+userID) })

	// 携带 admin token 登出（登出接口仅校验 token 有效性）
	apiClient.Post(t, "/admin/acl/index/logout", nil)
}

// TestUserToAssign 验证查询某用户可分配角色列表接口（GET /user/:userId/role）。
func TestUserToAssign(t *testing.T) {
	adminID := adminUserID(t)
	if adminID == "" {
		t.Fatalf("cannot resolve admin userId")
	}
	resp := apiClient.Get(t, "/admin/acl/user/"+adminID+"/role")
	var toAssign types.ResponseToAssignRole
	resp.decodeData(&toAssign)
	if len(toAssign.AllRolesList) == 0 {
		t.Fatalf("ToAssign should return at least one assignable role, got empty")
	}
}

// TestUserLockThenLoginDenied 验证锁定用户无法登录，解锁后可恢复。
func TestUserLockThenLoginDenied(t *testing.T) {
	const (
		username = "it_locked"
		password = "111111"
	)
	userID := signupUser(t, username, password)
	t.Cleanup(func() { apiClient.Delete(t, "/admin/acl/user/"+userID) })

	// 锁定
	apiClient.Post(t, "/admin/acl/user/lock", types.ParamUserLock{
		UserID: userID,
		Status: false,
	})
	// 锁定后登录应被拒绝（业务码与密码错误一致，避免泄露账号状态）
	apiClient.Call(t, "POST", "/admin/acl/index/login", types.ParamUserLogin{
		Username: username,
		Password: password,
	}, http.StatusOK, int(result.CodeInvalidPassword))

	// 解锁
	apiClient.Post(t, "/admin/acl/user/lock", types.ParamUserLock{
		UserID: userID,
		Status: true,
	})
	// 解锁后登录应成功
	resp := apiClient.Call(t, "POST", "/admin/acl/index/login", types.ParamUserLogin{
		Username: username,
		Password: password,
	}, http.StatusOK, int(result.CodeSuccess))
	if token := decodeToken(t, resp); token == "" {
		t.Fatalf("login token should not be empty after unlock")
	}
}

// signupUser 注册一个临时用户并返回其 userId。
func signupUser(t *testing.T, username, password string) string {
	t.Helper()
	apiClient.Post(t, "/admin/acl/user", types.ParamUserSignUp{
		Username: username,
		Name:     username,
		Password: password,
	})
	userID := findUserIDByUsername(t, username)
	if userID == "" {
		t.Fatalf("signup user %q not found in list", username)
	}
	return userID
}

// decodeToken 从登录响应的 data（token 字符串）中解析 token。
func decodeToken(t *testing.T, resp apiResponse) string {
	t.Helper()
	var token string
	if err := json.Unmarshal(resp.Data, &token); err != nil {
		t.Fatalf("decode login token: %v", err)
	}
	return token
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
