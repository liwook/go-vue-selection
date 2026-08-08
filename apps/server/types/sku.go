package types

type Sku struct {
	BaseModel
	SkuID       string `json:"skuId"`
	SpuID       string `json:"spuId"`
	Category3ID string `json:"category3Id"`
	TmID        string `json:"tmId"`
	SkuName     string `json:"skuName"`
	// WeightMg 重量，单位：毫克（最小重量单位，milligrams）。JSON number 类型，避免浮点误差。
	WeightMg int64 `json:"weightMg"`
	// PriceCent 价格，单位：分（最小货币单位，cents）。JSON number 类型，避免浮点误差。
	PriceCent int64  `json:"priceCent"`
	SkuDesc   string `json:"skuDesc"`
	IsSale    int8   `json:"isSale"`
}

type SkuAttrValue struct {
	BaseModel
	SkuAttrValueID string `json:"skuAttrValueId"`
	AttrID         string `json:"attrId"`  // 平台属性ID
	ValueID        string `json:"valueId"` // 属性值ID
	ValueName      string `json:"valueName"`
	AttrName       string `json:"attrName"`
	SkuID          string `json:"skuId"`
}

type SkuSaleAttrValue struct {
	BaseModel
	SkuSaleAttrValueID string `json:"skuSaleAttrValueId"`
	SaleAttrID         string `json:"saleAttrId"` // 销售属性ID
	SaleAttrValueID    string `json:"saleAttrValueId"`
	SaleAttrName       string `json:"saleAttrName"`
	SaleAttrValueName  string `json:"saleAttrValueName"`
	SkuID              string `json:"skuId"`
}

type SkuImg struct {
	BaseModel
	ImageID    string `json:"imgId"`
	SkuID      string `json:"skuId"`
	ImageName  string `json:"imgName"`
	ImageURL   string `json:"imgUrl"`
	SpuImageID string `json:"spuImgId"`
	IsDefault  string `json:"isDefault"`
}

type SkuAttrValueDTO struct {
	SkuAttrValueID string `json:"skuAttrValueId"`
	AttrID         string `json:"attrId"`  // 平台属性ID
	ValueID        string `json:"valueId"` // 属性值ID
	ValueName      string `json:"valueName"`
	AttrName       string `json:"attrName"`
	SkuID          string `json:"skuId"`
}

type SkuSaleAttrValueDTO struct {
	SkuSaleAttrValueID string `json:"skuSaleAttrValueId"`
	SaleAttrID         string `json:"saleAttrId"` // 销售属性ID
	SaleAttrValueID    string `json:"saleAttrValueId"`
	SaleAttrName       string `json:"saleAttrName"`
	SaleAttrValueName  string `json:"saleAttrValueName"`
	SkuID              string `json:"skuId"`
}

type SkuImgDTO struct {
	ImageID    string `json:"imgId"`
	SkuID      string `json:"skuId"`
	ImageName  string `json:"imgName"`
	ImageURL   string `json:"imgUrl"`
	SpuImageID string `json:"spuImgId"`
	IsDefault  string `json:"isDefault"`
}

type SkuInfo struct {
	SkuID       string `json:"skuId"`
	SpuID       string `json:"spuId"`
	Category3ID string `json:"category3Id"`
	TmID        string `json:"tmId"`
	SkuName     string `json:"skuName"`
	// WeightMg 重量，单位：毫克（最小重量单位，milligrams）。直接以“毫克”对外返回，避免浮点误差。
	WeightMg int64 `json:"weightMg"`
	// PriceCent 价格，单位：分（最小货币单位，cents）。直接以“分”对外返回，避免浮点误差。
	PriceCent            int64                  `json:"priceCent"`
	SkuDesc              string                 `json:"skuDesc"`
	IsSale               int8                   `json:"isSale"`
	SkuAttrValueList     []*SkuAttrValueDTO     `json:"skuAttrValueList"`
	SkuSaleAttrValueList []*SkuSaleAttrValueDTO `json:"skuSaleAttrValueList"`
	SkuImageList         []*SkuImgDTO           `json:"skuImageList"`
}

type ResponseSkuInfo struct {
	Sku
	SkuAttrValueList     []*SkuAttrValue     `json:"skuAttrValueList"`
	SkuSaleAttrValueList []*SkuSaleAttrValue `json:"skuSaleAttrValueList"`
	SkuImageList         []*SkuImg           `json:"skuImageList"`
}

type ResponseSkuInfoList struct {
	Records     []*ResponseSkuInfo `json:"records"`
	Total       int64              `json:"total"`
	Size        int64              `json:"size"`
	Current     int64              `json:"current"`
	SearchCount bool               `json:"searchCount"`
	Pages       int64              `json:"pages"`
}
