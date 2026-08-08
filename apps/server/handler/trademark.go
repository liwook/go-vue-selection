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

// TrademarkService 在使用方（controller）定义接口
type TrademarkService interface {
	CreateTrademark(context.Context, *types.ParamTmSave) error
	GetTrademarkList(context.Context, int64, int64) (*types.ResponseTmList, error)
	UpdateTrademark(context.Context, *types.ParamTmUpdate) error
	DeleteTrademark(context.Context, int64) error
	GetAllTrademarkList(context.Context) ([]types.Trademark, error)
}

type trademarkHandler struct {
	tmSvc TrademarkService
}

func NewTrademarkHandler(svc TrademarkService) *trademarkHandler {
	return &trademarkHandler{tmSvc: svc}
}

// RegisterRoutes 注册品牌相关路由（需在 JWT 中间件之后调用）
func (t *trademarkHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/trademark", t.CreateTrademark)
	r.GET("/trademark", t.GetTrademark)
	r.PUT("/trademark/:trademarkId", t.UpdateTrademark)
	r.DELETE("/trademark/:trademarkId", t.DeleteTrademark)
	r.GET("/trademark/all", t.GetAllTrademarkList)
}

// CreateTrademark 新增品牌接口
// @Summary 新增品牌接口
// @Description 新增品牌接口
// @Tags 品牌
// @Accept application/json
// @Produce application/json
// @Param object body types.ParamTmSave true "品牌信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /admin/product/trademark [post]
func (t *trademarkHandler) CreateTrademark(c *gin.Context) {
	// 1. 获取参数及参数校验
	p := new(types.ParamTmSave)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		slog.Error("CreateTrademark with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}

	// 2. 业务处理
	err := t.tmSvc.CreateTrademark(c.Request.Context(), p)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}

	// 3. 返回结果
	result.Success(c, nil)
}

// GetTrademark 获取品牌列表
// @Summary 获取品牌分页列表
// @Description 获取品牌列表
// @Tags 品牌
// @Accept application/json
// @Produce application/json
// @Param page query int true "当前页码"
// @Param limit query int true "每页记录数"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=types.ResponseTmList}
// @Router /admin/product/trademark [get]
func (t *trademarkHandler) GetTrademark(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)

	// 获取数据
	data, err := t.tmSvc.GetTrademarkList(c.Request.Context(), page, size)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	// 返回响应
	result.Success(c, data)
}

// UpdateTrademark 更新品牌
// @Summary 更新品牌接口
// @Description 处理更新品牌请求
// @Tags 品牌
// @Accept application/json
// @Produce application/json
// @Param trademarkId path string true "品牌 ID"
// @Param object body types.ParamTmUpdate true "品牌信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /admin/product/trademark/{trademarkId} [put]
func (t *trademarkHandler) UpdateTrademark(c *gin.Context) {
	// 1. 获取参数及参数校验
	p := new(types.ParamTmUpdate)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		slog.Error("UpdateTrademark with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}

	tmId := idconv.ToInt64Safe(c.Param("trademarkId"))
	p.TmID = idconv.ToStr(tmId)

	// 2. 业务处理
	err := t.tmSvc.UpdateTrademark(c.Request.Context(), p)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}

	// 3. 返回结果
	result.Success(c, nil)
}

// DeleteTrademark 删除品牌
// @Summary 删除品牌接口
// @Description 处理删除品牌请求
// @Tags 品牌
// @Accept application/json
// @Produce application/json
// @Param trademarkId path string true "品牌 ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /admin/product/trademark/{trademarkId} [delete]
func (t *trademarkHandler) DeleteTrademark(c *gin.Context) {
	idStr := c.Param("trademarkId")
	tmId := idconv.ToInt64Safe(idStr)

	err := t.tmSvc.DeleteTrademark(c.Request.Context(), tmId)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}

	result.Success(c, nil)
}

// GetAllTrademarkList 获取所有品牌列表
// @Summary 获取所有品牌列表
// @Description 获取所有品牌列表
// @Tags 品牌
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=[]types.Trademark}
// @Router /admin/product/trademark/all [get]
func (t *trademarkHandler) GetAllTrademarkList(c *gin.Context) {
	data, err := t.tmSvc.GetAllTrademarkList(c.Request.Context())
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, data)
}
