package repository

import (
	"context"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/dal/query"

	"gorm.io/gorm"
)

type CategoryRepo struct {
	db *gorm.DB
	q  *query.Query
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{db: db, q: query.Use(db)}
}

func (d *CategoryRepo) GetCategory1List(ctx context.Context) (list []*model.Category1, err error) {
	rows, err := d.q.Category1.WithContext(ctx).Order(d.q.Category1.Category1ID).Find()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *CategoryRepo) GetCategory2List(ctx context.Context, category1Id int64) (list []*model.Category2, err error) {
	rows, err := d.q.Category2.WithContext(ctx).
		Where(d.q.Category2.Category1ID.Eq(category1Id)).
		Order(d.q.Category2.Category2ID).Find()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *CategoryRepo) GetCategory3List(ctx context.Context, category2Id int64) (list []*model.Category3, err error) {
	rows, err := d.q.Category3.WithContext(ctx).
		Where(d.q.Category3.Category2ID.Eq(category2Id)).
		Order(d.q.Category3.Category3ID).Find()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// InsertCategory2 新增二级分类。Category2ID 由调用方（service 层）用雪花算法生成后传入，
// 本方法仅负责落库并原样回传（含 ID），供上层直接返回给前端。
func (d *CategoryRepo) InsertCategory2(ctx context.Context, c *model.Category2) (created *model.Category2, err error) {
	if err = d.q.Category2.WithContext(ctx).Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

// InsertCategory3 新增三级分类。Category3ID 由调用方（service 层）用雪花算法生成后传入，
// 本方法仅负责落库并原样回传（含 ID），供上层直接返回给前端。
func (d *CategoryRepo) InsertCategory3(ctx context.Context, c *model.Category3) (created *model.Category3, err error) {
	if err = d.q.Category3.WithContext(ctx).Create(c); err != nil {
		return nil, err
	}
	return c, nil
}
