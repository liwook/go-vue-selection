package repository

import (
	"context"
	"errors"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/dal/query"

	"gorm.io/gorm"
)

type RoleRepo struct {
	db *gorm.DB
	q  *query.Query
}

// ErrorRoleExist 角色名已存在（创建/更新时的业务唯一性冲突）。
var ErrorRoleExist = errors.New("角色已存在")

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db: db, q: query.Use(db)}
}

func (d *RoleRepo) GetAllRoleList(ctx context.Context) (roleList []*model.Role, err error) {
	roles, err := d.q.Role.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (d *RoleRepo) GetRoleList(ctx context.Context, roleName string, page, limit int64) (roles []*model.Role, total int64, err error) {
	q := d.q.Role.WithContext(ctx)
	if roleName != "" {
		q = q.Where(d.q.Role.RoleName.Like("%" + roleName + "%"))
	}
	total, err = q.Count()
	if err != nil {
		return nil, 0, err
	}
	if limit > 0 {
		offset := (page - 1) * limit
		q = q.Offset(int(offset)).Limit(int(limit))
	}
	roles, err = q.Order(d.q.Role.RoleID).Find()
	if err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func (d *RoleRepo) CheckRoleExist(ctx context.Context, roleName string) (err error) {
	count, err := d.q.Role.WithContext(ctx).Where(d.q.Role.RoleName.Eq(roleName)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrorRoleExist
	}
	return nil
}

func (d *RoleRepo) InsertRole(ctx context.Context, r *model.Role) (err error) {
	return d.q.Role.WithContext(ctx).Create(r)
}

func (d *RoleRepo) UpdateRole(ctx context.Context, r *model.Role) (err error) {
	_, err = d.q.Role.WithContext(ctx).
		Where(d.q.Role.RoleID.Eq(r.RoleID)).
		Updates(r)
	return err
}

func (d *RoleRepo) DeleteRole(ctx context.Context, roleId int64) (err error) {
	_, err = d.q.Role.WithContext(ctx).Where(d.q.Role.RoleID.Eq(roleId)).Delete()
	return err
}
