package seed

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/dal/query"
	"github.com/liwook/go-vue-selection/pkg/snowflake"

	"gorm.io/gorm"
)

// RunDemo 初始化商城演示数据（SPU/SKU/平台属性/销售属性等）。
// 这些数据的主键与外键均为雪花ID，原 SQL 里写死会导致与线上生成的 ID 不一致，
// 因此统一在此用 snowflake.GenID() 生成，并通过变量传递引用关系。
// 幂等：以 spu 表是否为空作为判断依据。
func RunDemo(db *gorm.DB) error {
	q := query.Use(db)
	ctx := context.Background()

	// 幂等判断：spu 表已有数据则跳过
	cnt, err := q.Spu.WithContext(ctx).Count()
	if err != nil {
		return fmt.Errorf("seed demo: count spu failed: %w", err)
	}
	if cnt > 0 {
		slog.Info("seed demo skipped: spu table already has data", "count", cnt)
		return nil
	}

	// 注：基础用户与角色（admin/user → role 1/2，以及 guest → role 3）已由 seed.Run 负责初始化，
	// 此处不再重复插入，避免与 seed.Run 的幂等判据冲突。本函数只负责商城演示数据。

	// 固定字典引用（来自 init-sql 中的固定小整数 ID）
	const (
		category3ID     int64 = 61 // category3: 手机
		tmID            int64 = 3  // trademark: OPPO（PG init-sql 中 tm_id=3，对应 spu 名 "oppo"）
		saleAttrColorID int64 = 1
		saleAttrVerID   int64 = 2
	)

	// ---------- 平台属性 attr / attr_value ----------
	attrPhoneID := snowflake.GenID()   // 手机一级
	attrBatteryID := snowflake.GenID() // 电池容量
	attrRamID := snowflake.GenID()     // 运行内存
	attrRomID := snowflake.GenID()     // 机身内存
	attrCpuID := snowflake.GenID()     // CPU型号

	attrs := []*model.Attr{
		{AttrID: attrPhoneID, AttrName: "手机一级", CategoryID: category3ID},
		{AttrID: attrBatteryID, AttrName: "电池容量", CategoryID: category3ID},
		{AttrID: attrRamID, AttrName: "运行内存", CategoryID: category3ID},
		{AttrID: attrRomID, AttrName: "机身内存", CategoryID: category3ID},
		{AttrID: attrCpuID, AttrName: "CPU型号", CategoryID: category3ID},
	}
	if err := q.Attr.WithContext(ctx).Create(attrs...); err != nil {
		return fmt.Errorf("seed demo: insert attr failed: %w", err)
	}

	avPhone1 := snowflake.GenID()
	avPhone2 := snowflake.GenID()
	avPhone3 := snowflake.GenID()
	avPhone4 := snowflake.GenID()
	avBat1 := snowflake.GenID()
	avBat2 := snowflake.GenID()
	avBat3 := snowflake.GenID()
	avRam1 := snowflake.GenID()
	avRam2 := snowflake.GenID()
	avRam3 := snowflake.GenID()
	avRom1 := snowflake.GenID()
	avRom2 := snowflake.GenID()
	avRom3 := snowflake.GenID()
	avRom4 := snowflake.GenID()
	avRom5 := snowflake.GenID()
	avRom6 := snowflake.GenID()
	avCpu1 := snowflake.GenID()
	avCpu2 := snowflake.GenID()

	attrValues := []*model.AttrValue{
		{AttrValueID: avPhone1, ValueName: "安卓手机", AttrID: attrPhoneID},
		{AttrValueID: avPhone2, ValueName: "苹果手机", AttrID: attrPhoneID},
		{AttrValueID: avPhone3, ValueName: "安卓手机222", AttrID: attrPhoneID},
		{AttrValueID: avPhone4, ValueName: "小米手机", AttrID: attrPhoneID},
		{AttrValueID: avBat1, ValueName: "1200mAh以下", AttrID: attrBatteryID},
		{AttrValueID: avBat2, ValueName: "3000mAh以上", AttrID: attrBatteryID},
		{AttrValueID: avBat3, ValueName: "1200mAh到3000mAh", AttrID: attrBatteryID},
		{AttrValueID: avRam1, ValueName: "128G", AttrID: attrRamID},
		{AttrValueID: avRam2, ValueName: "6G", AttrID: attrRamID},
		{AttrValueID: avRam3, ValueName: "256G", AttrID: attrRamID},
		{AttrValueID: avRom1, ValueName: "32G", AttrID: attrRomID},
		{AttrValueID: avRom2, ValueName: "128G", AttrID: attrRomID},
		{AttrValueID: avRom3, ValueName: "512G", AttrID: attrRomID},
		{AttrValueID: avRom4, ValueName: "64G", AttrID: attrRomID},
		{AttrValueID: avRom5, ValueName: "256G", AttrID: attrRomID},
		{AttrValueID: avRom6, ValueName: "1T", AttrID: attrRomID},
		{AttrValueID: avCpu1, ValueName: "骁龙730G", AttrID: attrCpuID},
		{AttrValueID: avCpu2, ValueName: "麒麟990", AttrID: attrCpuID},
	}
	if err := q.AttrValue.WithContext(ctx).Create(attrValues...); err != nil {
		return fmt.Errorf("seed demo: insert attr_value failed: %w", err)
	}

	// ---------- SPU ----------
	spuID := snowflake.GenID()
	if err := q.Spu.WithContext(ctx).Create(&model.Spu{
		SpuID:       spuID,
		SpuName:     "oppo find x7 pro",
		Description: "oppo find x7 pro 是 OPPO 旗下的一款手机，于 2021 年发布。",
		Category3ID: category3ID,
		TmID:        tmID,
	}); err != nil {
		return fmt.Errorf("seed demo: insert spu failed: %w", err)
	}

	// ---------- spu_image_list ----------
	imgURLs := []string{
		"/api/static/img/product/default/oppo1.jpeg",
		"/api/static/img/product/default/oppo2.jpeg",
		"/api/static/img/product/default/oppo3.jpeg",
		"/api/static/img/product/default/oppo4.jpeg",
	}
	spuImages := make([]*model.SpuImageList, 0, len(imgURLs))
	imageName := "oppo"
	for _, u := range imgURLs {
		imageURL := u
		spuImages = append(spuImages, &model.SpuImageList{
			ImageID:   snowflake.GenID(),
			ImageName: &imageName,
			ImageURL:  &imageURL,
			SpuID:     spuID,
		})
	}
	if err := q.SpuImageList.WithContext(ctx).Create(spuImages...); err != nil {
		return fmt.Errorf("seed demo: insert spu_image_list failed: %w", err)
	}

	// ---------- spu_sale_attr（引用固定 sale_attr_id 1/2/3） ----------
	spuSaleAttrs := []*model.SpuSaleAttr{
		{SpuSaleAttrID: snowflake.GenID(), BaseSaleAttrID: saleAttrColorID, SaleAttrName: new("颜色"), SpuID: spuID},
		{SpuSaleAttrID: snowflake.GenID(), BaseSaleAttrID: saleAttrVerID, SaleAttrName: new("版本"), SpuID: spuID},
	}
	if err := q.SpuSaleAttr.WithContext(ctx).Create(spuSaleAttrs...); err != nil {
		return fmt.Errorf("seed demo: insert spu_sale_attr failed: %w", err)
	}

	// ---------- sale_attr_value（引用 spu_id 与固定 sale_attr_id） ----------
	// 注：与 MySQL 源一致，仅演示颜色(1)/版本(2) 两个销售属性，不使用第 3 项（尺码）。
	savBlue := snowflake.GenID()
	sav12256 := snowflake.GenID()
	savBlack := snowflake.GenID()
	sav8128 := snowflake.GenID()

	saleAttrValues := []*model.SaleAttrValue{
		{SaleAttrValueID: savBlue, SpuID: spuID, SaleAttrID: saleAttrColorID, SaleAttrValueName: "粉色"},
		{SaleAttrValueID: savBlack, SpuID: spuID, SaleAttrID: saleAttrColorID, SaleAttrValueName: "红色"},
		{SaleAttrValueID: sav12256, SpuID: spuID, SaleAttrID: saleAttrVerID, SaleAttrValueName: "K1"},
		{SaleAttrValueID: sav8128, SpuID: spuID, SaleAttrID: saleAttrVerID, SaleAttrValueName: "K2"},
	}
	if err := q.SaleAttrValue.WithContext(ctx).Create(saleAttrValues...); err != nil {
		return fmt.Errorf("seed demo: insert sale_attr_value failed: %w", err)
	}

	// ---------- SKU（1 个，引用 spu_id / category3_id / tm_id；内容对照 init.sql 的 sku 表） ----------
	// price 单位：分（1000 元 = 100000 分）
	// weight 单位：毫克（1000 克 = 1000000 毫克）
	const (
		skuName   = "oppo find x7"
		skuDesc   = "松影墨韵 | 大漠银月 | 海阔天空"
		skuPrice  = int64(450000)
		skuWeight = int64(201000)
	)
	id := snowflake.GenID()
	skuIDs := []int64{id}
	skus := []*model.Sku{{
		SkuID:       id,
		SpuID:       spuID,
		Category3ID: category3ID,
		TmID:        tmID,
		SkuName:     skuName,
		WeightMg:    skuWeight,
		PriceCent:   skuPrice,
		SkuDesc:     skuDesc,
		IsSale:      false,
	}}
	if err := q.Sku.WithContext(ctx).Create(skus...); err != nil {
		return fmt.Errorf("seed demo: insert sku failed: %w", err)
	}

	// ---------- sku_attr_value（引用 attr_id / sku_id / attr_value_id） ----------
	// 每个 SKU 关联 3 条平台属性：运行内存 / 机身内存 / CPU型号
	type combo struct {
		attrID int64
		avID   int64
	}
	// 对照 init.sql 的 sku_attr_value：手机一级→安卓手机、电池容量→1200mAh以下、
	// 运行内存→128G、机身内存→32G、CPU型号→枭龙730G
	skuAttrMap := [][]combo{
		{{attrPhoneID, avPhone1}, {attrBatteryID, avBat1}, {attrRamID, avRam1}, {attrRomID, avRom1}, {attrCpuID, avCpu1}},
	}
	var skuAttrValues []*model.SkuAttrValue
	for i, skuID := range skuIDs {
		for _, c := range skuAttrMap[i] {
			skuAttrValues = append(skuAttrValues, &model.SkuAttrValue{
				SkuAttrValueID: snowflake.GenID(),
				AttrID:         c.attrID,
				SkuID:          skuID,
				AttrValueID:    c.avID,
			})
		}
	}
	if err := q.SkuAttrValue.WithContext(ctx).Create(skuAttrValues...); err != nil {
		return fmt.Errorf("seed demo: insert sku_attr_value failed: %w", err)
	}

	// ---------- sku_image（引用 sku_id 与 spu_image_id） ----------
	// 注：sku_image.is_default 是 boolean 列，而 gen 生成的 model 误判为 *string，
	// 为避免类型不匹配，这里用原生 SQL 精确插入。
	skuImageDef := []struct {
		spuImgIdx int
		isDefault bool
	}{
		{0, true}, {1, false}, {2, false}, {3, false},
	}
	for _, si := range skuImageDef {
		img := spuImages[si.spuImgIdx]
		if err := db.WithContext(ctx).Exec(
			`INSERT INTO app.sku_image (image_id, sku_id, image_url, spu_image_id, is_default)
			 VALUES ($1,$2,$3,$4,$5)`,
			snowflake.GenID(), skuIDs[0], img.ImageURL, img.ImageID, si.isDefault,
		).Error; err != nil {
			return fmt.Errorf("seed demo: insert sku_image failed: %w", err)
		}
	}

	// ---------- sku_sale_attr_value（引用 sku_id 与 sale_attr_value_id） ----------
	skuSaleMap := [][]int64{
		{savBlue, sav12256},
	}
	var skuSaleAttrValues []*model.SkuSaleAttrValue
	for i, skuID := range skuIDs {
		for _, savID := range skuSaleMap[i] {
			skuSaleAttrValues = append(skuSaleAttrValues, &model.SkuSaleAttrValue{
				SkuSaleAttrValueID: snowflake.GenID(),
				SkuID:              skuID,
				SaleAttrValueID:    savID,
			})
		}
	}
	if err := q.SkuSaleAttrValue.WithContext(ctx).Create(skuSaleAttrValues...); err != nil {
		return fmt.Errorf("seed demo: insert sku_sale_attr_value failed: %w", err)
	}

	slog.Info("seed demo done: inserted spu/sku/attr demo data")
	return nil
}
