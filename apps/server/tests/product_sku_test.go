//go:build integration

package tests

import (
	"testing"

	"github.com/liwook/go-vue-selection/types"
)

// TestSkuFlow 覆盖 SKU 写链路：save → onsale → cancelsale → info 验证 → delete。
func TestSkuFlow(t *testing.T) {
	const spuName = "IT-SKU-SPU"

	spuID, c3ID, tmID := prepareSpuForTest(t, spuName)

	// 1. save
	apiClient.Post(t, "/admin/product/sku", types.SkuInfo{
		SpuID:       spuID,
		Category3ID: c3ID,
		TmID:        tmID,
		SkuName:     "IT-SKU-001",
		WeightMg:    500,
		PriceCent:   9900,
		SkuDesc:     "integration test sku",
		SkuImageList: []*types.SkuImgDTO{
			{ImageName: "cover", ImageURL: "https://example.com/cover.png"},
		},
	})

	// 2. 从 FindBySpuId 反查 skuId
	skuID := findSkuIDBySpu(t, spuID)
	if skuID == "" {
		t.Fatalf("created sku for spu %q not found", spuID)
	}

	// 3. onsale（上架）
	apiClient.Put(t, "/admin/product/sku/"+skuID+"/onsale", nil)
	if got := skuIsSale(t, skuID); got != 1 {
		t.Fatalf("after onsale expect isSale=1, got=%d", got)
	}

	// 4. cancelsale（下架）
	apiClient.Put(t, "/admin/product/sku/"+skuID+"/cancelsale", nil)
	if got := skuIsSale(t, skuID); got != 0 {
		t.Fatalf("after cancelsale expect isSale=0, got=%d", got)
	}

	// 5. info 验证
	resp := apiClient.Get(t, "/admin/product/sku/"+skuID)
	var info types.ResponseSkuInfo
	resp.decodeData(&info)
	if info.SkuID != skuID {
		t.Fatalf("sku info id mismatch, got=%q want=%q", info.SkuID, skuID)
	}

	// 6. delete（同时作为清理）
	apiClient.Delete(t, "/admin/product/sku/"+skuID)
	if got := findSkuIDBySpu(t, spuID); got != "" {
		t.Fatalf("sku %q should be deleted, but still found", skuID)
	}
}

func findSkuIDBySpu(t *testing.T, spuID string) string {
	t.Helper()
	resp := apiClient.Get(t, "/admin/product/spu/"+spuID+"/sku")
	var list []*types.ResponseSkuInfo
	resp.decodeData(&list)
	if len(list) == 0 {
		return ""
	}
	return list[0].SkuID
}

func skuIsSale(t *testing.T, skuID string) int8 {
	t.Helper()
	resp := apiClient.Get(t, "/admin/product/sku/"+skuID)
	var info types.ResponseSkuInfo
	resp.decodeData(&info)
	return info.IsSale
}

// TestSkuSaleAttr 验证 SKU 销售属性子表可正确写入：从 SPU 销售属性值回填 SkuSaleAttrValueList。
func TestSkuSaleAttr(t *testing.T) {
	const spuName = "IT-SKU-SA"

	spuID, c3ID, tmID := prepareSpuForTest(t, spuName)

	// 1. 取 SPU 下的销售属性（含 saleAttrValueId 外键）
	resp := apiClient.Get(t, "/admin/product/spu/"+spuID+"/saleAttr")
	var saleAttrs []*types.SpuSaleAttr
	resp.decodeData(&saleAttrs)
	if len(saleAttrs) == 0 {
		t.Fatalf("spu %q should have sale attrs", spuID)
	}
	var saleAttrValues []*types.SkuSaleAttrValueDTO
	for _, sa := range saleAttrs {
		for _, v := range sa.SpuSaleAttrValue {
			saleAttrValues = append(saleAttrValues, &types.SkuSaleAttrValueDTO{
				SaleAttrID:        sa.BaseSaleAttrId,
				SaleAttrValueID:   v.SaleAttrValueID,
				SaleAttrName:      sa.SaleAttrName,
				SaleAttrValueName: v.SaleAttrValueName,
			})
		}
	}
	if len(saleAttrValues) == 0 {
		t.Fatalf("spu %q should have at least one sale attr value", spuID)
	}

	// 2. save SKU（携带销售属性子表）
	apiClient.Post(t, "/admin/product/sku", types.SkuInfo{
		SpuID:                spuID,
		Category3ID:          c3ID,
		TmID:                 tmID,
		SkuName:              "IT-SKU-SA-001",
		WeightMg:             500,
		PriceCent:            9900,
		SkuDesc:              "integration test sku with sale attr",
		SkuSaleAttrValueList: saleAttrValues,
	})

	// 3. 反查 skuId
	skuID := findSkuIDBySpu(t, spuID)
	if skuID == "" {
		t.Fatalf("created sku for spu %q not found", spuID)
	}
	t.Cleanup(func() { apiClient.Delete(t, "/admin/product/sku/"+skuID) })

	// 4. info 验证销售属性子表已写入
	infoResp := apiClient.Get(t, "/admin/product/sku/"+skuID)
	var info types.ResponseSkuInfo
	infoResp.decodeData(&info)
	if len(info.SkuSaleAttrValueList) == 0 {
		t.Fatalf("sku %q should have sale attr value list", skuID)
	}
	got := info.SkuSaleAttrValueList[0]
	if got.SaleAttrValueID != saleAttrValues[0].SaleAttrValueID {
		t.Fatalf("sku sale attr value id mismatch, got=%q want=%q",
			got.SaleAttrValueID, saleAttrValues[0].SaleAttrValueID)
	}
}
