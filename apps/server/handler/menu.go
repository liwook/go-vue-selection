package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/liwook/go-vue-selection/pkg/errs"
	"github.com/liwook/go-vue-selection/pkg/idconv"
	"github.com/liwook/go-vue-selection/types"
	"log/slog"
	"strings"

	"github.com/liwook/go-vue-selection/pkg/result"
)

// MenuService 在使用方（controller）定义接口
type MenuService interface {
	GetMenu(context.Context) ([]types.Menu, error)
	SaveMenu(context.Context, *types.ParamMenuSave) error
	UpdateMenu(context.Context, *types.ParamMenuUpdate) error
	DeleteMenu(context.Context, int64) error
	ToAssign(context.Context, int64) ([]types.Menu, error)
	DoAssign(context.Context, int64, []int64) error
}

type menuHandler struct {
	menuSvc MenuService
}

func NewMenuHandler(svc MenuService) *menuHandler {
	return &menuHandler{menuSvc: svc}
}

// RegisterRoutes 注册菜单管理相关路由（需在 JWT 中间件之后调用）
func (m *menuHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/permission", m.GetMenu)
	r.POST("/permission", m.SaveMenu)
	r.PUT("/permission/:permissionId", m.UpdateMenu)
	r.DELETE("/permission/:permissionId", m.DeleteMenu)
	r.GET("/permission/role/:roleId", m.ToAssign)
	r.POST("/permission/role/:roleId", m.DoAssign)
}

// GetMenu 获取菜单列表
// @Summary 获取菜单列表
// @Description 获取菜单列表
// @Tags 菜单
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=[]types.Menu}
// @Router /api/v1/acl/permission [get]
// @Failure 500 {object} result.ResponseData
func (m *menuHandler) GetMenu(c *gin.Context) {
	data, err := m.menuSvc.GetMenu(c.Request.Context())
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}

	result.Success(c, data)
}

// SaveMenu 新增菜单
// @Summary 新增菜单接口
// @Description 处理新增菜单请求
// @Tags 菜单
// @Accept application/json
// @Produce application/json
// @Param object body types.ParamMenuSave true "菜单信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/acl/permission [post]
// @Failure 500 {object} result.ResponseData
func (m *menuHandler) SaveMenu(c *gin.Context) {
	p := new(types.ParamMenuSave)
	if err := c.ShouldBindJSON(p); err != nil {
		slog.Error("SaveMenu with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}
	err := m.menuSvc.SaveMenu(c.Request.Context(), p)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)
}

// UpdateMenu 更新菜单
// @Summary 更新菜单接口
// @Description 处理更新菜单请求
// @Tags 菜单
// @Accept application/json
// @Produce application/json
// @Param permissionId path string true "菜单 ID"
// @Param object body types.ParamMenuUpdate true "菜单信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/acl/permission/{permissionId} [put]
// @Failure 500 {object} result.ResponseData
func (m *menuHandler) UpdateMenu(c *gin.Context) {
	p := new(types.ParamMenuUpdate)
	if err := c.ShouldBindJSON(p); err != nil {
		slog.Error("UpdateMenu with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}
	permissionId := idconv.ToInt64Safe(c.Param("permissionId"))
	p.MenuID = idconv.ToStr(permissionId)
	err := m.menuSvc.UpdateMenu(c.Request.Context(), p)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)
}

// DeleteMenu 删除菜单
// @Summary 删除菜单接口
// @Description 处理删除菜单请求
// @Tags 菜单
// @Accept application/json
// @Produce application/json
// @Param permissionId path string true "菜单 ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/acl/permission/{permissionId} [delete]
// @Failure 500 {object} result.ResponseData
func (m *menuHandler) DeleteMenu(c *gin.Context) {
	idStr := c.Param("permissionId")
	menuId := idconv.ToInt64Safe(idStr)
	err := m.menuSvc.DeleteMenu(c.Request.Context(), menuId)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)

}

// ToAssign 根据角色获取菜单
// @Summary 根据角色获取菜单接口
// @Description 根据角色获取菜单请求
// @Tags 菜单
// @Accept application/json
// @Produce application/json
// @Param roleId path string true "角色ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=[]types.Menu}
// @Router /api/v1/acl/permission/role/{roleId} [get]
// @Failure 500 {object} result.ResponseData
func (m *menuHandler) ToAssign(c *gin.Context) {
	idStr := c.Param("roleId")
	roleId := idconv.ToInt64Safe(idStr)
	data, err := m.menuSvc.ToAssign(c.Request.Context(), roleId)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, data)
}

// DoAssign 给角色分配权限
// @Summary 给角色分配权限
// @Description 给角色分配权限
// @Tags 菜单
// @Accept application/json
// @Produce application/json
// @Param roleId path string true "角色 ID"
// @Param permissionId query string true "菜单 ID列表（逗号分隔）"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/acl/permission/role/{roleId} [post]
// @Failure 500 {object} result.ResponseData
func (m *menuHandler) DoAssign(c *gin.Context) {
	idStr := c.Param("roleId")
	roleId := idconv.ToInt64Safe(idStr)

	permissionIdStr := c.Query("permissionId")
	strArr := strings.Split(permissionIdStr, ",")
	var intArr []int64
	for _, s := range strArr {
		num := idconv.ToInt64Safe(s)
		if num == 0 {
			continue
		}
		intArr = append(intArr, num)
	}

	err := m.menuSvc.DoAssign(c.Request.Context(), roleId, intArr)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}

	result.Success(c, nil)
}
