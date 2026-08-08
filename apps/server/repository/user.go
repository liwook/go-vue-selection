package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/dal/query"
	"github.com/liwook/go-vue-selection/pkg/idconv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
	q  *query.Query
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db, q: query.Use(db)}
}

// 用户相关错误
var (
	ErrorUserNotExist  = errors.New("用户名或密码错误")
	ErrorUserExist     = errors.New("用户名已存在")
	ErrorPasswordWrong = errors.New("用户名或密码错误")
)

// hashPassword 使用 bcrypt 对明文密码进行哈希
func hashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Login 按用户名查出用户并校验密码与账号状态，返回纯 model.User 供上层生成令牌。
// 只按用户名查出用户，密码用 bcrypt 在应用层比对，避免数据库存明文/可等值比较的密码。
func (d *UserRepo) Login(ctx context.Context, username, password string) (u *model.User, err error) {
	u, err = d.q.User.WithContext(ctx).
		Where(d.q.User.Username.Eq(username)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorUserNotExist
		}
		return nil, fmt.Errorf("登录失败(查询用户): %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, ErrorPasswordWrong
	}
	// 账号被禁用/锁定时拒绝登录，但仍使用统一的“用户名或密码错误”提示，避免泄露账号状态
	if !u.Status {
		return nil, ErrorUserNotExist
	}
	return u, nil
}

func (d *UserRepo) CheckUserExist(ctx context.Context, username string) (err error) {
	count, err := d.q.User.WithContext(ctx).Where(d.q.User.Username.Eq(username)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

func (d *UserRepo) InsertUser(ctx context.Context, user *model.User) (err error) {
	hashed, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	user.Status = true
	if err := d.q.User.WithContext(ctx).Create(user); err != nil {
		return fmt.Errorf("创建用户失败(插入用户): %w", err)
	}
	return nil
}

func (d *UserRepo) GetUserById(ctx context.Context, userID int64) (user *model.User, err error) {
	u, err := d.q.User.WithContext(ctx).Where(d.q.User.UserID.Eq(userID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (d *UserRepo) GetAssignRole(ctx context.Context, userID int64) (roleList []*model.Role, err error) {
	rels, err := d.q.UserRole.WithContext(ctx).Where(d.q.UserRole.UserID.Eq(userID)).Find()
	if err != nil {
		return nil, err
	}
	if len(rels) == 0 {
		return []*model.Role{}, nil
	}
	roleIDs := make([]int64, 0, len(rels))
	for _, r := range rels {
		roleIDs = append(roleIDs, r.RoleID)
	}
	roles, err := d.q.Role.WithContext(ctx).Where(d.q.Role.RoleID.In(roleIDs...)).Find()
	if err != nil {
		return nil, fmt.Errorf("查询用户角色失败(查询角色): %w", err)
	}
	return roles, nil
}

func (d *UserRepo) GetAssignMenuByUserId(ctx context.Context, userID int64) (menuList []*model.Menu, err error) {
	rels, err := d.q.UserRole.WithContext(ctx).Where(d.q.UserRole.UserID.Eq(userID)).Find()
	if err != nil {
		return nil, fmt.Errorf("查询用户菜单失败(查询用户角色关联): %w", err)
	}
	if len(rels) == 0 {
		return []*model.Menu{}, nil
	}
	roleIDs := make([]int64, 0, len(rels))
	for _, r := range rels {
		roleIDs = append(roleIDs, r.RoleID)
	}
	rmRels, err := d.q.RoleMenu.WithContext(ctx).Where(d.q.RoleMenu.RoleID.In(roleIDs...)).Find()
	if err != nil {
		return nil, fmt.Errorf("查询用户菜单失败(查询角色菜单关联): %w", err)
	}
	if len(rmRels) == 0 {
		return []*model.Menu{}, nil
	}
	menuIDSet := make(map[int64]struct{})
	menuIDs := make([]int64, 0)
	for _, rm := range rmRels {
		if _, ok := menuIDSet[rm.MenuID]; !ok {
			menuIDSet[rm.MenuID] = struct{}{}
			menuIDs = append(menuIDs, rm.MenuID)
		}
	}
	menus, err := d.q.Menu.WithContext(ctx).Where(d.q.Menu.MenuID.In(menuIDs...)).Find()
	if err != nil {
		return nil, fmt.Errorf("查询用户菜单失败(查询菜单): %w", err)
	}
	return menus, nil
}

func (d *UserRepo) GetUserList(ctx context.Context, username string, page, limit int64) (list []*model.User, total int64, err error) {
	q := d.q.User.WithContext(ctx)
	if username != "" {
		q = q.Where(d.q.User.Username.Like("%" + username + "%"))
	}
	total, err = q.Count()
	if err != nil {
		return nil, 0, err
	}
	if page > 0 && limit > 0 {
		q = q.Offset(int((page - 1) * limit)).Limit(int(limit))
	}
	list, err = q.Order(d.q.User.UserID.Desc()).Find()
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (d *UserRepo) UpdateUser(ctx context.Context, user *model.User) (err error) {
	updates := map[string]any{}
	// 仅更新调用方显式赋值的字段，nil/零值字段保持不变（部分更新语义）
	if user.Username != "" {
		updates["username"] = user.Username
	}
	if user.Name != nil {
		updates["name"] = user.Name
	}
	if user.Phone != nil {
		updates["phone"] = user.Phone
	}
	if user.Avatar != nil {
		updates["avatar"] = user.Avatar
	}
	// 仅当传入非空密码时才更新密码（并做 bcrypt 哈希），避免误清空密码
	if user.Password != "" {
		hashed, err := hashPassword(user.Password)
		if err != nil {
			return err
		}
		updates["password"] = hashed
	}
	if len(updates) == 0 {
		return nil
	}
	_, err = d.q.User.WithContext(ctx).
		Where(d.q.User.UserID.Eq(user.UserID)).
		Updates(updates)
	if err != nil {
		return fmt.Errorf("更新用户失败(更新用户): %w", err)
	}
	return nil
}

func (d *UserRepo) DeleteUser(ctx context.Context, userId int64) (err error) {
	_, err = d.q.User.WithContext(ctx).Where(d.q.User.UserID.Eq(userId)).Delete()
	return err
}

// LockUser 锁定/解锁用户（status=false 禁用，true 恢复）
func (d *UserRepo) LockUser(ctx context.Context, userID int64, status bool) (err error) {
	_, err = d.q.User.WithContext(ctx).
		Where(d.q.User.UserID.Eq(userID)).
		Update(d.q.User.Status, status)
	return err
}

func (d *UserRepo) DeleteAssignRoleByUserId(ctx context.Context, userID int64) (err error) {
	_, err = d.q.UserRole.WithContext(ctx).Where(d.q.UserRole.UserID.Eq(userID)).Delete()
	return err
}

// DoAssign 为用户重新分配角色：先删除原角色关联，再批量插入新关联。
// 使用事务保证原子性，两段写任一失败整体回滚，避免「原关联已删、新关联未插入」导致用户角色丢失。
func (d *UserRepo) DoAssign(ctx context.Context, userId string, roleIds []string) (err error) {
	uid := idconv.ToInt64Safe(userId)
	return d.q.Transaction(func(tx *query.Query) error {
		_, err := tx.UserRole.WithContext(ctx).Where(tx.UserRole.UserID.Eq(uid)).Delete()
		if err != nil {
			return fmt.Errorf("分配角色失败(删除原角色关联): %w", err)
		}
		if len(roleIds) == 0 {
			return nil
		}
		list := make([]*model.UserRole, 0, len(roleIds))
		for _, rid := range roleIds {
			list = append(list, &model.UserRole{
				UserID: uid,
				RoleID: idconv.ToInt64Safe(rid),
			})
		}
		if err := tx.UserRole.WithContext(ctx).Create(list...); err != nil {
			return fmt.Errorf("分配角色失败(批量插入角色关联): %w", err)
		}
		return nil
	})
}
