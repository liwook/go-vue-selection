//go:build integration

package tests

import (
	"testing"

	"github.com/liwook/go-vue-selection/types"
)

// TestTrademarkFlow 覆盖品牌写链路：create → update → list 验证 → delete。
func TestTrademarkFlow(t *testing.T) {
	const tmName = "IT品牌"

	// 1. create
	apiClient.Post(t, "/admin/product/trademark", types.ParamTmSave{
		TmName:  tmName,
		LogoUrl: "https://example.com/logo.png",
	})

	// 2. 从列表反查 tmId
	tmID := findTmIDByName(t, tmName)
	if tmID == "" {
		t.Fatalf("created trademark %q not found", tmName)
	}

	// 3. update
	apiClient.Put(t, "/admin/product/trademark/"+tmID, types.ParamTmUpdate{
		TmID:    tmID,
		TmName:  tmName,
		LogoUrl: "https://example.com/logo2.png",
	})

	// 4. list 验证
	if got := findTmIDByName(t, tmName); got != tmID {
		t.Fatalf("trademark %q missing after update, got=%q want=%q", tmName, got, tmID)
	}

	// 5. delete（同时作为清理）
	apiClient.Delete(t, "/admin/product/trademark/"+tmID)

	// 6. 删除后不再出现
	if got := findTmIDByName(t, tmName); got != "" {
		t.Fatalf("trademark %q should be deleted, but still found", tmName)
	}
}

func findTmIDByName(t *testing.T, name string) string {
	t.Helper()
	resp := apiClient.Get(t, "/admin/product/trademark?page=1&limit=50")
	var list types.ResponseTmList
	resp.decodeData(&list)
	for _, tm := range list.Records {
		if tm.TmName == name {
			return tm.TmID
		}
	}
	return ""
}
