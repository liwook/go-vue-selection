//go:build integration

package tests

import (
	"testing"

	"github.com/liwook/go-vue-selection/types"
)

// TestAttrFlow 覆盖基础属性写链路：建 c2/c3 → save → get 验证 → delete。
// 为避免跨文件依赖，分类在此自包含创建（复用包内 findC2ByName/findC3ByName）。
func TestAttrFlow(t *testing.T) {
	const (
		c2Name = "ITAttr二级"
		c3Name = "ITAttr三级"
		c1ID   = "1" // 图书
		attrName = "IT属性"
	)

	// 1. 准备分类（category2 / category3）
	apiClient.Post(t, "/admin/product/category2", types.ParamC2Create{
		Name:        c2Name,
		Category1ID: c1ID,
	})
	c2ID := findC2ByName(t, c1ID, c2Name)
	if c2ID == "" {
		t.Fatalf("created category2 %q not found", c2Name)
	}
	apiClient.Post(t, "/admin/product/category3", types.ParamC3Create{
		Name:        c3Name,
		Category2ID: c2ID,
	})
	c3ID := findC3ByName(t, c2ID, c3Name)
	if c3ID == "" {
		t.Fatalf("created category3 %q not found", c3Name)
	}

	// 2. save attr
	apiClient.Post(t, "/admin/product/attr", types.Attr{
		AttrName:   attrName,
		CategoryID: c3ID,
		AttrValueList: []*types.AttrValue{
			{ValueName: "红色"},
		},
	})

	// 3. get 验证并取 attrId
	attrID := findAttrIDByName(t, c1ID, c2ID, c3ID, attrName)
	if attrID == "" {
		t.Fatalf("created attr %q not found", attrName)
	}

	// 4. delete（同时作为清理）
	apiClient.Delete(t, "/admin/product/attr/"+attrID)

	// 5. 删除后不再出现
	if got := findAttrIDByName(t, c1ID, c2ID, c3ID, attrName); got != "" {
		t.Fatalf("attr %q should be deleted, but still found", attrName)
	}
}

func findAttrIDByName(t *testing.T, c1, c2, c3, name string) string {
	t.Helper()
	resp := apiClient.Get(t, "/admin/product/attr/"+c1+"/"+c2+"/"+c3)
	var list []*types.Attr
	resp.decodeData(&list)
	for _, a := range list {
		if a.AttrName == name {
			return a.AttrID
		}
	}
	return ""
}
