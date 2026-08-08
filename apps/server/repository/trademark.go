package repository

import (
	"context"
	"errors"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/dal/query"

	"gorm.io/gorm"
)

// ErrorTrademarkExist 品牌已存在错误
var ErrorTrademarkExist = errors.New("品牌已经存在")

type TrademarkRepo struct {
	db *gorm.DB
	q  *query.Query
}

func NewTrademarkRepo(db *gorm.DB) *TrademarkRepo {
	return &TrademarkRepo{db: db, q: query.Use(db)}
}

func (d *TrademarkRepo) CheckTrademarkExist(ctx context.Context, tmName string) (err error) {
	count, err := d.q.Trademark.WithContext(ctx).Where(d.q.Trademark.TmName.Eq(tmName)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrorTrademarkExist
	}
	return nil
}

func (d *TrademarkRepo) InsertTrademark(ctx context.Context, tm *model.Trademark) (err error) {
	return d.q.Trademark.WithContext(ctx).Create(tm)
}

func (d *TrademarkRepo) GetTrademarkList(ctx context.Context, page, limit int64) (list []*model.Trademark, total int64, err error) {
	total, err = d.q.Trademark.WithContext(ctx).Count()
	if err != nil {
		return nil, 0, err
	}
	if page > 0 && limit > 0 {
		list, err = d.q.Trademark.WithContext(ctx).
			Offset(int((page - 1) * limit)).
			Limit(int(limit)).
			Find()
	} else {
		list, err = d.q.Trademark.WithContext(ctx).Find()
	}
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (d *TrademarkRepo) UpdateTrademark(ctx context.Context, tm *model.Trademark) (err error) {
	_, err = d.q.Trademark.WithContext(ctx).
		Where(d.q.Trademark.TmID.Eq(tm.TmID)).
		Updates(tm)
	return err
}

func (d *TrademarkRepo) DeleteTrademark(ctx context.Context, tmId int64) (err error) {
	_, err = d.q.Trademark.WithContext(ctx).Where(d.q.Trademark.TmID.Eq(tmId)).Delete()
	return err
}

func (d *TrademarkRepo) GetAllTrademarkList(ctx context.Context) (list []*model.Trademark, err error) {
	list, err = d.q.Trademark.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}
	return list, nil
}
