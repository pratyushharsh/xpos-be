package model

type ProductEntityDao struct {
	PK   *string `json:"PK"`
	SK   *string `json:"SK"`
	GPK1 *string `json:"GPK1"`
	GSK1 *int64  `json:"GSK1"`
	Type *string `json:"Type"`
	*ProductEntity
}

type ProductEntity struct {
	Brand         *string   `json:"brand"`
	Category      *[]string `json:"category"`
	CreateTime    *int64    `json:"createTime"`
	Description   *string   `json:"description"`
	DisplayName   *string   `json:"displayName"`
	Enable        *bool     `json:"enable"`
	Hsn           *string   `json:"hsn"`
	ImageUrl      *[]string `json:"imageUrl"`
	LastSyncAt    *int64    `json:"lastSyncAt"`
	ListPrice     *float32  `json:"listPrice"`
	ProductId     *string   `json:"productId"`
	SalePrice     *float32  `json:"salePrice"`
	SkuCode       *string   `json:"skuCode"`
	SkuId         *string   `json:"skuId"`
	StoreId       *string   `json:"storeId"`
	TaxGroupId    *string   `json:"taxGroupId"`
	Uom           *string   `json:"uom"`
	LastChangedAt *int64    `json:"lastChangedAt"`
}
