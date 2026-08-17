package service

import (
	"context"
	"log/slog"
	"math"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/pkg/errs"
	"github.com/liwook/go-vue-selection/pkg/idconv"
	"github.com/liwook/go-vue-selection/pkg/result"
	"github.com/liwook/go-vue-selection/pkg/snowflake"
	"github.com/liwook/go-vue-selection/repository"
	"github.com/liwook/go-vue-selection/types"
)

type spuService struct {
	repo *repository.SpuRepo
}

func NewSpuService(repo *repository.SpuRepo) *spuService {
	return &spuService{repo: repo}
}

func (s *spuService) GetSaleAttrList(ctx context.Context) ([]*types.SaleAttr, error) {
	list, err := s.repo.GetSaleAttrList(ctx)
	if err != nil {
		slog.Error("查询销售属性失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeSpuDBErr, err)
	}
	res := make([]*types.SaleAttr, 0, len(list))
	for _, a := range list {
		res = append(res, toResponseSaleAttr(a))
	}
	return res, nil
}

func (s *spuService) SaveSpuInfo(ctx context.Context, p *types.Spu) (err error) {
	// SPU
	spuId := snowflake.GenID()
	spu := &model.Spu{
		SpuID:       spuId, // 原生 int64
		SpuName:     p.SpuName,
		Description: p.Description,
		Category3ID: idconv.ToInt64Safe(p.Category3ID),
		TmID:        idconv.ToInt64Safe(p.TmID),
	}

	// 图片列表
	imageList := make([]*model.SpuImageList, 0, len(p.SpuImageList))
	if len(p.SpuImageList) > 0 {
		for _, image := range p.SpuImageList {
			imageId := snowflake.GenID()
			imageList = append(imageList, &model.SpuImageList{
				ImageID:   imageId,
				ImageName: &image.ImageName,
				ImageURL:  &image.ImageUrl,
				SpuID:     spuId,
			})
		}
	}

	// SPU 销售属性
	spuSaleAttrList := make([]*model.SpuSaleAttr, 0, len(p.SpuSaleAttrList))
	if len(p.SpuSaleAttrList) > 0 {
		for _, spuSaleAttr := range p.SpuSaleAttrList {
			saleAttrId := snowflake.GenID()
			spuSaleAttrList = append(spuSaleAttrList, &model.SpuSaleAttr{
				SpuSaleAttrID:  saleAttrId,
				BaseSaleAttrID: idconv.ToInt64Safe(spuSaleAttr.BaseSaleAttrId),
				SaleAttrName:   &spuSaleAttr.SaleAttrName,
				SpuID:          spuId,
			})
		}
	}

	// SPU销售属性值
	saleAttrValueCount := 0
	for _, spuSaleAttr := range p.SpuSaleAttrList {
		saleAttrValueCount += len(spuSaleAttr.SpuSaleAttrValue)
	}
	spuSaleAttrValueList := make([]*model.SaleAttrValue, 0, saleAttrValueCount)
	if len(p.SpuSaleAttrList) > 0 {
		for _, spuSaleAttr := range p.SpuSaleAttrList {
			for _, spuSaleAttrValue := range spuSaleAttr.SpuSaleAttrValue {
				saleAttrValueId := snowflake.GenID()
				spuSaleAttrValueList = append(spuSaleAttrValueList, &model.SaleAttrValue{
					SaleAttrValueID:   saleAttrValueId,
					SaleAttrValueName: spuSaleAttrValue.SaleAttrValueName,
					SaleAttrID:        idconv.ToInt64Safe(spuSaleAttrValue.BaseSaleAttrId),
					SpuID:             spuId,
				})
			}
		}
	}

	err = s.repo.SaveSpuInfo(ctx, spu, imageList, spuSaleAttrList, spuSaleAttrValueList)
	if err != nil {
		slog.Error("保存SPU失败", slog.Any("error", err))
		return errs.Wrap(result.CodeSpuDBErr, err)
	}
	return err
}

func (s *spuService) GetSpuList(ctx context.Context, c3Id, page, limit int64) (spu *types.ResponseSpuList, err error) {
	spuList, count, err := s.repo.GetSpuList(ctx, c3Id, page, limit)
	if err != nil {
		slog.Error("查询SPU列表失败", slog.Int64("c3_id", c3Id), slog.Any("error", err))
		return nil, errs.Wrap(result.CodeSpuDBErr, err)
	}

	records := make([]*types.Spu, 0, len(spuList))
	for _, m := range spuList {
		records = append(records, toResponseSpu(m))
	}

	spu = &types.ResponseSpuList{
		Records:     records,
		Total:       count,
		Size:        limit,
		Current:     page,
		SearchCount: true,
		Pages:       int64(math.Ceil(float64(count) / float64(limit))),
	}

	return spu, err
}

func (s *spuService) GetSpuImageList(ctx context.Context, spuId int64) (spuImageList []*types.SpuImage, err error) {
	list, err := s.repo.GetSpuImageList(ctx, spuId)
	if err != nil {
		slog.Error("查询SPU图片失败", slog.Int64("spu_id", spuId), slog.Any("error", err))
		return nil, errs.Wrap(result.CodeSpuDBErr, err)
	}
	res := make([]*types.SpuImage, 0, len(list))
	for _, m := range list {
		res = append(res, toResponseSpuImage(m))
	}
	return res, nil
}

// GetSpuSaleAttrList 查询 SPU 下的销售属性，并拼装每个销售属性对应的属性值列表。
// 通过 BatchGetSpuSaleAttrValues 一次 SQL 批量取出该 spu 全部销售属性值，
// 按 sale_attr_id 分组后挂到对应销售属性上，避免「逐条销售属性各查一次」的 N+1 问题。
func (s *spuService) GetSpuSaleAttrList(ctx context.Context, spuId int64) (spuSaleAttrList []*types.SpuSaleAttr, err error) {
	modelList, err := s.repo.GetSpuSaleAttrList(ctx, spuId)
	if err != nil {
		slog.Error("查询SPU销售属性失败", slog.Int64("spu_id", spuId), slog.Any("error", err))
		return nil, errs.Wrap(result.CodeSpuDBErr, err)
	}

	// 批量拉取该 spu 下所有销售属性值，按 sale_attr_id 分组
	valueMap, err := s.repo.BatchGetSpuSaleAttrValues(ctx, spuId)
	if err != nil {
		slog.Error("查询SPU销售属性值失败", slog.Int64("spu_id", spuId), slog.Any("error", err))
		return nil, errs.Wrap(result.CodeSpuDBErr, err)
	}

	res := make([]*types.SpuSaleAttr, 0, len(modelList))
	for _, m := range modelList {
		spuSaleAttr := toResponseSpuSaleAttr(m)
		values := valueMap[m.BaseSaleAttrID]
		if len(values) > 0 {
			spuSaleAttr.SpuSaleAttrValue = make([]*types.SaleAttrValue, 0, len(values))
			for _, v := range values {
				spuSaleAttr.SpuSaleAttrValue = append(spuSaleAttr.SpuSaleAttrValue, toResponseSaleAttrValue(v))
			}
		}
		res = append(res, spuSaleAttr)
	}
	return res, nil
}

func (s *spuService) UpdateSpuInfo(ctx context.Context, spu *types.Spu) error {
	// 创建一个 SPU 对象
	newSpu := &model.Spu{
		SpuID:       idconv.ToInt64Safe(spu.SpuID),
		SpuName:     spu.SpuName,
		Description: spu.Description,
		Category3ID: idconv.ToInt64Safe(spu.Category3ID),
		TmID:        idconv.ToInt64Safe(spu.TmID),
	}

	// 创建一个图片列表
	imageList := make([]*model.SpuImageList, 0, len(spu.SpuImageList))
	if len(spu.SpuImageList) > 0 {
		for _, image := range spu.SpuImageList {
			imageId := snowflake.GenID()
			imageList = append(imageList, &model.SpuImageList{
				ImageID:   imageId,
				ImageName: &image.ImageName,
				ImageURL:  &image.ImageUrl,
				SpuID:     idconv.ToInt64Safe(spu.SpuID),
			})
		}
	}

	// SPU 销售属性
	spuSaleAttrList := make([]*model.SpuSaleAttr, 0, len(spu.SpuSaleAttrList))
	if len(spu.SpuSaleAttrList) > 0 {
		for _, spuSaleAttr := range spu.SpuSaleAttrList {
			spuSaleAttrId := snowflake.GenID()
			spuSaleAttrList = append(spuSaleAttrList, &model.SpuSaleAttr{
				SpuSaleAttrID:  spuSaleAttrId,
				BaseSaleAttrID: idconv.ToInt64Safe(spuSaleAttr.BaseSaleAttrId),
				SaleAttrName:   &spuSaleAttr.SaleAttrName,
				SpuID:          idconv.ToInt64Safe(spu.SpuID),
			})
		}
	}

	// SPU销售属性值
	saleAttrValueCount := 0
	for _, spuSaleAttr := range spu.SpuSaleAttrList {
		saleAttrValueCount += len(spuSaleAttr.SpuSaleAttrValue)
	}
	spuSaleAttrValueList := make([]*model.SaleAttrValue, 0, saleAttrValueCount)
	if len(spu.SpuSaleAttrList) > 0 {
		for _, spuSaleAttr := range spu.SpuSaleAttrList {
			if len(spuSaleAttr.SpuSaleAttrValue) > 0 {
				for _, spuSaleAttrValue := range spuSaleAttr.SpuSaleAttrValue {
					saleAttrValueId := snowflake.GenID()
					spuSaleAttrValueList = append(spuSaleAttrValueList, &model.SaleAttrValue{
						SaleAttrValueID:   saleAttrValueId,
						SaleAttrValueName: spuSaleAttrValue.SaleAttrValueName,
						SaleAttrID:        idconv.ToInt64Safe(spuSaleAttrValue.BaseSaleAttrId),
						SpuID:             idconv.ToInt64Safe(spu.SpuID),
					})
				}
			}

		}
	}

	err := s.repo.UpdateSpuInfo(ctx, newSpu, imageList, spuSaleAttrList, spuSaleAttrValueList)
	if err != nil {
		slog.Error("更新SPU失败", slog.String("spu_id", spu.SpuID), slog.Any("error", err))
		return errs.Wrap(result.CodeSpuDBErr, err)
	}
	return nil
}

func (s *spuService) DeleteSpu(ctx context.Context, spuId int64) error {
	if err := s.repo.DeleteSpu(ctx, spuId); err != nil {
		slog.Error("删除SPU失败", slog.Int64("spu_id", spuId), slog.Any("error", err))
		if errs.IsForeignKeyViolation(err) {
			return errs.Wrap(result.CodeSpuInUse, err)
		}
		return errs.Wrap(result.CodeSpuDBErr, err)
	}
	return nil
}

// 以下为 model → types（DTO）的转换函数，与 repo 层解耦：
// repo 只返回纯 model，由 service 在返回给 handler 前完成对外契约的转换。

func toResponseSpu(s *model.Spu) *types.Spu {
	return &types.Spu{
		SpuID:       idconv.ToStr(s.SpuID),
		SpuName:     s.SpuName,
		Description: s.Description,
		Category3ID: idconv.ToStr(s.Category3ID),
		TmID:        idconv.ToStr(s.TmID),
		BaseModel: types.BaseModel{
			CreateTime: s.CreateTime,
			UpdateTime: s.UpdateTime,
		},
	}
}

func toResponseSpuImage(i *model.SpuImageList) *types.SpuImage {
	return &types.SpuImage{
		ImageID:   idconv.ToStr(i.ImageID),
		ImageName: derefStr(i.ImageName),
		ImageUrl:  derefStr(i.ImageURL),
		SpuID:     idconv.ToStr(i.SpuID),
		BaseModel: types.BaseModel{
			CreateTime: i.CreateTime,
			UpdateTime: i.UpdateTime,
		},
	}
}

func toResponseSpuSaleAttr(s *model.SpuSaleAttr) *types.SpuSaleAttr {
	return &types.SpuSaleAttr{
		SpuSaleAttrID:  idconv.ToStr(s.SpuSaleAttrID),
		BaseSaleAttrId: idconv.ToStr(s.BaseSaleAttrID),
		SaleAttrName:   derefStr(s.SaleAttrName),
		SpuID:          idconv.ToStr(s.SpuID),
	}
}

func toResponseSaleAttrValue(v *model.SaleAttrValue) *types.SaleAttrValue {
	return &types.SaleAttrValue{
		SaleAttrValueID:   idconv.ToStr(v.SaleAttrValueID),
		SaleAttrValueName: v.SaleAttrValueName,
		BaseSaleAttrId:    idconv.ToStr(v.SaleAttrID),
		SpuID:             idconv.ToStr(v.SpuID),
		BaseModel: types.BaseModel{
			CreateTime: v.CreateTime,
			UpdateTime: v.UpdateTime,
		},
	}
}

func toResponseSaleAttr(a *model.SaleAttr) *types.SaleAttr {
	return &types.SaleAttr{
		SaleAttrID:   idconv.ToStr(a.SaleAttrID),
		SaleAttrName: a.SaleAttrName,
		BaseModel: types.BaseModel{
			CreateTime: a.CreateTime,
			UpdateTime: a.UpdateTime,
		},
	}
}
