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

type menuService struct {
	repo *repository.MenuRepo
}

func NewMenuService(repo *repository.MenuRepo) *menuService {
	return &menuService{repo: repo}
}

func (m *menuService) GetMenu(ctx context.Context) (data []types.Menu, err error) {
	// 1. 查询所有菜单
	modelList, err := m.repo.GetMenuList(ctx)
	if err != nil {
		slog.Error("查询菜单列表失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeMenuDBErr, err)
	}
	if len(modelList) == 0 {
		return []types.Menu{}, nil
	}
	menuList := modelToMenuList(modelList)
	menuValueList := make([]types.Menu, 0, len(menuList))
	for _, m := range menuList {
		menuValueList = append(menuValueList, *m)
	}

	// 2. 格式化成返回需要的树形格式
	data, err = buildMenuTree(menuValueList)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (m *menuService) SaveMenu(ctx context.Context, p *types.ParamMenuSave) (err error) {
	// 1. 判断菜单是否已经存在
	if err = m.repo.CheckMenuExist(ctx, p.Name); err != nil {
		if errors.Is(err, repository.ErrorMenuExist) {
			return errs.Wrap(result.CodeMenuExist, err)
		}
		slog.Error("校验菜单是否存在失败", slog.String("name", p.Name), slog.Any("error", err))
		return errs.Wrap(result.CodeMenuDBErr, err)
	}

	menuId := snowflake.GenID()
	pid := idconv.ToInt64Safe(p.ParentID)
	// 2. 生成一个菜单实例
	menu := &model.Menu{
		MenuID: menuId, // 原生 int64
		Name:   p.Name,
		Pid:    &pid,
		Code:   &p.CODE,
		ToCode: nil,
		Type:   int32(p.TYPE),
		Status: true,
		Level:  int32(p.LEVEL),
	}
	// 3. 保存进数据库
	err = m.repo.InsertMenu(ctx, menu)
	if err != nil {
		slog.Error("创建菜单失败", slog.String("name", p.Name), slog.Any("error", err))
		return errs.Wrap(result.CodeMenuDBErr, err)
	}

	return err
}

func (m *menuService) UpdateMenu(ctx context.Context, p *types.ParamMenuUpdate) (err error) {
	pid := idconv.ToInt64Safe(p.ParentID)
	// 构造一个菜单实例
	menu := &model.Menu{
		MenuID: idconv.ToInt64Safe(p.MenuID),
		Name:   p.Name,
		Pid:    &pid,
		Code:   &p.CODE,
		Level:  int32(p.LEVEL),
	}

	// 保存进数据库
	err = m.repo.UpdateMenu(ctx, menu)
	if err != nil {
		slog.Error("更新菜单失败", slog.String("menu_id", p.MenuID), slog.Any("error", err))
		return errs.Wrap(result.CodeMenuDBErr, err)
	}

	return err
}

func (m *menuService) DeleteMenu(ctx context.Context, menuId int64) (err error) {
	// 根据当前菜单ID，查询是否包含子菜单
	count, err := m.repo.QuerySubMenuByID(ctx, menuId)
	if err != nil {
		slog.Error("查询子菜单失败", slog.Int64("menu_id", menuId), slog.Any("error", err))
		return errs.Wrap(result.CodeMenuDBErr, err)
	}
	if count > 0 {
		return errs.Wrap(result.CodeMenuNodeExist, errors.New("节点下有子节点，不可以删除"))
	}

	// 根据菜单ID删除菜单
	err = m.repo.DeleteMenuByID(ctx, menuId)
	if err != nil {
		slog.Error("删除菜单失败", slog.Int64("menu_id", menuId), slog.Any("error", err))
		if errs.IsForeignKeyViolation(err) {
			return errs.Wrap(result.CodeMenuInUse, err)
		}
		return errs.Wrap(result.CodeMenuDBErr, err)
	}

	return err
}

func (m *menuService) ToAssign(ctx context.Context, roleId int64) (data []types.Menu, err error) {
	// 1. 查询所有菜单
	modelList, err := m.repo.GetMenuList(ctx)
	if err != nil {
		slog.Error("查询菜单列表失败", slog.Any("error", err))
		return nil, errs.Wrap(result.CodeMenuDBErr, err)
	}
	if len(modelList) == 0 {
		return []types.Menu{}, nil
	}
	menuList := modelToMenuList(modelList)
	menuValueList := make([]types.Menu, 0, len(menuList))
	for _, m := range menuList {
		menuValueList = append(menuValueList, *m)
	}

	// 2. 查询已经具有的权限
	assignMenuIdList, err := m.repo.GetAssignMenu(ctx, roleId)
	if err != nil {
		slog.Error("查询角色菜单失败", slog.Int64("role_id", roleId), slog.Any("error", err))
		return nil, errs.Wrap(result.CodeMenuDBErr, err)
	}

	// 3. 将所有权限中已经具有的权限 select 值置为 true
	if len(assignMenuIdList) != 0 {
		for i := range menuValueList {
			for _, menuId := range assignMenuIdList {
				if menuValueList[i].MenuID == menuId {
					menuValueList[i].SELECT = true
				}
			}
		}
	}

	// 3. 格式化成返回需要的树形格式
	data, err = buildMenuTree(menuValueList)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (m *menuService) DoAssign(ctx context.Context, roleId int64, menuIds []int64) (err error) {
	// repository.DoAssign 内部已通过事务完成「删除原菜单关联 + 重新分配」，
	// 两步写整体原子，无需在此额外删除。
	err = m.repo.DoAssign(ctx, roleId, menuIds)
	if err != nil {
		slog.Error("分配角色菜单失败", slog.Int64("role_id", roleId), slog.Any("error", err))
		return errs.Wrap(result.CodeMenuDBErr, err)
	}

	return err
}

// buildMenuTree 将扁平菜单列表封装成树形结构（以 ParentID 为 "0" 的节点为根）。
func buildMenuTree(menuList []types.Menu) (tree []types.Menu, err error) {
	for _, menu := range menuList {
		if menu.ParentID == "0" {
			tree = append(tree, findMenuChildren(menu, menuList))
		}
	}
	return tree, nil
}

// findMenuChildren 递归查找 menu 的子节点并挂到 menu.CHILDREN 上，返回完整的父节点（带子树）。
func findMenuChildren(menu types.Menu, menuList []types.Menu) (parent types.Menu) {
	for _, child := range menuList {
		if menu.MenuID == child.ParentID {
			menu.CHILDREN = append(menu.CHILDREN, findMenuChildren(child, menuList))
		}
	}
	return menu
}
