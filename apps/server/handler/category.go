package handler

import (
	"context"
	"log/slog"

	"github.com/liwook/go-vue-selection/pkg/errs"
	"github.com/liwook/go-vue-selection/pkg/idconv"
	"github.com/liwook/go-vue-selection/types"

	"github.com/gin-gonic/gin"

	"github.com/liwook/go-vue-selection/pkg/result"
)

// CategoryService 在使用方（controller）定义接口
type CategoryService interface {
	GetCategory1(context.Context) ([]*types.Category1, error)
	GetCategory2(context.Context, int64) ([]*types.Category2, error)
	GetCategory3(context.Context, int64) ([]*types.Category3, error)
	CreateCategory2(context.Context, *types.ParamC2Create) (*types.Category2, error)
	CreateCategory3(context.Context, *types.ParamC3Create) (*types.Category3, error)
}

type categoryHandler struct {
	categorySvc CategoryService
}

func NewCategoryHandler(svc CategoryService) *categoryHandler {
	return &categoryHandler{categorySvc: svc}
}

// RegisterRoutes 注册商品分类相关路由（需在 JWT 中间件之后调用）
func (c *categoryHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/category1", c.GetCategory1)
	r.GET("/category2/:category1Id", c.GetCategory2)
	r.GET("/category3/:category2Id", c.GetCategory3)
	r.POST("/category2", c.CreateCategory2)
	r.POST("/category3", c.CreateCategory3)
}

// GetCategory1 获取一级分类
// @Summary 获取一级分类接口
// @Description 获取一级分类
// @Tags 商品分类
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=[]types.Category1}
// @Router /api/v1/product/category1 [get]
func (c *categoryHandler) GetCategory1(ctx *gin.Context) {
	data, err := c.categorySvc.GetCategory1(ctx.Request.Context())
	if err != nil {
		result.Error(ctx, errs.CodeOf(err))
		return
	}
	// 返回响应
	result.Success(ctx, data)
}

// GetCategory2 获取二级分类
// @Summary 获取二级分类接口
// @Description 获取二级分类
// @Tags 商品分类
// @Accept application/json
// @Produce application/json
// @Param category1Id path string true "分类一 ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=[]types.Category2}
// @Router /api/v1/product/category2/{category1Id} [get]
func (c *categoryHandler) GetCategory2(ctx *gin.Context) {
	idStr := ctx.Param("category1Id")
	category1Id := idconv.ToInt64Safe(idStr)
	data, err := c.categorySvc.GetCategory2(ctx.Request.Context(), category1Id)
	if err != nil {
		result.Error(ctx, errs.CodeOf(err))
		return
	}
	// 返回响应
	result.Success(ctx, data)
}

// GetCategory3 获取三级分类
// @Summary 获取三级分类接口
// @Description 获取三级分类
// @Tags 商品分类
// @Accept application/json
// @Produce application/json
// @Param category2Id path string true "分类二 ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=[]types.Category3}
// @Router /api/v1/product/category3/{category2Id} [get]
func (c *categoryHandler) GetCategory3(ctx *gin.Context) {
	idStr := ctx.Param("category2Id")
	category2Id := idconv.ToInt64Safe(idStr)
	data, err := c.categorySvc.GetCategory3(ctx.Request.Context(), category2Id)
	if err != nil {
		result.Error(ctx, errs.CodeOf(err))
		return
	}
	// 返回响应
	result.Success(ctx, data)
}

// CreateCategory2 新增二级分类
// @Summary 新增二级分类
// @Description 新增二级分类，由后端生成二级分类 ID，返回新建的分类信息
// @Tags 商品分类
// @Accept application/json
// @Produce application/json
// @Param object body types.ParamC2Create true "二级分类信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=types.Category2}
// @Router /api/v1/product/category2 [post]
func (c *categoryHandler) CreateCategory2(ctx *gin.Context) {
	p := new(types.ParamC2Create)
	if err := ctx.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		slog.Error("CreateCategory2 with invalid param", slog.Any("error", err))
		result.Error(ctx, result.CodeInvalidParam)
		return
	}

	// 2. 业务处理
	data, err := c.categorySvc.CreateCategory2(ctx.Request.Context(), p)
	if err != nil {
		result.Error(ctx, errs.CodeOf(err))
		return
	}

	// 3. 返回结果（含后端生成的二级分类 ID）
	result.Success(ctx, data)
}

// CreateCategory3 新增三级分类
// @Summary 新增三级分类
// @Description 新增三级分类，由后端生成三级分类 ID，返回新建的分类信息
// @Tags 商品分类
// @Accept application/json
// @Produce application/json
// @Param object body types.ParamC3Create true "三级分类信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=types.Category3}
// @Router /api/v1/product/category3 [post]
func (c *categoryHandler) CreateCategory3(ctx *gin.Context) {
	p := new(types.ParamC3Create)
	if err := ctx.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		slog.Error("CreateCategory3 with invalid param", slog.Any("error", err))
		result.Error(ctx, result.CodeInvalidParam)
		return
	}

	// 2. 业务处理
	data, err := c.categorySvc.CreateCategory3(ctx.Request.Context(), p)
	if err != nil {
		result.Error(ctx, errs.CodeOf(err))
		return
	}

	// 3. 返回结果（含后端生成的三级分类 ID）
	result.Success(ctx, data)
}
