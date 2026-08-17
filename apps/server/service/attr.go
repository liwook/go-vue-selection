package service

import (
	"context"
	"log/slog"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/pkg/errs"
	"github.com/liwook/go-vue-selection/pkg/idconv"
	"github.com/liwook/go-vue-selection/pkg/result"
	"github.com/liwook/go-vue-selection/pkg/snowflake"
	"github.com/liwook/go-vue-selection/repository"
	"github.com/liwook/go-vue-selection/types"
)

type attrService struct {
	repo *repository.AttrRepo
}

func NewAttrService(repo *repository.AttrRepo) *attrService {
	return &attrService{repo: repo}
}

func (a *attrService) UpdateAttr(ctx context.Context, p *types.Attr) (err error) {
	// 把 API 层的 types.Attr 转换为存储层的 model.Attr
	attr := &model.Attr{
		AttrID:     idconv.ToInt64Safe(p.AttrID),
		AttrName:   p.AttrName,
		CategoryID: idconv.ToInt64Safe(p.CategoryID),
	}
	values := make([]*model.AttrValue, 0, len(p.AttrValueList))
	for i := 0; i < len(p.AttrValueList); i++ {
		v := p.AttrValueList[i]
		values = append(values, &model.AttrValue{
			AttrValueID: idconv.ToInt64Safe(v.AttrValueID),
			ValueName:   v.ValueName,
			AttrID:      idconv.ToInt64Safe(v.AttrID),
		})
	}
	err = a.repo.UpdateAttrAndAttrValue(ctx, attr, values)
	if err != nil {
		slog.Error("更新属性失败",
			slog.Int64("attr_id", attr.AttrID),
			slog.Int64("category_id", attr.CategoryID),
			slog.Any("error", err))
		return errs.Wrap(result.CodeAttrDBErr, err)
	}
	return err
}

func (a *attrService) CreateAttr(ctx context.Context, p *types.Attr) (err error) {
	// 创建属性实例
	attr := &model.Attr{
		AttrID:     snowflake.GenID(), // 原生 int64
		AttrName:   p.AttrName,
		CategoryID: idconv.ToInt64Safe(p.CategoryID),
	}

	// 创建属性值实例
	values := make([]*model.AttrValue, 0, len(p.AttrValueList))
	for i := 0; i < len(p.AttrValueList); i++ {
		values = append(values, &model.AttrValue{
			AttrValueID: snowflake.GenID(), // 原生 int64
			ValueName:   p.AttrValueList[i].ValueName,
			AttrID:      attr.AttrID,
		})
	}

	// 写入数据库
	err = a.repo.InsertAttrAndAttrValue(ctx, attr, values)
	if err != nil {
		slog.Error("新建属性失败", slog.Any("error", err))
		return errs.Wrap(result.CodeAttrDBErr, err)
	}

	return err
}

func (a *attrService) GetAttr(ctx context.Context, c1Id, c2Id, c3Id int64) (attrs []*types.Attr, err error) {
	modelAttrs, err := a.repo.GetAttr(ctx, c1Id, c2Id, c3Id)
	if err != nil {
		slog.Error("查询属性失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeAttrDBErr, err)
	}

	for _, ma := range modelAttrs {
		attr := modelToAttr(ma)
		modelValues, err := a.repo.GetAttrValue(ctx, ma.AttrID)
		if err != nil {
			slog.Error("查询属性值失败", slog.Any("error", err))
			return nil, errs.Wrap(result.CodeAttrDBErr, err)
		}
		values := make([]*types.AttrValue, 0, len(modelValues))
		for _, mv := range modelValues {
			values = append(values, modelToAttrValue(mv))
		}
		attr.AttrValueList = values
		attrs = append(attrs, attr)
	}
	return
}

func (a *attrService) DeleteAttr(ctx context.Context, attrId int64) (err error) {
	if err = a.repo.DeleteAttr(ctx, attrId); err != nil {
		slog.Error("删除属性失败", slog.Int64("attr_id", attrId), slog.Any("error", err))
		if errs.IsForeignKeyViolation(err) {
			return errs.Wrap(result.CodeAttrInUse, err)
		}
		return errs.Wrap(result.CodeAttrDBErr, err)
	}
	return nil
}
