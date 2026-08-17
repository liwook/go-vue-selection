package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/pkg/errs"
	"github.com/liwook/go-vue-selection/pkg/idconv"
	"github.com/liwook/go-vue-selection/pkg/result"
	"github.com/liwook/go-vue-selection/pkg/snowflake"
	"github.com/liwook/go-vue-selection/repository"
	"github.com/liwook/go-vue-selection/types"
)

type roleService struct {
	repo *repository.RoleRepo
}

func NewRoleService(repo *repository.RoleRepo) *roleService {
	return &roleService{repo: repo}
}

func (r *roleService) GetRoleList(ctx context.Context, roleName string, page, limit int64) (data *types.ResponseRoleList, err error) {
	roles, total, err := r.repo.GetRoleList(ctx, roleName, page, limit)
	if err != nil {
		slog.Error("查询角色列表失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeRoleDBErr, err)
	}
	records := make([]*types.Role, 0, len(roles))
	for _, role := range roles {
		records = append(records, modelToRole(role))
	}
	data = &types.ResponseRoleList{
		Records:     records,
		Total:       total,
		Size:        limit,
		Current:     page,
		Pages:       pageCount(total, limit),
		SearchCount: true,
	}
	return
}

func (r *roleService) SaveRole(ctx context.Context, p *types.ParamRoleSave) (err error) {
	// 1. 判断角色是否已经存在
	if err = r.repo.CheckRoleExist(ctx, p.RoleName); err != nil {
		if errors.Is(err, repository.ErrorRoleExist) {
			return errs.Wrap(result.CodeRoleExist, err)
		}
		return errs.Wrap(result.CodeRoleDBErr, err)
	}

	roleId := snowflake.GenID()
	// 2. 生成一个角色实例
	role := &model.Role{
		RoleID:   roleId, // 原生 int64
		RoleName: p.RoleName,
		Remark:   &p.Remark,
	}

	// 3. 保存数据库
	err = r.repo.InsertRole(ctx, role)
	if err != nil {
		slog.Error("创建角色失败", slog.Any("error", err))
		return errs.Wrap(result.CodeRoleDBErr, err)
	}

	return err
}

func (r *roleService) UpdateRole(ctx context.Context, p *types.ParamRoleUpdate) (err error) {
	// 构造一个角色实例
	role := &model.Role{
		RoleID:   idconv.ToInt64Safe(p.RoleID),
		RoleName: p.RoleName,
		Remark:   &p.Remark,
	}

	// 保存进数据库
	err = r.repo.UpdateRole(ctx, role)
	if err != nil {
		slog.Error("更新角色失败", slog.Any("error", err))
		return errs.Wrap(result.CodeRoleDBErr, err)
	}

	return err
}

func (r *roleService) DeleteRole(ctx context.Context, roleId int64) (err error) {
	if err = r.repo.DeleteRole(ctx, roleId); err != nil {
		slog.Error("删除角色失败", slog.Int64("role_id", roleId), slog.Any("error", err))
		if errs.IsForeignKeyViolation(err) {
			return errs.Wrap(result.CodeRoleInUse, err)
		}
		return errs.Wrap(result.CodeRoleDBErr, err)
	}
	return nil
}
