//go:build integration

package tests

import (
	"net/http"
	"testing"

	"github.com/liwook/go-vue-selection/pkg/result"
	"github.com/liwook/go-vue-selection/types"
)

// TestTrademarkFlow 覆盖品牌写链路：create → update → list 验证 → delete。
func TestTrademarkFlow(t *testing.T) {
	const tmName = "IT品牌"

	// 1. create
	apiClient.Post(t, "/api/v1/product/trademark", types.ParamTmSave{
		TmName:  tmName,
		LogoUrl: "https://example.com/logo.png",
	})

	// 2. 从列表反查 tmId
	tmID := findTmIDByName(t, tmName)
	if tmID == "" {
		t.Fatalf("created trademark %q not found", tmName)
	}

	// 3. update
	apiClient.Put(t, "/api/v1/product/trademark/"+tmID, types.ParamTmUpdate{
		TmID:    tmID,
		TmName:  tmName,
		LogoUrl: "https://example.com/logo2.png",
	})

	// 4. list 验证
	if got := findTmIDByName(t, tmName); got != tmID {
		t.Fatalf("trademark %q missing after update, got=%q want=%q", tmName, got, tmID)
	}

	// 5. delete（同时作为清理）
	apiClient.Delete(t, "/api/v1/product/trademark/"+tmID)

	// 6. 删除后不再出现
	if got := findTmIDByName(t, tmName); got != "" {
		t.Fatalf("trademark %q should be deleted, but still found", tmName)
	}
}

func findTmIDByName(t *testing.T, name string) string {
	t.Helper()
	resp := apiClient.Get(t, "/api/v1/product/trademark?page=1&limit=50")
	var list types.ResponseTmList
	resp.decodeData(&list)
	for _, tm := range list.Records {
		if tm.TmName == name {
			return tm.TmID
		}
	}
	return ""
}

// TestTrademarkNegative 验证品牌名重复创建返回品牌业务错误。
func TestTrademarkNegative(t *testing.T) {
	const tmName = "IT品牌-重复"
	// 首次创建成功
	apiClient.Post(t, "/api/v1/product/trademark", types.ParamTmSave{
		TmName:  tmName,
		LogoUrl: "https://example.com/logo.png",
	})
	defer func() {
		if id := findTmIDByName(t, tmName); id != "" {
			apiClient.Delete(t, "/api/v1/product/trademark/"+id)
		}
	}()
	// 二次创建同名应失败
	apiClient.Call(t, "POST", "/api/v1/product/trademark", types.ParamTmSave{
		TmName:  tmName,
		LogoUrl: "https://example.com/logo.png",
	}, http.StatusOK, int(result.CodeTrademarkErr))
}
