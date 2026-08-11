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

// SkuService 在使用方（controller）定义接口
type SkuService interface {
	SaveSkuInfo(context.Context, *types.SkuInfo) error
	FindBySpuId(context.Context, int64) ([]*types.ResponseSkuInfo, error)
	GetSkuList(context.Context, int64, int64) (*types.ResponseSkuInfoList, error)
	OnSaleSku(context.Context, int64) error
	CancelSaleSku(context.Context, int64) error
	DeleteSku(context.Context, int64) error
	GetSkuInfo(context.Context, int64) (*types.ResponseSkuInfo, error)
}

type skuHandler struct {
	skuSvc SkuService
}

func NewSkuHandler(svc SkuService) *skuHandler {
	return &skuHandler{skuSvc: svc}
}

// RegisterRoutes 注册商品 SKU 相关路由（需在 JWT 中间件之后调用）
func (s *skuHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/sku", s.SaveSkuInfo)
	r.GET("/sku", s.GetSkuList)
	r.GET("/spu/:spuId/sku", s.FindBySpuId)
	r.PUT("/sku/:skuId/onsale", s.OnSaleSku)
	r.PUT("/sku/:skuId/cancelsale", s.CancelSaleSku)
	r.DELETE("/sku/:skuId", s.DeleteSku)
	r.GET("/sku/:skuId", s.GetSkuInfo)
}

// SaveSkuInfo 新增SKU接口
// @Summary 新增SKU接口
// @Description 新增SKU接口
// @Tags 商品SKU
// @Accept application/json
// @Produce application/json
// @Param object body types.SkuInfo true "SKU信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/product/sku [post]
// @Failure 500 {object} result.ResponseData
func (s *skuHandler) SaveSkuInfo(c *gin.Context) {
	p := new(types.SkuInfo)
	if err := c.ShouldBindJSON(p); err != nil {
		slog.Error("SaveSkuInfo with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}

	err := s.skuSvc.SaveSkuInfo(c.Request.Context(), p)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}

	result.Success(c, nil)
}

// FindBySpuId 根据 SPU ID 查询 SKU 列表
// @Summary 根据 SPU ID 查询 SKU 接口
// @Description 处理根据 SPU ID 查询 SKU列表请求
// @Tags 商品SKU
// @Accept application/json
// @Produce application/json
// @Param spuId path string true "SPU ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=[]types.ResponseSkuInfo}
// @Router /api/v1/product/spu/{spuId}/sku [get]
// @Failure 500 {object} result.ResponseData
func (s *skuHandler) FindBySpuId(c *gin.Context) {
	idStr := c.Param("spuId")
	spuId := idconv.ToInt64Safe(idStr)

	skuList, err := s.skuSvc.FindBySpuId(c.Request.Context(), spuId)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}

	result.Success(c, skuList)
}

// GetSkuList 获取SKU分页列表
// @Summary 获取SKU分页列表
// @Description 获取SKU分页列表
// @Tags 商品SKU
// @Accept application/json
// @Produce application/json
// @Param page query int true "当前页码"
// @Param limit query int true "每页记录数"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=types.ResponseSkuInfoList}
// @Router /api/v1/product/sku [get]
// @Failure 500 {object} result.ResponseData
func (s *skuHandler) GetSkuList(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)

	skuList, err := s.skuSvc.GetSkuList(c.Request.Context(), page, size)

	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}

	result.Success(c, skuList)
}

// OnSaleSku 上架SKU
// @Summary 上架SKU接口
// @Description 处理上架SKU请求
// @Tags 商品SKU
// @Accept application/json
// @Produce application/json
// @Param skuId path string true "SKU ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/product/sku/{skuId}/onsale [put]
// @Failure 500 {object} result.ResponseData
func (s *skuHandler) OnSaleSku(c *gin.Context) {
	idStr := c.Param("skuId")
	skuId := idconv.ToInt64Safe(idStr)
	err := s.skuSvc.OnSaleSku(c.Request.Context(), skuId)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)

}

// CancelSaleSku 下架SKU
// @Summary 下架SKU接口
// @Description 处理下架SKU请求
// @Tags 商品SKU
// @Accept application/json
// @Produce application/json
// @Param skuId path string true "SKU ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/product/sku/{skuId}/cancelsale [put]
// @Failure 500 {object} result.ResponseData
func (s *skuHandler) CancelSaleSku(c *gin.Context) {
	idStr := c.Param("skuId")
	skuId := idconv.ToInt64Safe(idStr)
	err := s.skuSvc.CancelSaleSku(c.Request.Context(), skuId)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)

}

// DeleteSku 删除SKU
// @Summary 删除SKU接口
// @Description 处理删除SKU请求
// @Tags 商品SKU
// @Accept application/json
// @Produce application/json
// @Param skuId path string true "SKU ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/product/sku/{skuId} [delete]
// @Failure 500 {object} result.ResponseData
func (s *skuHandler) DeleteSku(c *gin.Context) {
	idStr := c.Param("skuId")
	skuId := idconv.ToInt64Safe(idStr)
	err := s.skuSvc.DeleteSku(c.Request.Context(), skuId)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)
}

// GetSkuInfo 获取SKU详情
// @Summary 获取SKU详情接口
// @Description 处理获取SKU详情请求
// @Tags 商品SKU
// @Accept application/json
// @Produce application/json
// @Param skuId path string true "SKU ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=types.ResponseSkuInfo}
// @Router /api/v1/product/sku/{skuId} [get]
// @Failure 500 {object} result.ResponseData
func (s *skuHandler) GetSkuInfo(c *gin.Context) {
	idStr := c.Param("skuId")
	skuId := idconv.ToInt64Safe(idStr)
	skuInfo, err := s.skuSvc.GetSkuInfo(c.Request.Context(), skuId)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, skuInfo)
}
