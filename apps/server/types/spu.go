package types

type SaleAttr struct {
	BaseModel

	SaleAttrID   string `json:"saleAttrId"` // 销售属性ID
	SaleAttrName string `json:"name"`       // 销售属性名称
}

type SaleAttrValue struct {
	BaseModel
	SaleAttrValueID   string `json:"saleAttrValueId"`
	SaleAttrValueName string `json:"saleAttrValueName"`
	BaseSaleAttrId    string `json:"baseSaleAttrId"`
	SpuID             string `json:"spuId"`
}

type SpuSaleAttr struct {
	SpuSaleAttrID    string           `json:"spuSaleAttrId"`
	BaseSaleAttrId   string           `json:"baseSaleAttrId"`
	SaleAttrName     string           `json:"saleAttrName"`
	SpuID            string           `json:"spuId"`
	SpuSaleAttrValue []*SaleAttrValue `json:"spuSaleAttrValueList"`
}
type SpuImage struct {
	BaseModel
	ImageID   string `json:"imgId"`
	ImageName string `json:"imgName"`
	ImageUrl  string `json:"imgUrl"`
	SpuID     string `json:"spuId"`
}

type Spu struct {
	BaseModel
	SpuID           string         `json:"spuId"`
	SpuName         string         `json:"spuName"`
	Description     string         `json:"description"`
	Category3ID     string         `json:"category3Id"`
	TmID            string         `json:"tmId"`
	SpuImageList    []*SpuImage    `json:"spuImageList"`
	SpuSaleAttrList []*SpuSaleAttr `json:"spuSaleAttrList"`
}

type ResponseSpuList struct {
	Records     []*Spu `json:"records"`
	Total       int64  `json:"total"`
	Size        int64  `json:"size"`
	Current     int64  `json:"current"`
	SearchCount bool   `json:"searchCount"`
	Pages       int64  `json:"pages"`
}
