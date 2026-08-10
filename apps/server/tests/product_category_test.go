//go:build integration

package tests

import (
	"net/http"
	"testing"

	"github.com/liwook/go-vue-selection/pkg/result"
	"github.com/liwook/go-vue-selection/types"
)

// TestCategoryFlow 覆盖分类写链路：createC2 → createC3 → getC2 验证 → getC3 验证。
// 注：商品分类无 DELETE 接口，新增数据保留（符合文档说明）。
func TestCategoryFlow(t *testing.T) {
	const (
		c2Name = "IT二级分类"
		c3Name = "IT三级分类"
		c1ID   = "1" // init.sql 固定：图书
	)

	// 1. createC2
	apiClient.Post(t, "/api/v1/product/category2", types.ParamC2Create{
		Name:        c2Name,
		Category1ID: c1ID,
	})

	// 2. getC2 反查 category2Id
	c2ID := findC2ByName(t, c1ID, c2Name)
	if c2ID == "" {
		t.Fatalf("created category2 %q not found", c2Name)
	}

	// 3. createC3
	apiClient.Post(t, "/api/v1/product/category3", types.ParamC3Create{
		Name:        c3Name,
		Category2ID: c2ID,
	})

	// 4. getC3 验证
	c3ID := findC3ByName(t, c2ID, c3Name)
	if c3ID == "" {
		t.Fatalf("created category3 %q not found", c3Name)
	}

	// 5. 二次验证 getC2 仍包含该二级分类
	if got := findC2ByName(t, c1ID, c2Name); got != c2ID {
		t.Fatalf("category2 %q missing, got=%q want=%q", c2Name, got, c2ID)
	}

	_ = c3ID
}

func findC2ByName(t *testing.T, c1ID, name string) string {
	t.Helper()
	resp := apiClient.Get(t, "/api/v1/product/category2/"+c1ID)
	var list []types.Category2
	resp.decodeData(&list)
	for _, c := range list {
		if c.Name == name {
			return c.Category2ID
		}
	}
	return ""
}

func findC3ByName(t *testing.T, c2ID, name string) string {
	t.Helper()
	resp := apiClient.Get(t, "/api/v1/product/category3/"+c2ID)
	var list []types.Category3
	resp.decodeData(&list)
	for _, c := range list {
		if c.Name == name {
			return c.Category3ID
		}
	}
	return ""
}

// TestCategoryNegative 验证三级分类绑定不存在的父级（外键约束失败）返回分类 DB 错误。
func TestCategoryNegative(t *testing.T) {
	// 使用一个不可能存在的二级分类 ID，插入会因外键约束失败
	apiClient.Call(t, "POST", "/api/v1/product/category3", types.ParamC3Create{
		Name:        "IT三级分类-孤儿",
		Category2ID: "0",
	}, http.StatusOK, int(result.CodeCategoryDBErr))
}
