//go:build integration

package tests

import (
	"io"
	"net/http"
	"testing"
)

// TestHealth 冒烟测试：验证测试地基（建库 + 建表 + seed + httptest 拉起 router）可用，
// 直接请求 /health 不应走统一响应体，而是返回纯文本 "I'm OK!"。
func TestHealth(t *testing.T) {
	if apiClient == nil {
		t.Fatal("apiClient not initialized, TestMain may have failed")
	}

	resp, err := apiClient.http.Get(apiClient.baseURL + "/health")
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /health: expect status 200, got %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read /health body: %v", err)
	}
	const want = "I'm OK!"
	if string(body) != want {
		t.Fatalf("GET /health: expect body %q, got %q", want, string(body))
	}
}
