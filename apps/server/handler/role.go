package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/liwook/go-vue-selection/pkg/errs"
	"github.com/liwook/go-vue-selection/pkg/idconv"
	"github.com/liwook/go-vue-selection/types"
	"log/slog"

	"github.com/liwook/go-vue-selection/pkg/result"
)

// RoleService 在使用方（controller）定义接口
type RoleService interface {
	GetRoleList(context.Context, string, int64, int64) (*types.ResponseRoleList, error)
	SaveRole(context.Context, *types.ParamRoleSave) error
	UpdateRole(context.Context, *types.ParamRoleUpdate) error
	DeleteRole(context.Context, int64) error
}

type roleHandler struct {
	roleSvc RoleService
}

func NewRoleHandler(svc RoleService) *roleHandler {
	return &roleHandler{roleSvc: svc}
}

// RegisterRoutes 注册角色管理相关路由（需在 JWT 中间件之后调用）
func (r *roleHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/role", r.GetRole)
	rg.POST("/role", r.SaveRole)
	rg.PUT("/role/:roleId", r.UpdateRole)
	rg.DELETE("/role/:roleId", r.DeleteRole)
}

// GetRole 获取角色列表
// @Summary 获取角色分页列表
// @Description 获取角色列表
// @Tags 角色
// @Accept application/json
// @Produce application/json
// @Param page query int true "当前页码"
// @Param limit query int true "每页记录数"
// @Param roleName query string false "角色名"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=types.ResponseRoleList}
// @Router /api/v1/acl/role [get]
func (r *roleHandler) GetRole(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	roleName := c.Query("roleName")

	// 获取数据
	data, err := r.roleSvc.GetRoleList(c.Request.Context(), roleName, page, size)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}

	// 返回数据
	result.Success(c, data)
}

// SaveRole 新增角色
// @Summary 新增角色接口
// @Description 处理新增角色请求
// @Tags 角色
// @Accept application/json
// @Produce application/json
// @Param object body types.ParamRoleSave true "角色信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/acl/role [post]
func (r *roleHandler) SaveRole(c *gin.Context) {
	p := new(types.ParamRoleSave)
	if err := c.ShouldBindJSON(p); err != nil {
		slog.Error("SaveRole with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}

	err := r.roleSvc.SaveRole(c.Request.Context(), p)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}

	result.Success(c, nil)
}

// UpdateRole 更新角色
// @Summary 更新角色接口
// @Description 处理更新角色请求
// @Tags 角色
// @Accept application/json
// @Produce application/json
// @Param roleId path string true "角色 ID"
// @Param object body types.ParamRoleUpdate true "角色信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/acl/role/{roleId} [put]
func (r *roleHandler) UpdateRole(c *gin.Context) {
	p := new(types.ParamRoleUpdate)
	if err := c.ShouldBindJSON(p); err != nil {
		slog.Error("UpdateRole with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}
	roleId := idconv.ToInt64Safe(c.Param("roleId"))
	p.RoleID = idconv.ToStr(roleId)
	err := r.roleSvc.UpdateRole(c.Request.Context(), p)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)
}

// DeleteRole 删除角色
// @Summary 删除角色接口
// @Description 处理删除角色请求
// @Tags 角色
// @Accept application/json
// @Produce application/json
// @Param roleId path string true "角色 ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/acl/role/{roleId} [delete]
func (r *roleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("roleId")
	roleId := idconv.ToInt64Safe(idStr)
	err := r.roleSvc.DeleteRole(c.Request.Context(), roleId)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)
}
