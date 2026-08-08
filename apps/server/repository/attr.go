package repository

import (
	"context"
	"fmt"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/dal/query"
	"github.com/liwook/go-vue-selection/pkg/snowflake"

	"gorm.io/gorm"
)

type AttrRepo struct {
	db *gorm.DB
	q  *query.Query
}

func NewAttrRepo(db *gorm.DB) *AttrRepo {
	return &AttrRepo{db: db, q: query.Use(db)}
}

func (d *AttrRepo) GetAttr(ctx context.Context, c1Id, c2Id, c3Id int64) (data []*model.Attr, err error) {
	// 按三级分类查询属性
	rows, err := d.q.Attr.WithContext(ctx).
		Where(d.q.Attr.CategoryID.Eq(c3Id)).
		Find()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *AttrRepo) GetAttrValue(ctx context.Context, attrID int64) (data []*model.AttrValue, err error) {
	rows, err := d.q.AttrValue.WithContext(ctx).Where(d.q.AttrValue.AttrID.Eq(attrID)).Find()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *AttrRepo) InsertAttrAndAttrValue(ctx context.Context, attr *model.Attr, values []*model.AttrValue) (err error) {
	// 使用事务保证原子性：先插入 attr 父表，再批量插入 attr_value 子表（依赖父表生成的 attr_id），
	// 任一环节失败整体回滚，避免留下无属性值的孤立 attr 记录。
	return d.q.Transaction(func(tx *query.Query) error {
		if err := tx.Attr.WithContext(ctx).Create(attr); err != nil {
			return fmt.Errorf("新建属性失败(插入属性): %w", err)
		}
		if len(values) > 0 {
			for _, v := range values {
				if v.AttrValueID == 0 {
					v.AttrValueID = snowflake.GenID()
				}
				v.AttrID = attr.AttrID
			}
			if err := tx.AttrValue.WithContext(ctx).Create(values...); err != nil {
				return fmt.Errorf("新建属性失败(批量插入属性值): %w", err)
			}
		}
		return nil
	})
}

func (d *AttrRepo) UpdateAttrAndAttrValue(ctx context.Context, attr *model.Attr, values []*model.AttrValue) (err error) {
	// 使用事务保证原子性：更新 attr 父表 + 逐条更新/插入 attr_value 子表 + 删除过期子表记录，
	// 多步写任一失败整体回滚，避免出现「父表已更新但子表部分丢失/残留」的不一致状态。
	return d.q.Transaction(func(tx *query.Query) error {
		_, err := tx.Attr.WithContext(ctx).
			Where(tx.Attr.AttrID.Eq(attr.AttrID)).
			Updates(map[string]any{"attr_name": attr.AttrName})
		if err != nil {
			return fmt.Errorf("更新属性失败(更新属性名): %w", err)
		}

		keepIDs := make([]int64, 0)
		for _, v := range values {
			if v.AttrValueID != 0 {
				keepIDs = append(keepIDs, v.AttrValueID)
				_, err = tx.AttrValue.WithContext(ctx).
					Where(tx.AttrValue.AttrValueID.Eq(v.AttrValueID)).
					Updates(map[string]any{"value_name": v.ValueName})
				if err != nil {
					return fmt.Errorf("更新属性失败(更新属性值): %w", err)
				}
			} else {
				v.AttrValueID = snowflake.GenID()
				v.AttrID = attr.AttrID
				if err := tx.AttrValue.WithContext(ctx).Create(v); err != nil {
					return fmt.Errorf("更新属性失败(新建属性值): %w", err)
				}
				keepIDs = append(keepIDs, v.AttrValueID)
			}
		}
		// 删除不在列表中的 AttrValue
		if len(keepIDs) > 0 {
			_, err = tx.AttrValue.WithContext(ctx).
				Where(tx.AttrValue.AttrID.Eq(attr.AttrID), tx.AttrValue.AttrValueID.NotIn(keepIDs...)).
				Delete()
		} else {
			_, err = tx.AttrValue.WithContext(ctx).
				Where(tx.AttrValue.AttrID.Eq(attr.AttrID)).
				Delete()
		}
		if err != nil {
			return fmt.Errorf("更新属性失败(删除过期属性值): %w", err)
		}
		return nil
	})
}

func (d *AttrRepo) DeleteAttr(ctx context.Context, attrId int64) (err error) {
	// 使用事务保证原子性：attr_value 与 attr 之间未配置外键（init.sql 仅建了索引），
	// 不会随 attr 级联删除，故先删父表再删子表两段写必须整体事务化，避免只删掉一半。
	return d.q.Transaction(func(tx *query.Query) error {
		_, err := tx.Attr.WithContext(ctx).Where(tx.Attr.AttrID.Eq(attrId)).Delete()
		if err != nil {
			return fmt.Errorf("删除属性失败(删除属性): %w", err)
		}
		_, err = tx.AttrValue.WithContext(ctx).Where(tx.AttrValue.AttrID.Eq(attrId)).Delete()
		if err != nil {
			return fmt.Errorf("删除属性失败(删除属性值): %w", err)
		}
		return nil
	})
}
