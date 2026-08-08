package service

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/pkg/errs"
	"github.com/liwook/go-vue-selection/pkg/idconv"
	"github.com/liwook/go-vue-selection/pkg/result"
	"github.com/liwook/go-vue-selection/pkg/snowflake"
	"github.com/liwook/go-vue-selection/repository"
	"github.com/liwook/go-vue-selection/types"
)

type categoryService struct {
	repo *repository.CategoryRepo
}

func NewCategoryService(repo *repository.CategoryRepo) *categoryService {
	return &categoryService{repo: repo}
}

func (c *categoryService) GetCategory1(ctx context.Context) (data []*types.Category1, err error) {
	rows, err := c.repo.GetCategory1List(ctx)
	if err != nil {
		slog.Error("查询一级分类失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeCategoryDBErr, err)
	}
	data = make([]*types.Category1, 0, len(rows))
	for _, r := range rows {
		data = append(data, modelToCategory1(r))
	}
	return
}

func (c *categoryService) GetCategory2(ctx context.Context, category1Id int64) (data []*types.Category2, err error) {
	rows, err := c.repo.GetCategory2List(ctx, category1Id)
	if err != nil {
		slog.Error("查询二级分类失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeCategoryDBErr, err)
	}
	data = make([]*types.Category2, 0, len(rows))
	for _, r := range rows {
		data = append(data, modelToCategory2(r))
	}
	return
}

func (c *categoryService) GetCategory3(ctx context.Context, category2Id int64) (data []*types.Category3, err error) {
	rows, err := c.repo.GetCategory3List(ctx, category2Id)
	if err != nil {
		slog.Error("查询三级分类失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeCategoryDBErr, err)
	}
	data = make([]*types.Category3, 0, len(rows))
	for _, r := range rows {
		data = append(data, modelToCategory3(r))
	}
	return
}

// CreateCategory2 新增二级分类，返回包含雪花算法生成 ID 的新分类。
func (c *categoryService) CreateCategory2(ctx context.Context, p *types.ParamC2Create) (data *types.Category2, err error) {
	m := &model.Category2{
		Category2ID: snowflake.GenID(), // 原生 int64
		Name:        p.Name,
		Category1ID: idconv.ToInt64Safe(p.Category1ID),
	}
	if m, err = c.repo.InsertCategory2(ctx, m); err != nil {
		slog.Error("创建二级分类失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeCategoryDBErr, err)
	}
	data = &types.Category2{
		Category2ID: strconv.FormatInt(m.Category2ID, 10),
		Name:        m.Name,
		Category1ID: strconv.FormatInt(m.Category1ID, 10),
	}
	return data, nil
}

// CreateCategory3 新增三级分类，返回包含雪花算法生成 ID 的新分类。
func (c *categoryService) CreateCategory3(ctx context.Context, p *types.ParamC3Create) (data *types.Category3, err error) {
	m := &model.Category3{
		Category3ID: snowflake.GenID(), // 原生 int64
		Name:        p.Name,
		Category2ID: idconv.ToInt64Safe(p.Category2ID),
	}
	if m, err = c.repo.InsertCategory3(ctx, m); err != nil {
		slog.Error("创建三级分类失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeCategoryDBErr, err)
	}
	data = &types.Category3{
		Category3ID: strconv.FormatInt(m.Category3ID, 10),
		Name:        m.Name,
		Category2ID: strconv.FormatInt(m.Category2ID, 10),
	}
	return data, nil
}
