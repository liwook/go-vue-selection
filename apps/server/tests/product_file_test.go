//go:build integration

package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFileUpload 覆盖文件上传链路：构造临时文件 → multipart 上传 → 断言返回 URL。
func TestFileUpload(t *testing.T) {
	// 1. 造一个临时图片文件用于上传
	dir := t.TempDir()
	src := filepath.Join(dir, "it-upload.png")
	if err := os.WriteFile(src, []byte("fake-image-bytes"), 0o644); err != nil {
		t.Fatalf("write temp upload file: %v", err)
	}

	// 2. 上传（multipart，字段名 file）
	resp := apiClient.Upload(t, "/api/v1/product/file/upload", "file", src)

	// 3. 断言返回 data 为文件 URL（实际返回形如 /api/img/日期/文件名.png）
	var url string
	resp.decodeData(&url)
	if !strings.Contains(url, "img") || url == "" {
		t.Fatalf("upload returned unexpected url: %q", url)
	}
}
