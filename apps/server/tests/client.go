//go:build integration

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/liwook/go-vue-selection/config"
	"github.com/liwook/go-vue-selection/pkg/result"
)

// apiResponse 对应后端统一的 {code, message, data} 响应结构。
// StatusCode 为 HTTP 层状态码（200/500…），用于与业务 code 一起做断言。
type apiResponse struct {
	StatusCode int             `json:"-"`
	Code       int             `json:"code"`
	Message    string          `json:"message"`
	Data       json.RawMessage `json:"data"`
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
	if err := c.do(http.MethodPost, "/api/v1/acl/index/login", "", bytes.NewReader(body), &resp); err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login http status=%d", resp.StatusCode)
	}
	if resp.Code != int(result.CodeSuccess) {
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
// 同时记录 HTTP 状态码到 out.StatusCode。
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
	ar.StatusCode = resp.StatusCode
	if out != nil {
		*out = ar
	}
	return nil
}

// request 发送带 Token 的请求，并同时断言 HTTP 状态码与业务 code。
//   - expectHTTP: 期望的 HTTP 状态码，传 0 表示默认 200。
//   - expectCode: 期望的业务 code，传 0 表示默认 result.CodeSuccess(200)；
//     传 -1 表示“不校验业务 code”（例如只想验证 HTTP 500）。
func (c *APIClient) request(t *testing.T, method, path string, body io.Reader, expectHTTP, expectCode int) apiResponse {
	t.Helper()
	if expectHTTP == 0 {
		expectHTTP = http.StatusOK
	}
	var resp apiResponse
	if err := c.do(method, path, c.token, body, &resp); err != nil {
		t.Fatalf("%s %s failed: %v", method, path, err)
	}
	if resp.StatusCode != expectHTTP {
		t.Fatalf("%s %s: expect HTTP %d, got %d (body code=%d message=%s)",
			method, path, expectHTTP, resp.StatusCode, resp.Code, resp.Message)
	}
	if expectCode == -1 {
		return resp
	}
	if expectCode == 0 {
		expectCode = int(result.CodeSuccess)
	}
	if resp.Code != expectCode {
		t.Fatalf("%s %s: expect code %d, got %d (message=%s)",
			method, path, expectCode, resp.Code, resp.Message)
	}
	return resp
}

// Get 发送带 Token 的 GET 请求（HTTP 200 + code 200）。
func (c *APIClient) Get(t *testing.T, path string) apiResponse {
	t.Helper()
	return c.request(t, http.MethodGet, path, nil, 0, 0)
}

// Post 发送带 Token 的 POST 请求（HTTP 200 + code 200）。
func (c *APIClient) Post(t *testing.T, path string, body interface{}) apiResponse {
	t.Helper()
	return c.request(t, http.MethodPost, path, marshalBody(t, body), 0, 0)
}

// Put 发送带 Token 的 PUT 请求（HTTP 200 + code 200）。
func (c *APIClient) Put(t *testing.T, path string, body interface{}) apiResponse {
	t.Helper()
	return c.request(t, http.MethodPut, path, marshalBody(t, body), 0, 0)
}

// Delete 发送带 Token 的 DELETE 请求（HTTP 200 + code 200）。
func (c *APIClient) Delete(t *testing.T, path string) apiResponse {
	t.Helper()
	return c.request(t, http.MethodDelete, path, nil, 0, 0)
}

// Call 低层封装：可自定义期望的 HTTP 状态码与业务 code。
//   - expectCode 传 -1 表示不校验业务 code。
func (c *APIClient) Call(t *testing.T, method, path string, body interface{}, expectHTTP, expectCode int) apiResponse {
	t.Helper()
	var reader io.Reader
	if body != nil {
		reader = marshalBody(t, body)
	}
	return c.request(t, method, path, reader, expectHTTP, expectCode)
}

// Upload 发送 multipart/form-data 文件上传请求（字段名固定为 file）。
func (c *APIClient) Upload(t *testing.T, path, field, filename string) apiResponse {
	t.Helper()
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("open upload file %s: %v", filename, err)
	}
	defer f.Close()

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	part, err := mw.CreateFormFile(field, filepath.Base(filename))
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := io.Copy(part, f); err != nil {
		t.Fatalf("copy file to form: %v", err)
	}
	mw.Close()

	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, &buf)
	if err != nil {
		t.Fatalf("new upload request: %v", err)
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		t.Fatalf("upload request: %v", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read upload body: %v", err)
	}
	var ar apiResponse
	if err := json.Unmarshal(raw, &ar); err != nil {
		t.Fatalf("unmarshal upload response (status=%d body=%s): %v", resp.StatusCode, string(raw), err)
	}
	ar.StatusCode = resp.StatusCode
	if ar.StatusCode != http.StatusOK {
		t.Fatalf("upload %s: expect HTTP 200, got %d (code=%d message=%s)",
			path, ar.StatusCode, ar.Code, ar.Message)
	}
	if ar.Code != int(result.CodeSuccess) {
		t.Fatalf("upload %s: expect code 200, got %d (message=%s)", path, ar.Code, ar.Message)
	}
	return ar
}

// decodeData 将响应中的 data 反序列化为 out（指针）。
func (r apiResponse) decodeData(out interface{}) {
	if len(r.Data) == 0 || string(r.Data) == "null" {
		return
	}
	if err := json.Unmarshal(r.Data, out); err != nil {
		// 不是 fatal：解码失败由调用方判断是否需要
		_ = err
	}
}

// marshalBody 将 body 序列化为 JSON 读取器，出错时直接让测试失败。
func marshalBody(t *testing.T, body interface{}) io.Reader {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}
	return bytes.NewReader(b)
}

// apiClient 为包级单例，由 setupAPIClient 在 TestMain 中初始化。
var apiClient *APIClient
