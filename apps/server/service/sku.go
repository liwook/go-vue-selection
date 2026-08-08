package service

import (
	"context"
	"log/slog"
	"math"
	"strconv"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/pkg/errs"
	"github.com/liwook/go-vue-selection/pkg/idconv"
	"github.com/liwook/go-vue-selection/pkg/result"
	"github.com/liwook/go-vue-selection/pkg/snowflake"
	"github.com/liwook/go-vue-selection/repository"
	"github.com/liwook/go-vue-selection/types"
)

type skuService struct {
	skuRepo *repository.SkuRepo
	spuRepo *repository.SpuRepo
}

func NewSkuService(skuRepo *repository.SkuRepo, spuRepo *repository.SpuRepo) *skuService {
	return &skuService{skuRepo: skuRepo, spuRepo: spuRepo}
}

// 以下 toResponse* 函数把 dal/model 转换为对外 DTO（types）。
// model→DTO 的转换统一放在 service 层，repo 层只返回纯 model，保持分层干净。

func toResponseSku(s *model.Sku) types.Sku {
	m := types.Sku{
		SkuID:       idconv.ToStr(s.SkuID),
		SpuID:       idconv.ToStr(s.SpuID),
		Category3ID: idconv.ToStr(s.Category3ID),
		TmID:        idconv.ToStr(s.TmID),
		SkuName:     s.SkuName,
		PriceCent:   s.PriceCent, // 分（与数据库列一致，无换算）
		WeightMg:    s.WeightMg,  // 毫克（与数据库列一致，无换算）
		SkuDesc:     s.SkuDesc,
		BaseModel: types.BaseModel{
			CreateTime: s.CreateTime,
			UpdateTime: s.UpdateTime,
		},
	}
	if s.IsSale {
		m.IsSale = 1
	} else {
		m.IsSale = 0
	}
	return m
}

func toResponseSkuImg(i *model.SkuImage) *types.SkuImg {
	var spuImageID string
	if i.SpuImageID != nil {
		spuImageID = idconv.ToStr(*i.SpuImageID)
	}
	return &types.SkuImg{
		ImageID:    idconv.ToStr(i.ImageID),
		SkuID:      idconv.ToStr(i.SkuID),
		ImageURL:   i.ImageURL,
		SpuImageID: spuImageID,
		IsDefault:  strconv.FormatBool(i.IsDefault),
		BaseModel: types.BaseModel{
			CreateTime: i.CreateTime,
			UpdateTime: i.UpdateTime,
		},
	}
}

func toResponseSkuAttrValue(v *model.SkuAttrValue) *types.SkuAttrValue {
	return &types.SkuAttrValue{
		SkuAttrValueID: idconv.ToStr(v.SkuAttrValueID),
		AttrID:         idconv.ToStr(v.AttrID),
		ValueID:        idconv.ToStr(v.AttrValueID),
		SkuID:          idconv.ToStr(v.SkuID),
		BaseModel: types.BaseModel{
			CreateTime: v.CreateTime,
			UpdateTime: v.UpdateTime,
		},
	}
}

func toResponseSkuSaleAttrValue(v *model.SkuSaleAttrValue) *types.SkuSaleAttrValue {
	return &types.SkuSaleAttrValue{
		SkuSaleAttrValueID: idconv.ToStr(v.SkuSaleAttrValueID),
		SaleAttrValueID:    idconv.ToStr(v.SaleAttrValueID),
		SkuID:              idconv.ToStr(v.SkuID),
		BaseModel: types.BaseModel{
			CreateTime: v.CreateTime,
			UpdateTime: v.UpdateTime,
		},
	}
}

// toResponseSkuInfo 把单条 SKU 及其关联子表数据拼装成对外 DTO。
// sku 必传；imgs / attrVals / saleVals 为该 sku 关联的子表数据（来自 repo 的批量查询）。
func toResponseSkuInfo(
	sku *model.Sku,
	imgs []*model.SkuImage,
	attrVals []*model.SkuAttrValue,
	saleVals []*model.SkuSaleAttrValue,
) *types.ResponseSkuInfo {
	imgList := make([]*types.SkuImg, 0, len(imgs))
	for _, i := range imgs {
		imgList = append(imgList, toResponseSkuImg(i))
	}
	attrList := make([]*types.SkuAttrValue, 0, len(attrVals))
	for _, v := range attrVals {
		attrList = append(attrList, toResponseSkuAttrValue(v))
	}
	saleList := make([]*types.SkuSaleAttrValue, 0, len(saleVals))
	for _, v := range saleVals {
		saleList = append(saleList, toResponseSkuSaleAttrValue(v))
	}
	return &types.ResponseSkuInfo{
		Sku:                  toResponseSku(sku),
		SkuImageList:         imgList,
		SkuAttrValueList:     attrList,
		SkuSaleAttrValueList: saleList,
	}
}

func (s *skuService) SaveSkuInfo(ctx context.Context, skuInfo *types.SkuInfo) error {
	// SKU
	skuId := snowflake.GenID()
	sku := &model.Sku{
		SkuID:       skuId, // 原生 int64
		SpuID:       idconv.ToInt64Safe(skuInfo.SpuID),
		Category3ID: idconv.ToInt64Safe(skuInfo.Category3ID),
		TmID:        idconv.ToInt64Safe(skuInfo.TmID),
		SkuName:     skuInfo.SkuName,
		WeightMg:    skuInfo.WeightMg,
		PriceCent:   skuInfo.PriceCent,
		SkuDesc:     skuInfo.SkuDesc,
		IsSale:      skuInfo.IsSale != 0,
	}

	// SKU 图片列表
	skuImageList := make([]*model.SkuImage, 0, len(skuInfo.SkuImageList))
	if len(skuInfo.SkuImageList) > 0 {
		for _, image := range skuInfo.SkuImageList {
			var spuImageID *int64
			if image.SpuImageID != "" {
				v := idconv.ToInt64Safe(image.SpuImageID)
				spuImageID = &v
			}
			skuImageList = append(skuImageList, &model.SkuImage{
				ImageID:    snowflake.GenID(), // 原生 int64
				SkuID:      skuId,
				ImageURL:   image.ImageURL,
				SpuImageID: spuImageID,
				IsDefault:  image.IsDefault == "1",
			})
		}
	}

	// SKU 属性
	skuAttrValueList := make([]*model.SkuAttrValue, 0, len(skuInfo.SkuAttrValueList))
	if len(skuInfo.SkuAttrValueList) > 0 {
		for _, attrValue := range skuInfo.SkuAttrValueList {
			skuAttrValueList = append(skuAttrValueList, &model.SkuAttrValue{
				SkuAttrValueID: snowflake.GenID(), // 原生 int64
				AttrID:         idconv.ToInt64Safe(attrValue.AttrID),
				AttrValueID:    idconv.ToInt64Safe(attrValue.ValueID),
				SkuID:          skuId,
			})
		}
	}

	// SKU 销售属性
	skuSaleAttrValueList := make([]*model.SkuSaleAttrValue, 0, len(skuInfo.SkuSaleAttrValueList))
	if len(skuInfo.SkuSaleAttrValueList) > 0 {
		for _, saleAttrValue := range skuInfo.SkuSaleAttrValueList {
			skuSaleAttrValueList = append(skuSaleAttrValueList, &model.SkuSaleAttrValue{
				SkuSaleAttrValueID: snowflake.GenID(), // 原生 int64
				SaleAttrValueID:    idconv.ToInt64Safe(saleAttrValue.SaleAttrValueID),
				SkuID:              skuId,
			})
		}
	}

	err := s.skuRepo.SaveSkuInfo(ctx, sku, skuImageList, skuAttrValueList, skuSaleAttrValueList)
	if err != nil {
		slog.Error("保存SKU失败", slog.String("spu_id", skuInfo.SpuID), slog.Any("error", err))
		return errs.Wrap(result.CodeSkuDBErr, err)
	}
	return nil
}

func (s *skuService) FindBySpuId(ctx context.Context, spuId int64) ([]*types.ResponseSkuInfo, error) {
	// repo 返回纯 model 列表，这里收集 skuID 后通过三次批量查询一次性拉取关联数据，
	// 再拼装成 DTO。无论该 SPU 下有多少 SKU，关联数据都只需 3 次 SQL，消除了按条查询的 N+1。
	skuList, err := s.skuRepo.FindBySpuId(ctx, spuId)
	if err != nil {
		slog.Error("按SPU查询SKU失败", slog.Int64("spu_id", spuId), slog.Any("error", err))
		return nil, errs.Wrap(result.CodeSkuDBErr, err)
	}
	return s.assembleSkuInfos(ctx, skuList), nil
}

func (s *skuService) GetSkuList(ctx context.Context, page, limit int64) (*types.ResponseSkuInfoList, error) {
	// repo 返回纯 model 列表（主表 sku），这里收集 skuID 后通过三次批量查询一次性拉取关联数据，
	// 再拼装成 DTO。无论列表有多少条 SKU，关联数据都只需 3 次 SQL，
	// 把原本 1 + N×3 的查询降为固定的 1（计数）+ 1（列表）+ 3（批量关联）= 5 次，消除了 N+1。
	skuList, count, err := s.skuRepo.GetSkuList(ctx, page, limit)
	if err != nil {
		slog.Error("查询SKU列表失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeSkuDBErr, err)
	}
	data := &types.ResponseSkuInfoList{
		Records:     s.assembleSkuInfos(ctx, skuList),
		Total:       count,
		Size:        limit,
		Current:     page,
		SearchCount: true,
		Pages:       int64(math.Ceil(float64(count) / float64(limit))),
	}
	return data, nil
}

// assembleSkuInfos 把 SKU model 列表与其关联子表数据（图片/平台属性值/销售属性值）拼装成对外 DTO 列表。
// 关联数据通过 repo 的批量方法一次性查询，避免逐条 SKU 各发 3 次查询产生的 N+1 问题。
func (s *skuService) assembleSkuInfos(ctx context.Context, skuList []*model.Sku) []*types.ResponseSkuInfo {
	if len(skuList) == 0 {
		return []*types.ResponseSkuInfo{}
	}

	// 收集本批 SKU 的 ID
	skuIDs := make([]int64, 0, len(skuList))
	for _, sku := range skuList {
		skuIDs = append(skuIDs, sku.SkuID)
	}

	// 三次批量查询（每次一条 SQL，用 IN 覆盖整批 SKU），替代 N+1 的逐条查询
	imgMap, err := s.skuRepo.BatchGetSkuImages(ctx, skuIDs)
	if err != nil {
		slog.Error("批量查询SKU图片失败", slog.Any("error", err))
		return nil
	}
	attrMap, err := s.skuRepo.BatchGetSkuAttrValues(ctx, skuIDs)
	if err != nil {
		slog.Error("批量查询SKU平台属性值失败", slog.Any("error", err))
		return nil
	}
	saleMap, err := s.skuRepo.BatchGetSkuSaleAttrValues(ctx, skuIDs)
	if err != nil {
		slog.Error("批量查询SKU销售属性值失败", slog.Any("error", err))
		return nil
	}

	// 按 skuID 拼装成 DTO，保持与入参相同的顺序
	res := make([]*types.ResponseSkuInfo, 0, len(skuList))
	for _, sku := range skuList {
		res = append(res, toResponseSkuInfo(sku, imgMap[sku.SkuID], attrMap[sku.SkuID], saleMap[sku.SkuID]))
	}
	return res
}

func (s *skuService) OnSaleSku(ctx context.Context, skuId int64) (err error) {
	if err = s.skuRepo.OnSaleSku(ctx, skuId); err != nil {
		slog.Error("上架SKU失败", slog.Int64("sku_id", skuId), slog.Any("error", err))
		return errs.Wrap(result.CodeSkuDBErr, err)
	}
	return nil
}

func (s *skuService) CancelSaleSku(ctx context.Context, skuId int64) error {
	if err := s.spuRepo.CancelSaleSku(ctx, skuId); err != nil {
		slog.Error("下架SKU失败", slog.Int64("sku_id", skuId), slog.Any("error", err))
		return errs.Wrap(result.CodeSkuDBErr, err)
	}
	return nil
}

func (s *skuService) DeleteSku(ctx context.Context, skuId int64) error {
	if err := s.skuRepo.DeleteSku(ctx, skuId); err != nil {
		slog.Error("删除SKU失败", slog.Int64("sku_id", skuId), slog.Any("error", err))
		return errs.Wrap(result.CodeSkuDBErr, err)
	}
	return nil
}

func (s *skuService) GetSkuInfo(ctx context.Context, skuId int64) (skuInfo *types.ResponseSkuInfo, err error) {
	// repo 返回纯 model，详情页同样用批量方法拉取关联数据（单条 SKU 时用 IN 也仅 3 次 SQL）。
	sku, err := s.skuRepo.GetSku(ctx, skuId)
	if err != nil {
		slog.Error("查询SKU失败", slog.Int64("sku_id", skuId), slog.Any("error", err))
		return nil, errs.Wrap(result.CodeSkuDBErr, err)
	}

	imgMap, err := s.skuRepo.BatchGetSkuImages(ctx, []int64{skuId})
	if err != nil {
		slog.Error("查询SKU图片失败", slog.Int64("sku_id", skuId), slog.Any("error", err))
		return nil, errs.Wrap(result.CodeSkuDBErr, err)
	}
	attrMap, err := s.skuRepo.BatchGetSkuAttrValues(ctx, []int64{skuId})
	if err != nil {
		slog.Error("查询SKU平台属性值失败", slog.Int64("sku_id", skuId), slog.Any("error", err))
		return nil, errs.Wrap(result.CodeSkuDBErr, err)
	}
	saleMap, err := s.skuRepo.BatchGetSkuSaleAttrValues(ctx, []int64{skuId})
	if err != nil {
		slog.Error("查询SKU销售属性值失败", slog.Int64("sku_id", skuId), slog.Any("error", err))
		return nil, errs.Wrap(result.CodeSkuDBErr, err)
	}

	skuInfo = toResponseSkuInfo(sku, imgMap[skuId], attrMap[skuId], saleMap[skuId])
	return skuInfo, nil
}
