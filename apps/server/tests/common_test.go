//go:build integration

package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/liwook/go-vue-selection/pkg/result"
)

// rawCall 发送一个不经过 apiClient token 的原始请求（可自定义 Authorization 头），
// 返回解析后的统一响应结构与底层 HTTP 状态码。
func rawCall(t *testing.T, method, path, authHeader string) apiResponse {
	t.Helper()
	req, err := http.NewRequest(method, apiClient.baseURL+path, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	resp, err := apiClient.http.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var ar apiResponse
	_ = json.Unmarshal(raw, &ar)
	ar.StatusCode = resp.StatusCode
	return ar
}

// TestCommonNoRoute 验证访问不存在的路由返回 CodeNoRoute（兜底）。
func TestCommonNoRoute(t *testing.T) {
	ar := rawCall(t, http.MethodGet, "/api/v1/acl/nonexist", "Bearer "+apiClient.token)
	if ar.StatusCode != http.StatusOK {
		t.Fatalf("expect HTTP 200, got %d", ar.StatusCode)
	}
	if ar.Code != int(result.CodeNoRoute) {
		t.Fatalf("expect CodeNoRoute(%d), got %d (message=%s)", result.CodeNoRoute, ar.Code, ar.Message)
	}
}

// TestCommonAuthMissing 验证缺失 Authorization 头访问受保护接口返回 CodeNeedLogin。
func TestCommonAuthMissing(t *testing.T) {
	ar := rawCall(t, http.MethodGet, "/api/v1/acl/user/info", "")
	if ar.StatusCode != http.StatusOK {
		t.Fatalf("expect HTTP 200, got %d", ar.StatusCode)
	}
	if ar.Code != int(result.CodeNeedLogin) {
		t.Fatalf("expect CodeNeedLogin(%d), got %d (message=%s)", result.CodeNeedLogin, ar.Code, ar.Message)
	}
}

// TestCommonAuthInvalid 验证伪造/失效的 token 返回 CodeInvalidToken。
func TestCommonAuthInvalid(t *testing.T) {
	ar := rawCall(t, http.MethodGet, "/api/v1/acl/user/info", "Bearer not-a-valid-token")
	if ar.StatusCode != http.StatusOK {
		t.Fatalf("expect HTTP 200, got %d", ar.StatusCode)
	}
	if ar.Code != int(result.CodeInvalidToken) {
		t.Fatalf("expect CodeInvalidToken(%d), got %d (message=%s)", result.CodeInvalidToken, ar.Code, ar.Message)
	}
}

// TestCommonCORS 验证跨域预检请求返回 204 并携带 Access-Control-Allow-Origin。
func TestCommonCORS(t *testing.T) {
	req, err := http.NewRequest(http.MethodOptions, apiClient.baseURL+"/api/v1/product/trademark", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	resp, err := apiClient.http.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expect HTTP 204 for CORS preflight, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Access-Control-Allow-Origin"); got == "" {
		t.Fatalf("expect Access-Control-Allow-Origin header in CORS preflight")
	}
}
