package repository

import (
	"context"
	"fmt"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/dal/query"

	"gorm.io/gorm"
)

type SpuRepo struct {
	db *gorm.DB
	q  *query.Query
}

func NewSpuRepo(db *gorm.DB) *SpuRepo {
	return &SpuRepo{db: db, q: query.Use(db)}
}

// 说明：model→types（DTO）的转换逻辑已下沉到 service 层（见 service/spu.go 的 toResponse* 函数）。
// repo 层只返回纯 model，不再依赖 types 包。这样做的好处：
//   1) 分层干净：repo 不反向依赖对外 JSON 契约 types，未来把 repo 抽象成接口时，
//      接口签名只需依赖 dal/model，service 作为消费方定义接口，不会出现 service 反向依赖 repo 结构体的问题；
//   2) 列表/批量查询在 repo 侧统一用批量方法返回，配合 service 的拼装即可消除 N+1 查询。

func (d *SpuRepo) GetSaleAttrList(ctx context.Context) (data []*model.SaleAttr, err error) {
	list, err := d.q.SaleAttr.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (d *SpuRepo) SaveSpuInfo(ctx context.Context, spu *model.Spu, imageList []*model.SpuImageList, spuSaleAttrList []*model.SpuSaleAttr, spuSaleAttrValueList []*model.SaleAttrValue) error {
	// 使用事务保证原子性：连续写入 spu 父表与多张子表，任何一步失败都会整体回滚，
	// 避免出现「父表已插入但子表缺失」的残缺数据。
	return d.q.Transaction(func(tx *query.Query) error {
		if err := tx.Spu.WithContext(ctx).Create(spu); err != nil {
			return fmt.Errorf("保存SPU失败(插入SPU): %w", err)
		}

		if len(imageList) > 0 {
			if err := tx.SpuImageList.WithContext(ctx).Create(imageList...); err != nil {
				return fmt.Errorf("保存SPU失败(批量插入SPU图片): %w", err)
			}
		}

		if len(spuSaleAttrList) > 0 {
			if err := tx.SpuSaleAttr.WithContext(ctx).Create(spuSaleAttrList...); err != nil {
				return fmt.Errorf("保存SPU失败(批量插入SPU销售属性): %w", err)
			}
		}

		if len(spuSaleAttrValueList) > 0 {
			if err := tx.SaleAttrValue.WithContext(ctx).Create(spuSaleAttrValueList...); err != nil {
				return fmt.Errorf("保存SPU失败(批量插入SPU销售属性值): %w", err)
			}
		}

		return nil
	})
}

func (d *SpuRepo) GetSpuList(ctx context.Context, c3Id, page, limit int64) (spuList []*model.Spu, count int64, err error) {
	count, err = d.q.Spu.WithContext(ctx).Where(d.q.Spu.Category3ID.Eq(c3Id)).Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*model.Spu
	if page > 0 && limit > 0 {
		list, err = d.q.Spu.WithContext(ctx).
			Where(d.q.Spu.Category3ID.Eq(c3Id)).
			Offset(int((page - 1) * limit)).
			Limit(int(limit)).
			Find()
	} else {
		list, err = d.q.Spu.WithContext(ctx).Where(d.q.Spu.Category3ID.Eq(c3Id)).Find()
	}
	if err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (d *SpuRepo) GetSpuImageList(ctx context.Context, spuId int64) (spuImageList []*model.SpuImageList, err error) {
	list, err := d.q.SpuImageList.WithContext(ctx).Where(d.q.SpuImageList.SpuID.Eq(spuId)).Find()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (d *SpuRepo) GetSpuSaleAttrList(ctx context.Context, spuId int64) (spuSaleAttrList []*model.SpuSaleAttr, err error) {
	list, err := d.q.SpuSaleAttr.WithContext(ctx).Where(d.q.SpuSaleAttr.SpuID.Eq(spuId)).Find()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (d *SpuRepo) GetSpuSaleAttrValueList(ctx context.Context, spuId, baseSaleAttrId int64) (saleAttrValueList []*model.SaleAttrValue, err error) {
	list, err := d.q.SaleAttrValue.WithContext(ctx).
		Where(d.q.SaleAttrValue.SpuID.Eq(spuId), d.q.SaleAttrValue.SaleAttrID.Eq(baseSaleAttrId)).
		Find()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// BatchGetSpuSaleAttrValues 批量查询指定 spu 下所有销售属性值，按 sale_attr_id 分组返回。
// 一次 SQL（WHERE spu_id = ?）即可覆盖该 spu 全部销售属性值的拼装，
// 替代逐条销售属性各查一次的 N+1（见 service.GetSpuSaleAttrList 的组装逻辑）。
func (d *SpuRepo) BatchGetSpuSaleAttrValues(ctx context.Context, spuId int64) (map[int64][]*model.SaleAttrValue, error) {
	res := make(map[int64][]*model.SaleAttrValue)
	list, err := d.q.SaleAttrValue.WithContext(ctx).Where(d.q.SaleAttrValue.SpuID.Eq(spuId)).Find()
	if err != nil {
		return nil, fmt.Errorf("批量查询SPU销售属性值失败: %w", err)
	}
	for _, item := range list {
		res[item.SaleAttrID] = append(res[item.SaleAttrID], item)
	}
	return res, nil
}

func (d *SpuRepo) UpdateSpuInfo(ctx context.Context, spu *model.Spu, imageList []*model.SpuImageList, spuSaleAttrList []*model.SpuSaleAttr, spuSaleAttrValueList []*model.SaleAttrValue) error {
	// 使用事务保证原子性：先删除该 SPU 下全部旧子表行（spu_image_list / spu_sale_attr /
	// sale_attr_value，三张子表在 init.sql 中均无外键级联，不会随父表更新自动清理），
	// 再更新 spu 父表并批量插入新子表行，避免每次保存都让子表记录无限翻倍累积。
	return d.q.Transaction(func(tx *query.Query) error {
		if _, err := tx.SpuImageList.WithContext(ctx).Where(tx.SpuImageList.SpuID.Eq(spu.SpuID)).Delete(); err != nil {
			return fmt.Errorf("更新SPU失败(删除旧SPU图片): %w", err)
		}
		if _, err := tx.SpuSaleAttr.WithContext(ctx).Where(tx.SpuSaleAttr.SpuID.Eq(spu.SpuID)).Delete(); err != nil {
			return fmt.Errorf("更新SPU失败(删除旧SPU销售属性): %w", err)
		}
		if _, err := tx.SaleAttrValue.WithContext(ctx).Where(tx.SaleAttrValue.SpuID.Eq(spu.SpuID)).Delete(); err != nil {
			return fmt.Errorf("更新SPU失败(删除旧SPU销售属性值): %w", err)
		}

		_, err := tx.Spu.WithContext(ctx).
			Where(tx.Spu.SpuID.Eq(spu.SpuID)).
			Updates(map[string]any{
				"spu_name":     spu.SpuName,
				"description":  spu.Description,
				"category3_id": spu.Category3ID,
				"tm_id":        spu.TmID,
			})
		if err != nil {
			return fmt.Errorf("更新SPU失败(更新SPU): %w", err)
		}

		if len(imageList) > 0 {
			if err := tx.SpuImageList.WithContext(ctx).Create(imageList...); err != nil {
				return fmt.Errorf("更新SPU失败(批量插入SPU图片): %w", err)
			}
		}

		if len(spuSaleAttrList) > 0 {
			if err := tx.SpuSaleAttr.WithContext(ctx).Create(spuSaleAttrList...); err != nil {
				return fmt.Errorf("更新SPU失败(批量插入SPU销售属性): %w", err)
			}
		}

		if len(spuSaleAttrValueList) > 0 {
			if err := tx.SaleAttrValue.WithContext(ctx).Create(spuSaleAttrValueList...); err != nil {
				return fmt.Errorf("更新SPU失败(批量插入SPU销售属性值): %w", err)
			}
		}

		return nil
	})
}

func (d *SpuRepo) DeleteSpu(ctx context.Context, spuId int64) (err error) {
	// 三张子表（spu_image_list / spu_sale_attr / sale_attr_value）在 init.sql 中均未配置
	// 外键级联，不会随父表 spu 删除而自动清理，故需手动按 spu_id 先删子表再删父表。
	// 整个删除用事务包裹（SaleAttrValue 删除 + Spu 删除），任一失败整体回滚，
	// 避免只删掉一半留下孤儿数据。
	_, err = d.q.SaleAttrValue.WithContext(ctx).Where(d.q.SaleAttrValue.SpuID.Eq(spuId)).Delete()
	if err != nil {
		return fmt.Errorf("删除SPU失败(删除SPU销售属性值): %w", err)
	}
	_, err = d.q.SpuImageList.WithContext(ctx).Where(d.q.SpuImageList.SpuID.Eq(spuId)).Delete()
	if err != nil {
		return fmt.Errorf("删除SPU失败(删除SPU图片): %w", err)
	}
	_, err = d.q.SpuSaleAttr.WithContext(ctx).Where(d.q.SpuSaleAttr.SpuID.Eq(spuId)).Delete()
	if err != nil {
		return fmt.Errorf("删除SPU失败(删除SPU销售属性): %w", err)
	}
	_, err = d.q.Spu.WithContext(ctx).Where(d.q.Spu.SpuID.Eq(spuId)).Delete()
	if err != nil {
		return fmt.Errorf("删除SPU失败(删除SPU): %w", err)
	}
	return nil
}

func (d *SpuRepo) CancelSaleSku(ctx context.Context, skuId int64) (err error) {
	_, err = d.q.Sku.WithContext(ctx).
		Where(d.q.Sku.SkuID.Eq(skuId)).
		Update(d.q.Sku.IsSale, false)
	return err
}
