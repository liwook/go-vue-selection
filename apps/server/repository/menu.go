package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/dal/query"
	"github.com/liwook/go-vue-selection/pkg/idconv"

	"gorm.io/gorm"
)

type MenuRepo struct {
	db *gorm.DB
	q  *query.Query
}

// ErrorMenuExist 菜单名称已存在（创建/更新时的业务唯一性冲突）。
var ErrorMenuExist = errors.New("菜单名称已存在")

func NewMenuRepo(db *gorm.DB) *MenuRepo {
	return &MenuRepo{db: db, q: query.Use(db)}
}

func (d *MenuRepo) GetMenuList(ctx context.Context) (menuList []*model.Menu, err error) {
	menus, err := d.q.Menu.WithContext(ctx).Order(d.q.Menu.MenuID).Find()
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (d *MenuRepo) GetMenu(ctx context.Context, menuID int64) (menu *model.Menu, err error) {
	m, err := d.q.Menu.WithContext(ctx).Where(d.q.Menu.MenuID.Eq(menuID)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return m, nil
}

func (d *MenuRepo) CheckMenuExist(ctx context.Context, name string) (err error) {
	count, err := d.q.Menu.WithContext(ctx).Where(d.q.Menu.Name.Eq(name)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrorMenuExist
	}
	return nil
}

func (d *MenuRepo) InsertMenu(ctx context.Context, m *model.Menu) (err error) {
	return d.q.Menu.WithContext(ctx).Create(m)
}

func (d *MenuRepo) UpdateMenu(ctx context.Context, m *model.Menu) (err error) {
	_, err = d.q.Menu.WithContext(ctx).
		Where(d.q.Menu.MenuID.Eq(m.MenuID)).
		Updates(m)
	return err
}

func (d *MenuRepo) QuerySubMenuByID(ctx context.Context, menuID int64) (total int64, err error) {
	return d.q.Menu.WithContext(ctx).Where(d.q.Menu.Pid.Eq(menuID)).Count()
}

func (d *MenuRepo) DeleteMenuByID(ctx context.Context, menuID int64) (err error) {
	_, err = d.q.Menu.WithContext(ctx).Where(d.q.Menu.MenuID.Eq(menuID)).Delete()
	return err
}

func (d *MenuRepo) GetAssignMenu(ctx context.Context, roleId int64) (menuIds []string, err error) {
	rels, err := d.q.RoleMenu.WithContext(ctx).Where(d.q.RoleMenu.RoleID.Eq(roleId)).Find()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rels))
	for _, r := range rels {
		ids = append(ids, idconv.ToStr(r.MenuID))
	}
	return ids, nil
}

func (d *MenuRepo) DeleteAssignMenuByRoleId(ctx context.Context, roleId int64) (err error) {
	_, err = d.q.RoleMenu.WithContext(ctx).Where(d.q.RoleMenu.RoleID.Eq(roleId)).Delete()
	return err
}

func (d *MenuRepo) DoAssign(ctx context.Context, roleId int64, menuIds []int64) (err error) {
	// 使用事务保证原子性：先删除角色原菜单关联，再批量插入新关联。
	// 两段写任一失败整体回滚，避免「原关联已删、新关联未插入」导致角色权限丢失。
	return d.q.Transaction(func(tx *query.Query) error {
		_, err := tx.RoleMenu.WithContext(ctx).Where(tx.RoleMenu.RoleID.Eq(roleId)).Delete()
		if err != nil {
			return fmt.Errorf("分配菜单权限失败(删除原菜单关联): %w", err)
		}
		if len(menuIds) == 0 {
			return nil
		}
		rows := make([]*model.RoleMenu, 0, len(menuIds))
		for _, id := range menuIds {
			rows = append(rows, &model.RoleMenu{
				RoleID: roleId,
				MenuID: id,
			})
		}
		if err := tx.RoleMenu.WithContext(ctx).Create(rows...); err != nil {
			return fmt.Errorf("分配菜单权限失败(批量插入菜单关联): %w", err)
		}
		return nil
	})
}
