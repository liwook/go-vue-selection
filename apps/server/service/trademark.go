package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/pkg/errs"
	"github.com/liwook/go-vue-selection/pkg/idconv"
	"github.com/liwook/go-vue-selection/pkg/result"
	"github.com/liwook/go-vue-selection/pkg/snowflake"
	"github.com/liwook/go-vue-selection/repository"
	"github.com/liwook/go-vue-selection/types"
)

type trademarkService struct {
	repo *repository.TrademarkRepo
}

func NewTrademarkService(repo *repository.TrademarkRepo) *trademarkService {
	return &trademarkService{repo: repo}
}

func (t *trademarkService) CreateTrademark(ctx context.Context, p *types.ParamTmSave) (err error) {
	// 1. 判断品牌是否已经存在
	if err = t.repo.CheckTrademarkExist(ctx, p.TmName); err != nil {
		if errors.Is(err, repository.ErrorTrademarkExist) {
			return errs.Wrap(result.CodeTrademarkErr, err)
		}
		slog.Error("校验品牌是否存在失败", slog.String("tm_name", p.TmName), slog.Any("error", err))
		return errs.Wrap(result.CodeTrademarkDBErr, err)
	}
	// 2. 生成 TmID
	TmID := snowflake.GenID()

	// 3. 构造一个品牌实例
	trademark := &model.Trademark{
		TmID:    TmID, // 原生 int64
		TmName:  p.TmName,
		LogoURL: &p.LogoUrl,
	}

	// 4. 保存进数据库
	err = t.repo.InsertTrademark(ctx, trademark)
	if err != nil {
		slog.Error("创建品牌失败", slog.String("tm_name", p.TmName), slog.Any("error", err))
		return errs.Wrap(result.CodeTrademarkDBErr, err)
	}
	return nil
}

func (t *trademarkService) GetTrademarkList(ctx context.Context, page, limit int64) (data *types.ResponseTmList, err error) {
	list, total, err := t.repo.GetTrademarkList(ctx, page, limit)
	if err != nil {
		slog.Error("查询品牌列表失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeTrademarkDBErr, err)
	}
	records := make([]*types.Trademark, 0, len(list))
	for _, tm := range list {
		records = append(records, modelToTrademark(tm))
	}
	data = &types.ResponseTmList{
		Records:     records,
		Total:       total,
		Size:        limit,
		Current:     page,
		SearchCount: true,
		Pages:       pageCount(total, limit),
	}
	return
}

func (t *trademarkService) UpdateTrademark(ctx context.Context, p *types.ParamTmUpdate) error {
	// 构造一个品牌实例
	trademark := &model.Trademark{
		TmID:    idconv.ToInt64Safe(p.TmID),
		TmName:  p.TmName,
		LogoURL: &p.LogoUrl,
	}

	// 保存进数据库
	err := t.repo.UpdateTrademark(ctx, trademark)
	if err != nil {
		slog.Error("更新品牌失败", slog.String("tm_id", p.TmID), slog.Any("error", err))
		return errs.Wrap(result.CodeTrademarkDBErr, err)
	}
	return nil
}

func (t *trademarkService) DeleteTrademark(ctx context.Context, tmId int64) error {
	if err := t.repo.DeleteTrademark(ctx, tmId); err != nil {
		slog.Error("删除品牌失败", slog.Int64("tm_id", tmId), slog.Any("error", err))
		if errs.IsForeignKeyViolation(err) {
			return errs.Wrap(result.CodeTrademarkInUse, err)
		}
		return errs.Wrap(result.CodeTrademarkDBErr, err)
	}
	return nil
}

func (t *trademarkService) GetAllTrademarkList(ctx context.Context) (data []types.Trademark, err error) {
	list, err := t.repo.GetAllTrademarkList(ctx)
	if err != nil {
		slog.Error("查询全部品牌失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeTrademarkDBErr, err)
	}
	data = make([]types.Trademark, 0, len(list))
	for _, tm := range list {
		data = append(data, *modelToTrademark(tm))
	}
	return data, nil
}
