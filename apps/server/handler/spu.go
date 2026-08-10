package handler

import (
	"context"

	"github.com/liwook/go-vue-selection/pkg/errs"
	"github.com/liwook/go-vue-selection/pkg/idconv"
	"github.com/liwook/go-vue-selection/types"
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/liwook/go-vue-selection/pkg/result"
)

// SpuService 在使用方（controller）定义接口
type SpuService interface {
	GetSaleAttrList(context.Context) ([]*types.SaleAttr, error)
	SaveSpuInfo(context.Context, *types.Spu) error
	GetSpuList(context.Context, int64, int64, int64) (*types.ResponseSpuList, error)
	GetSpuImageList(context.Context, int64) ([]*types.SpuImage, error)
	GetSpuSaleAttrList(context.Context, int64) ([]*types.SpuSaleAttr, error)
	UpdateSpuInfo(context.Context, *types.Spu) error
	DeleteSpu(context.Context, int64) error
}

type spuHandler struct {
	spuSvc SpuService
}

func NewSpuHandler(svc SpuService) *spuHandler {
	return &spuHandler{spuSvc: svc}
}

// RegisterRoutes 注册商品 SPU 相关路由（需在 JWT 中间件之后调用）
func (s *spuHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/baseSaleAttr", s.GetSaleAttrList)
	r.POST("/spu", s.SaveSpuInfo)
	r.GET("/spu", s.GetSpuList)
	r.PUT("/spu/:spuId", s.UpdateSpuInfo)
	r.DELETE("/spu/:spuId", s.DeleteSpu)
	r.GET("/spu/:spuId/images", s.GetSpuImageList)
	r.GET("/spu/:spuId/saleAttr", s.GetSpuSaleAttrList)
}

// GetSaleAttrList 获取所有销售列表
// @Summary 获取所有销售列表
// @Description 获取所有销售列表
// @Tags 商品SPU
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=[]types.SaleAttr}
// @Router /api/v1/product/baseSaleAttr [get]
func (s *spuHandler) GetSaleAttrList(c *gin.Context) {
	saleAttrList, err := s.spuSvc.GetSaleAttrList(c.Request.Context())
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, saleAttrList)
}

// SaveSpuInfo 新增SPU接口
// @Summary 新增SPU接口
// @Description 新增SPU接口
// @Tags 商品SPU
// @Accept application/json
// @Produce application/json
// @Param object body types.Spu true "SPU信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/product/spu [post]
func (s *spuHandler) SaveSpuInfo(c *gin.Context) {
	p := new(types.Spu)
	if err := c.ShouldBindJSON(p); err != nil {
		slog.Error("SaveSpu with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}

	err := s.spuSvc.SaveSpuInfo(c.Request.Context(), p)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}

	result.Success(c, nil)
}

// GetSpuList 获取SPU分页列表
// @Summary 获取SPU分页列表
// @Description 获取SPU分页列表
// @Tags 商品SPU
// @Accept application/json
// @Produce application/json
// @Param page query int true "当前页码"
// @Param limit query int true "每页记录数"
// @Param category3Id query int true "三级分类 ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=types.ResponseSpuList}
// @Router /api/v1/product/spu [get]
func (s *spuHandler) GetSpuList(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	// 获取三级分类id
	c3IdStr := c.Query("category3Id")
	c3Id := idconv.ToInt64Safe(c3IdStr)

	// 获取列表数据
	responseSpuList, err := s.spuSvc.GetSpuList(c.Request.Context(), c3Id, page, size)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, responseSpuList)

}

// GetSpuImageList 获取商品图片列表
// @Summary 获取商品图片列表接口
// @Description 处理获取商品图片列表请求
// @Tags 商品SPU
// @Accept application/json
// @Produce application/json
// @Param spuId path string true "SPU ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=[]types.SpuImage}
// @Router /api/v1/product/spu/{spuId}/images [get]
func (s *spuHandler) GetSpuImageList(c *gin.Context) {
	idStr := c.Param("spuId")
	spuId := idconv.ToInt64Safe(idStr)
	spuImageList, err := s.spuSvc.GetSpuImageList(c.Request.Context(), spuId)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, spuImageList)
}

// GetSpuSaleAttrList 获取商品销售属性列表
// @Summary 获取商品销售属性列表接口
// @Description 处理获取商品销售属性列表请求
// @Tags 商品SPU
// @Accept application/json
// @Produce application/json
// @Param spuId path string true "SPU ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=[]types.SpuSaleAttr}
// @Router /api/v1/product/spu/{spuId}/saleAttr [get]
func (s *spuHandler) GetSpuSaleAttrList(c *gin.Context) {
	idStr := c.Param("spuId")
	spuId := idconv.ToInt64Safe(idStr)
	spuSaleAttrList, err := s.spuSvc.GetSpuSaleAttrList(c.Request.Context(), spuId)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, spuSaleAttrList)
}

// UpdateSpuInfo 更新SPU接口
// @Summary 更新SPU接口
// @Description 更新SPU接口
// @Tags 商品SPU
// @Accept application/json
// @Produce application/json
// @Param spuId path string true "SPU ID"
// @Param object body types.Spu true "SPU信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/product/spu/{spuId} [put]
func (s *spuHandler) UpdateSpuInfo(c *gin.Context) {
	p := new(types.Spu)
	if err := c.ShouldBindJSON(p); err != nil {
		slog.Error("UpdateSpuInfo with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}

	spuId := idconv.ToInt64Safe(c.Param("spuId"))
	p.SpuID = idconv.ToStr(spuId)

	err := s.spuSvc.UpdateSpuInfo(c.Request.Context(), p)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)
}

// DeleteSpu 删除SPU
// @Summary 删除SPU接口
// @Description 处理删除SPU请求
// @Tags 商品SPU
// @Accept application/json
// @Produce application/json
// @Param spuId path string true "SPU ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/product/spu/{spuId} [delete]
func (s *spuHandler) DeleteSpu(c *gin.Context) {
	idStr := c.Param("spuId")
	spuId := idconv.ToInt64Safe(idStr)
	err := s.spuSvc.DeleteSpu(c.Request.Context(), spuId)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)
}
