package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/pkg/errs"
	"github.com/liwook/go-vue-selection/pkg/idconv"
	"github.com/liwook/go-vue-selection/pkg/jwt"
	"github.com/liwook/go-vue-selection/pkg/result"
	"github.com/liwook/go-vue-selection/pkg/snowflake"
	"github.com/liwook/go-vue-selection/repository"
	"github.com/liwook/go-vue-selection/types"
)

type userService struct {
	userRepo *repository.UserRepo
	roleRepo *repository.RoleRepo
	jwt      *jwt.JWT
}

func NewUserService(userRepo *repository.UserRepo, roleRepo *repository.RoleRepo, jwt *jwt.JWT) *userService {
	return &userService{userRepo: userRepo, roleRepo: roleRepo, jwt: jwt}
}

func (u *userService) Login(ctx context.Context, p *types.ParamUserLogin) (token string, err error) {
	// repository.Login 返回纯 model.User（已完成密码校验与账号状态校验）
	user, err := u.userRepo.Login(ctx, p.Username, p.Password)
	if err != nil {
		if errors.Is(err, repository.ErrorUserNotExist) || errors.Is(err, repository.ErrorPasswordWrong) {
			return "", errs.Wrap(result.CodeInvalidPassword, err)
		}
		slog.Error("用户登录失败", slog.String("username", p.Username), slog.Any("error", err))
		return "", errs.Wrap(result.CodeUserDBErr, err)
	}
	// 生成 JWT
	token, err = u.jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		slog.Error("生成用户令牌失败", slog.String("username", p.Username), slog.Any("error", err))
		return "", errs.Wrap(result.CodeUserDBErr, err)
	}
	return token, err
}

func (u *userService) SignUp(ctx context.Context, p *types.ParamUserSignUp) (err error) {
	// 1. 判断用户是否存在
	if err = u.userRepo.CheckUserExist(ctx, p.Username); err != nil {
		if errors.Is(err, repository.ErrorUserExist) {
			return errs.Wrap(result.CodeUserExist, err)
		}
		slog.Error("校验用户是否存在失败", slog.String("username", p.Username), slog.Any("error", err))
		return errs.Wrap(result.CodeUserDBErr, err)
	}

	// 2. 生成UID
	userID := snowflake.GenID()

	// 构造一个User实例
	user := &model.User{
		UserID:   userID, // 原生 int64
		Username: p.Username,
		Name:     &p.Name,
		Password: p.Password, // 密码哈希在 repo 层完成
	}

	// 3. 保存进数据库
	err = u.userRepo.InsertUser(ctx, user)
	if err != nil {
		slog.Error("创建用户失败", slog.String("username", p.Username), slog.Any("error", err))
		return errs.Wrap(result.CodeUserDBErr, err)
	}

	return err
}

func (u *userService) GetUserInfo(ctx context.Context, userID int64) (*types.ResponseUserInfo, error) {
	user, err := u.userRepo.GetUserById(ctx, userID)
	if err != nil {
		slog.Error("查询用户信息失败",
			slog.Int64("user_id", userID),
			slog.Any("error", err))
		return nil, errs.Wrap(result.CodeUserDBErr, err)
	}
	// 获取用户角色
	modelRoles, err := u.userRepo.GetAssignRole(ctx, userID)
	if err != nil {
		slog.Error("查询用户角色失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeUserDBErr, err)
	}
	roleNameList := make([]string, 0, len(modelRoles))
	for _, role := range modelRoles {
		roleNameList = append(roleNameList, role.RoleName)
	}

	// 获取当前用户拥有的所有菜单
	modelMenus, err := u.userRepo.GetAssignMenuByUserId(ctx, userID)
	if err != nil {
		slog.Error("查询用户菜单失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeUserDBErr, err)
	}
	menuList := modelToMenuList(modelMenus)
	// 获取用户路由 和 用户按钮权限
	userRoutes := make([]string, 0, 12)
	userButtons := make([]string, 0, 12)
	for _, menu := range menuList {
		if menu.LEVEL != 4 {
			userRoutes = append(userRoutes, menu.CODE)
		} else {
			userButtons = append(userButtons, menu.CODE)
		}
	}

	userInfo := &types.ResponseUserInfo{
		Routes:  userRoutes,
		Buttons: userButtons,
		Roles:   roleNameList,
		Name:    derefStr(user.Name),
		Avatar:  derefStr(user.Avatar),
	}
	return userInfo, err

}

func (u *userService) GetUserList(ctx context.Context, username string, page, limit int64) (data *types.ResponseUserList, err error) {
	list, total, err := u.userRepo.GetUserList(ctx, username, page, limit)
	if err != nil {
		slog.Error("查询用户列表失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeUserDBErr, err)
	}
	records := make([]*types.ResponseUser, 0, len(list))
	for _, mu := range list {
		records = append(records, modelToResponseUser(mu))
	}
	data = &types.ResponseUserList{
		Records: records,
		Total:   total,
		Size:    limit,
		Current: page,
		Pages:   pageCount(total, limit),
	}
	return
}

func (u *userService) UpdateUser(ctx context.Context, userID int64, p *types.ParamUserUpdate) (err error) {
	user := &model.User{
		UserID: userID,
	}
	// 仅当字段非 nil 时才赋值，实现“只更新前端显式发送的非空字段”
	// 注意：Username（登录名）不允许用户自行修改，此处不处理
	if p.Name != nil {
		user.Name = p.Name
	}
	if p.Avatar != nil {
		user.Avatar = p.Avatar
	}
	err = u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		slog.Error("更新用户失败", slog.Int64("user_id", userID), slog.Any("error", err))
		return errs.Wrap(result.CodeUserDBErr, err)
	}
	return err

}

// adminRoleName 是具备用户管理权限的管理员角色名（数据源：init-sql 中 role_id=1）
const adminRoleName = "管理员"

// IsAdmin 判断指定用户是否为管理员（实时查库，保证权限数据最新）
func (u *userService) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	roles, err := u.userRepo.GetAssignRole(ctx, userID)
	if err != nil {
		slog.Error("查询用户角色失败", slog.Int64("user_id", userID), slog.Any("error", err))
		return false, errs.Wrap(result.CodeUserDBErr, err)
	}
	for _, r := range roles {
		if r.RoleName == adminRoleName {
			return true, nil
		}
	}
	return false, nil
}

// DeleteUserByAdmin 管理员删除指定用户（带权限校验，operatorID 为操作者）
func (u *userService) DeleteUserByAdmin(ctx context.Context, operatorID, targetUserID int64) (err error) {
	ok, err := u.IsAdmin(ctx, operatorID)
	if err != nil {
		return err
	}
	if !ok {
		return errs.Wrap(result.CodeNoPermission, errors.New("无删除用户权限"))
	}
	if err = u.userRepo.DeleteUser(ctx, targetUserID); err != nil {
		slog.Error("删除用户失败", slog.Int64("user_id", targetUserID), slog.Any("error", err))
		return errs.Wrap(result.CodeUserDBErr, err)
	}
	return nil
}

// DeleteSelfAccount 注销当前登录用户自己（直接删除，不经过管理员权限校验）
func (u *userService) DeleteSelfAccount(ctx context.Context, userID int64) (err error) {
	if err = u.userRepo.DeleteUser(ctx, userID); err != nil {
		slog.Error("注销用户失败", slog.Int64("user_id", userID), slog.Any("error", err))
		return errs.Wrap(result.CodeUserDBErr, err)
	}
	return nil
}

func (u *userService) ToAssign(ctx context.Context, userId int64) (data *types.ResponseToAssignRole, err error) {
	modelAllRoles, err := u.roleRepo.GetAllRoleList(ctx)
	if err != nil {
		slog.Error("查询全部角色失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeUserDBErr, err)
	}
	allRolesList := make([]*types.Role, 0, len(modelAllRoles))
	for _, r := range modelAllRoles {
		allRolesList = append(allRolesList, modelToRole(r))
	}
	modelAssignRoles, err := u.userRepo.GetAssignRole(ctx, userId)
	if err != nil {
		slog.Error("查询用户角色失败", slog.Int64("user_id", userId), slog.Any("error", err))
		return nil, errs.Wrap(result.CodeUserDBErr, err)
	}
	assignRoles := make([]*types.Role, 0, len(modelAssignRoles))
	for _, r := range modelAssignRoles {
		assignRoles = append(assignRoles, modelToRole(r))
	}
	toAssign := &types.ResponseToAssignRole{
		AssignRoles:  assignRoles,
		AllRolesList: allRolesList,
	}
	return toAssign, err
}

func (u *userService) DoAssignRole(ctx context.Context, p *types.ParamDoAssignRole) (err error) {
	// repository.DoAssign 内部已通过事务完成「删除原角色关联 + 重新分配」，
	// 两步写整体原子，无需在此额外删除。
	err = u.userRepo.DoAssign(ctx, p.UserID, p.RoleIDList)
	if err != nil {
		slog.Error("分配用户角色失败", slog.String("user_id", p.UserID), slog.Any("error", err))
		return errs.Wrap(result.CodeUserDBErr, err)
	}
	return err
}

// LockUser 锁定/解锁用户
func (u *userService) LockUser(ctx context.Context, p *types.ParamUserLock) (err error) {
	if err = u.userRepo.LockUser(ctx, idconv.ToInt64Safe(p.UserID), p.Status); err != nil {
		slog.Error("锁定/解锁用户失败",
			slog.String("user_id", p.UserID),
			slog.Bool("status", p.Status),
			slog.Any("error", err))
		return errs.Wrap(result.CodeUserDBErr, err)
	}
	return err
}
