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

// AttrService 在使用方（controller）定义接口
type AttrService interface {
	UpdateAttr(context.Context, *types.Attr) error
	CreateAttr(context.Context, *types.Attr) error
	GetAttr(context.Context, int64, int64, int64) ([]*types.Attr, error)
	DeleteAttr(context.Context, int64) error
}

type attrHandler struct {
	attrSvc AttrService
}

func NewAttrHandler(svc AttrService) *attrHandler {
	return &attrHandler{attrSvc: svc}
}

// RegisterRoutes 注册商品属性相关路由（需在 JWT 中间件之后调用）
func (a *attrHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/attr", a.SaveAttrInfo)
	r.GET("/attr/:category1Id/:category2Id/:category3Id", a.GetAttr)
	r.DELETE("/attr/:attrId", a.DeleteAttr)
}

// SaveAttrInfo 添加或者修改已有的属性
// @Summary 添加或者修改已有的属性的接口
// @Description 处理添加或者修改已有的属性请求
// @Tags 商品属性
// @Accept application/json
// @Produce application/json
// @Param object body types.Attr true "属性信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /admin/product/attr [post]
func (a *attrHandler) SaveAttrInfo(c *gin.Context) {
	p := new(types.Attr)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		slog.Error("SaveAttrInfo with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}

	// 2. 业务处理
	// 有属性ID则进行更新
	if p.AttrID != "" && p.AttrID != "0" {
		err := a.attrSvc.UpdateAttr(c.Request.Context(), p)
		if err != nil {
			result.Error(c, errs.CodeOf(err))
			return
		}
	} else {
		// 没有属性ID则进行创建
		err := a.attrSvc.CreateAttr(c.Request.Context(), p)
		if err != nil {
			result.Error(c, errs.CodeOf(err))
			return
		}
	}

	// 3. 返回结果
	result.Success(c, nil)
}

// GetAttr 获取分类下已有的属性与属性值
// @Summary 获取分类下已有的属性与属性值接口
// @Description 处理获取分类下已有的属性与属性值请求
// @Tags 商品属性
// @Accept application/json
// @Produce application/json
// @Param category1Id path string true "分类一 ID"
// @Param category2Id path string true "分类二 ID"
// @Param category3Id path string true "分类三 ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=[]types.Attr}
// @Router /admin/product/attr/{category1Id}/{category2Id}/{category3Id} [get]
func (a *attrHandler) GetAttr(c *gin.Context) {
	c1Id, c2Id, c3Id, err := getAllCategoryID(c)
	if err != nil {
		slog.Error("GetAttr with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}
	// 获取数据
	data, err := a.attrSvc.GetAttr(c.Request.Context(), c1Id, c2Id, c3Id)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	// 返回响应
	result.Success(c, data)
}

// DeleteAttr 删除基础属性
// @Summary 删除基础属性接口
// @Description 处理删除基础属性请求
// @Tags 商品属性
// @Accept application/json
// @Produce application/json
// @Param attrId path string true "属性 ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /admin/product/attr/{attrId} [delete]
func (a *attrHandler) DeleteAttr(c *gin.Context) {
	idStr := c.Param("attrId")
	attrId := idconv.ToInt64Safe(idStr)
	err := a.attrSvc.DeleteAttr(c.Request.Context(), attrId)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}

	result.Success(c, nil)
}
