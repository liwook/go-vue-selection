package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/dal/query"

	"gorm.io/gorm"
)

type SkuRepo struct {
	db *gorm.DB
	q  *query.Query
}

func NewSkuRepo(db *gorm.DB) *SkuRepo {
	return &SkuRepo{db: db, q: query.Use(db)}
}

// 说明：model→types（DTO）的转换逻辑已下沉到 service 层（见 service/sku.go 的 toResponse* 函数）。
// repo 层只返回纯 model，不再依赖 types 包。这样做的好处：
//   1) 分层干净：repo 不反向依赖对外 JSON 契约 types，未来把 repo 抽象成接口时，
//      接口签名只需依赖 dal/model，service 作为消费方定义接口，不会出现 service 反向依赖 repo 结构体的问题；
//   2) 列表/批量查询在 repo 侧统一用批量方法返回，配合 service 的拼装即可消除 N+1 查询。

func (d *SkuRepo) SaveSkuInfo(ctx context.Context, sku *model.Sku, skuImageList []*model.SkuImage, skuAttrValueList []*model.SkuAttrValue, skuSaleAttrValueList []*model.SkuSaleAttrValue) error {
	// 使用事务保证原子性：连续写入 sku 父表与多张子表，任何一步失败都会整体回滚，
	// 避免出现「父表已插入但子表缺失」的残缺数据。
	return d.q.Transaction(func(tx *query.Query) error {
		if err := tx.Sku.WithContext(ctx).Create(sku); err != nil {
			return fmt.Errorf("保存SKU失败(插入SKU): %w", err)
		}

		if len(skuImageList) > 0 {
			if err := tx.SkuImage.WithContext(ctx).Create(skuImageList...); err != nil {
				return fmt.Errorf("保存SKU失败(批量插入SKU图片): %w", err)
			}
		}

		if len(skuAttrValueList) > 0 {
			if err := tx.SkuAttrValue.WithContext(ctx).Create(skuAttrValueList...); err != nil {
				return fmt.Errorf("保存SKU失败(批量插入SKU平台属性值): %w", err)
			}
		}

		if len(skuSaleAttrValueList) > 0 {
			if err := tx.SkuSaleAttrValue.WithContext(ctx).Create(skuSaleAttrValueList...); err != nil {
				return fmt.Errorf("保存SKU失败(批量插入SKU销售属性值): %w", err)
			}
		}

		return nil
	})
}

// GetSkuList 返回 SKU model 列表（不含关联的子表数据）。
// 注意：这里只查主表 sku，关联的图片/平台属性值/销售属性值由 service 层
// 通过 BatchGet* 批量方法一次性拉取后拼装，避免逐条 SKU 各发 3 次查询的 N+1 问题。
func (d *SkuRepo) GetSkuList(ctx context.Context, page, limit int64) ([]*model.Sku, int64, error) {
	count, err := d.q.Sku.WithContext(ctx).Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*model.Sku
	if page > 0 && limit > 0 {
		list, err = d.q.Sku.WithContext(ctx).
			Offset(int((page - 1) * limit)).
			Limit(int(limit)).
			Find()
	} else {
		list, err = d.q.Sku.WithContext(ctx).Find()
	}
	if err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

// FindBySpuId 按 spu_id 返回 SKU model 列表（不含关联子表数据），同样交由 service 批量拼装。
func (d *SkuRepo) FindBySpuId(ctx context.Context, spuId int64) ([]*model.Sku, error) {
	list, err := d.q.Sku.WithContext(ctx).Where(d.q.Sku.SpuID.Eq(spuId)).Find()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// BatchGetSkuImages 批量查询 SKU 图片，按 sku_id 分组返回。
// 一次 SQL（WHERE sku_id IN (...)）即可覆盖整页 SKU，替代逐条 SKU 各查一次的 N+1。
func (d *SkuRepo) BatchGetSkuImages(ctx context.Context, skuIDs []int64) (map[int64][]*model.SkuImage, error) {
	res := make(map[int64][]*model.SkuImage)
	if len(skuIDs) == 0 {
		return res, nil
	}
	list, err := d.q.SkuImage.WithContext(ctx).
		Where(d.q.SkuImage.SkuID.In(skuIDs...)).Find()
	if err != nil {
		return nil, fmt.Errorf("批量查询SKU图片失败: %w", err)
	}
	for _, item := range list {
		res[item.SkuID] = append(res[item.SkuID], item)
	}
	return res, nil
}

// BatchGetSkuAttrValues 批量查询 SKU 平台属性值，按 sku_id 分组返回。
func (d *SkuRepo) BatchGetSkuAttrValues(ctx context.Context, skuIDs []int64) (map[int64][]*model.SkuAttrValue, error) {
	res := make(map[int64][]*model.SkuAttrValue)
	if len(skuIDs) == 0 {
		return res, nil
	}
	list, err := d.q.SkuAttrValue.WithContext(ctx).
		Where(d.q.SkuAttrValue.SkuID.In(skuIDs...)).Find()
	if err != nil {
		return nil, fmt.Errorf("批量查询SKU平台属性值失败: %w", err)
	}
	for _, item := range list {
		res[item.SkuID] = append(res[item.SkuID], item)
	}
	return res, nil
}

// BatchGetSkuSaleAttrValues 批量查询 SKU 销售属性值，按 sku_id 分组返回。
func (d *SkuRepo) BatchGetSkuSaleAttrValues(ctx context.Context, skuIDs []int64) (map[int64][]*model.SkuSaleAttrValue, error) {
	res := make(map[int64][]*model.SkuSaleAttrValue)
	if len(skuIDs) == 0 {
		return res, nil
	}
	list, err := d.q.SkuSaleAttrValue.WithContext(ctx).
		Where(d.q.SkuSaleAttrValue.SkuID.In(skuIDs...)).Find()
	if err != nil {
		return nil, fmt.Errorf("批量查询SKU销售属性值失败: %w", err)
	}
	for _, item := range list {
		res[item.SkuID] = append(res[item.SkuID], item)
	}
	return res, nil
}

func (d *SkuRepo) OnSaleSku(ctx context.Context, skuId int64) (err error) {
	_, err = d.q.Sku.WithContext(ctx).
		Where(d.q.Sku.SkuID.Eq(skuId)).
		Update(d.q.Sku.IsSale, true)
	return err
}

// DeleteSku 删除 SKU 主表记录即可。
// sku_image / sku_attr_value / sku_sale_attr_value 三张从表在 init.sql 中已通过
// ON DELETE CASCADE 外键（fk_sku_img_sku / fk_sku_attr_val_sku / fk_sku_sale_av_sku）
// 级联删除，无需在代码里逐个手动清理，避免冗余的 3 次 DELETE。
func (d *SkuRepo) DeleteSku(ctx context.Context, skuId int64) (err error) {
	_, err = d.q.Sku.WithContext(ctx).Where(d.q.Sku.SkuID.Eq(skuId)).Delete()
	if err != nil {
		return fmt.Errorf("删除SKU失败(删除SKU): %w", err)
	}
	return nil
}

// GetSku 返回单条 SKU 的 model（不含关联子表数据）。
// 详情页（GetSkuInfo）由 service 层继续调用 BatchGet* 拉取关联数据。
// ErrSkuNotFound 表示 SKU 不存在（查无记录）。
var ErrSkuNotFound = errors.New("sku not found")

func (d *SkuRepo) GetSku(ctx context.Context, skuId int64) (sku *model.Sku, err error) {
	skuList, err := d.q.Sku.WithContext(ctx).Where(d.q.Sku.SkuID.Eq(skuId)).Find()
	if err != nil {
		slog.Error("SkuRepo.GetSku", slog.Any("error", err))
		return nil, err
	}
	if len(skuList) == 0 {
		return nil, ErrSkuNotFound
	}
	return skuList[0], nil
}

// GetAttrNameByAttrId 只投影 attr_name 列，避免 SELECT * 拉取整行。
func (d *SkuRepo) GetAttrNameByAttrId(ctx context.Context, attrId int64) (attrName string, err error) {
	a, err := d.q.Attr.WithContext(ctx).
		Select(d.q.Attr.AttrName).
		Where(d.q.Attr.AttrID.Eq(attrId)).
		First()
	if err != nil {
		return "", err
	}
	return a.AttrName, nil
}

// GetAttrValueNameByValueId 只投影 value_name 列。
func (d *SkuRepo) GetAttrValueNameByValueId(ctx context.Context, valueId int64) (attrValueName string, err error) {
	v, err := d.q.AttrValue.WithContext(ctx).
		Select(d.q.AttrValue.ValueName).
		Where(d.q.AttrValue.AttrValueID.Eq(valueId)).
		First()
	if err != nil {
		return "", err
	}
	return v.ValueName, nil
}

// GetSaleAttrNameBySaleAttrId 只投影 sale_attr_name 列。
func (d *SkuRepo) GetSaleAttrNameBySaleAttrId(ctx context.Context, saleAttrId int64) (saleAttrName string, err error) {
	s, err := d.q.SpuSaleAttr.WithContext(ctx).
		Select(d.q.SpuSaleAttr.SaleAttrName).
		Where(d.q.SpuSaleAttr.SpuSaleAttrID.Eq(saleAttrId)).
		First()
	if err != nil {
		return "", err
	}
	return derefStr(s.SaleAttrName), nil
}

// GetSaleAttrValueNameBySaleAttrValueId 只投影 sale_attr_value_name 列。
func (d *SkuRepo) GetSaleAttrValueNameBySaleAttrValueId(ctx context.Context, saleAttrValueId int64) (saleAttrValueName string, err error) {
	v, err := d.q.SaleAttrValue.WithContext(ctx).
		Select(d.q.SaleAttrValue.SaleAttrValueName).
		Where(d.q.SaleAttrValue.SaleAttrValueID.Eq(saleAttrValueId)).
		First()
	if err != nil {
		return "", err
	}
	return v.SaleAttrValueName, nil
}
