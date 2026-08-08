//go:build integration

package tests

import (
	"testing"

	"github.com/liwook/go-vue-selection/types"
)

// TestSpuFlow 覆盖 SPU 写链路：准备 tm/c3/baseSaleAttr → save → update → list 验证 → delete。
func TestSpuFlow(t *testing.T) {
	const spuName = "IT-SPU-001"

	spuID, c3ID, tmID := prepareSpuForTest(t, spuName)
	_ = tmID

	// 6. update
	apiClient.Put(t, "/admin/product/spu/"+spuID, types.Spu{
		SpuID:       spuID,
		SpuName:     spuName,
		Description: "integration test spu (updated)",
		Category3ID: c3ID,
		TmID:        tmID,
	})

	// 7. list 验证
	if got := findSpuIDByName(t, c3ID, spuName); got != spuID {
		t.Fatalf("spu %q missing after update, got=%q want=%q", spuName, got, spuID)
	}

	// 8. delete（同时作为清理）
	apiClient.Delete(t, "/admin/product/spu/"+spuID)
	if got := findSpuIDByName(t, c3ID, spuName); got != "" {
		t.Fatalf("spu %q should be deleted, but still found", spuName)
	}
}

// prepareSpuForTest 创建独立的 trademark/category3/baseSaleAttr 并保存一个 SPU，
// 返回 spuId、category3Id、tmId，供 SPU/SKU 测试共用。
func prepareSpuForTest(t *testing.T, spuName string) (spuID, c3ID, tmID string) {
	t.Helper()
	const (
		c2Name = "ITSpu二级"
		c3Name = "ITSpu三级"
		c1ID   = "1"
	)
	tmName := "IT品牌-" + spuName // 不同 spuName 生成不同商标名，避免同库冲突

	// 1. 准备 trademark
	apiClient.Post(t, "/admin/product/trademark", types.ParamTmSave{
		TmName:  tmName,
		LogoUrl: "https://example.com/logo.png",
	})
	tmID = findTmIDByName(t, tmName)
	if tmID == "" {
		t.Fatalf("created trademark %q not found", tmName)
	}

	// 2. 准备分类 c3
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
	c3ID = findC3ByName(t, c2ID, c3Name)
	if c3ID == "" {
		t.Fatalf("created category3 %q not found", c3Name)
	}

	// 3. 取基础销售属性 id
	saleAttrID := firstBaseSaleAttrID(t)
	if saleAttrID == "" {
		t.Fatalf("no base sale attr available")
	}

	// 4. save
	apiClient.Post(t, "/admin/product/spu", types.Spu{
		SpuName:     spuName,
		Description: "integration test spu",
		Category3ID: c3ID,
		TmID:        tmID,
		SpuImageList: []*types.SpuImage{
			{ImageName: "cover", ImageUrl: "https://example.com/cover.png"},
		},
		SpuSaleAttrList: []*types.SpuSaleAttr{
			{
				BaseSaleAttrId: saleAttrID,
				SaleAttrName:   "颜色",
				SpuSaleAttrValue: []*types.SaleAttrValue{
					{BaseSaleAttrId: saleAttrID, SaleAttrValueName: "红色"},
				},
			},
		},
	})

	// 5. 从 list 反查 spuId
	spuID = findSpuIDByName(t, c3ID, spuName)
	if spuID == "" {
		t.Fatalf("created spu %q not found", spuName)
	}
	return spuID, c3ID, tmID
}

func firstBaseSaleAttrID(t *testing.T) string {
	t.Helper()
	resp := apiClient.Get(t, "/admin/product/baseSaleAttr")
	var list []*types.SaleAttr
	resp.decodeData(&list)
	if len(list) == 0 {
		return ""
	}
	return list[0].SaleAttrID
}

func findSpuIDByName(t *testing.T, c3ID, name string) string {
	t.Helper()
	resp := apiClient.Get(t, "/admin/product/spu?page=1&limit=50&category3Id="+c3ID)
	var list types.ResponseSpuList
	resp.decodeData(&list)
	for _, s := range list.Records {
		if s.SpuName == name {
			return s.SpuID
		}
	}
	return ""
}
