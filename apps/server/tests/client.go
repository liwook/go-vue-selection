//go:build integration

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/liwook/go-vue-selection/config"
)

// apiResponse 对应后端统一的 {code, message, data} 响应结构。
type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// APIClient 封装集成测试的 HTTP 调用：自动携带登录 Token，统一解析响应。
type APIClient struct {
	baseURL string
	token   string
	http    *http.Client
}

// setupAPIClient 使用 httptest 拉起的 router 初始化客户端，并以 admin/111111 登录获取 Token。
func setupAPIClient(r http.Handler, pg *config.PostgresConfig) error {
	ts := httptest.NewServer(r)
	apiClient = &APIClient{
		baseURL: ts.URL,
		http:    ts.Client(),
	}
	if err := apiClient.login("admin", "111111"); err != nil {
		return fmt.Errorf("login admin failed: %w", err)
	}
	return nil
}

// login 调用登录接口换取 Token 并保存。
func (c *APIClient) login(username, password string) error {
	body, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	var resp apiResponse
	if err := c.do(http.MethodPost, "/admin/acl/index/login", "", bytes.NewReader(body), &resp); err != nil {
		return err
	}
	if resp.Code != 200 {
		return fmt.Errorf("login failed, code=%d message=%s", resp.Code, resp.Message)
	}
	// 登录接口 data 直接为 token 字符串
	var token string
	if err := json.Unmarshal(resp.Data, &token); err != nil {
		return fmt.Errorf("unmarshal login token: %w", err)
	}
	if token == "" {
		return fmt.Errorf("login returned empty token")
	}
	c.token = token
	return nil
}

// do 发送请求并解析统一的 {code,message,data} 结构到 out（可为 nil）。
// 当 expectCode 为 0 时不校验业务 code；否则校验 resp.Code == expectCode。
func (c *APIClient) do(method, path, token string, body io.Reader, out *apiResponse) error {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	var ar apiResponse
	if err := json.Unmarshal(raw, &ar); err != nil {
		return fmt.Errorf("unmarshal response (status=%d body=%s): %w", resp.StatusCode, string(raw), err)
	}
	if out != nil {
		*out = ar
	}
	return nil
}

// Get 发送带 Token 的 GET 请求。
func (c *APIClient) Get(t *testing.T, path string, expectCode int) apiResponse {
	t.Helper()
	return c.authed(t, http.MethodGet, path, nil, expectCode)
}

// Post 发送带 Token 的 POST 请求。
func (c *APIClient) Post(t *testing.T, path string, body interface{}, expectCode int) apiResponse {
	t.Helper()
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(b)
	}
	return c.authed(t, http.MethodPost, path, reader, expectCode)
}

func (c *APIClient) authed(t *testing.T, method, path string, body io.Reader, expectCode int) apiResponse {
	t.Helper()
	var resp apiResponse
	if err := c.do(method, path, c.token, body, &resp); err != nil {
		t.Fatalf("%s %s failed: %v", method, path, err)
	}
	if expectCode != 0 && resp.Code != expectCode {
		t.Fatalf("%s %s: expect code %d, got %d (message=%s)", method, path, expectCode, resp.Code, resp.Message)
	}
	return resp
}

// apiClient 为包级单例，由 setupAPIClient 在 TestMain 中初始化。
var apiClient *APIClient
