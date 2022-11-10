package model

// PK STORE#<Store Id>#<>
// SK

type TransactionHeaderEntityDao struct {
	PK   *string `json:"PK"`
	SK   *string `json:"SK"`
	GPK1 *string `json:"GPK1"`
	GSK1 *int64  `json:"GSK1"`
	Type *string `json:"Type"`
	*TransactionHeaderEntity
}

type TransactionHeaderEntity struct {
	AssociateId      *string                   `json:"associateId"`
	AssociateName    *string                   `json:"associateName"`
	BeginDatetime    *int64                    `json:"beginDatetime"`
	BillingAddress   *string                   `json:"billingAddress"`
	BusinessDate     *int64                    `json:"businessDate"`
	CustomerId       *string                   `json:"customerId"`
	CustomerName     *string                   `json:"customerName"`
	CustomerPhone    *string                   `json:"customerPhone"`
	DiscountTotal    *float64                  `json:"discountTotal"`
	EndDateTime      *int64                    `json:"endDateTime"`
	IsVoid           *bool                     `json:"isVoid"`
	LastChangedAt    *int64                    `json:"lastChangedAt"`
	LastSyncedAt     *int64                    `json:"lastSyncedAt"`
	LineItems        *[]TransactionLineItem    `json:"lineItems"`
	Notes            *string                   `json:"notes"`
	PaymentLineItems *[]TransactionPaymentLine `json:"paymentLineItems"`
	RoundTotal       *float32                  `json:"roundTotal"`
	ShippingAddress  *string                   `json:"shippingAddress"`
	Status           *string                   `json:"status"`
	StoreCurrency    *string                   `json:"storeCurrency"`
	StoreId          *int                      `json:"storeId"`
	StoreLocale      *string                   `json:"storeLocale"`
	Subtotal         *float64                  `json:"subtotal"`
	TaxTotal         *float64                  `json:"taxTotal"`
	Total            *float64                  `json:"total"`
	TransId          *int                      `json:"transId"`
	TransactionType  *string                   `json:"transactionType"`
}

type TransactionLineItem struct {
	BaseUnitPrice        *float64                           `json:"baseUnitPrice"`
	BusinessDate         *int64                             `json:"businessDate"`
	Category             *string                            `json:"category"`
	Currency             *string                            `json:"currency"`
	DiscountAmount       *float64                           `json:"discountAmount"`
	ExtendedAmount       *float64                           `json:"extendedAmount"`
	GrossAmount          *float64                           `json:"grossAmount"`
	Hsn                  *string                            `json:"hsn"`
	IsVoid               *bool                              `json:"isVoid"`
	ItemDescription      *string                            `json:"itemDescription"`
	ItemId               *string                            `json:"itemId"`
	ItemIdEntryMethod    *string                            `json:"itemIdEntryMethod"`
	LineItemSeq          *int                               `json:"lineItemSeq"`
	LineModifiers        *[]TransactionLineItemModifier     `json:"lineModifiers"`
	NetAmount            *float64                           `json:"netAmount"`
	NonExchangeableFlag  *bool                              `json:"nonExchangeableFlag"`
	NonReturnableFlag    *bool                              `json:"nonReturnableFlag"`
	OriginalBusinessDate *int64                             `json:"originalBusinessDate"`
	OriginalLineItemSeq  *int64                             `json:"originalLineItemSeq"`
	OriginalPosId        *int                               `json:"originalPosId"`
	OriginalTransSeq     *int                               `json:"originalTransSeq"`
	PosId                *int                               `json:"posId"`
	PriceEntryMethod     *string                            `json:"priceEntryMethod"`
	PriceOverride        *bool                              `json:"priceOverride"`
	PriceOverrideAmount  *float64                           `json:"priceOverrideAmount"`
	PriceOverrideReason  *string                            `json:"priceOverrideReason"`
	Quantity             *float32                           `json:"quantity"`
	ReturnComment        *string                            `json:"returnComment"`
	ReturnFlag           *bool                              `json:"returnFlag"`
	ReturnReasonCode     *string                            `json:"returnReasonCode"`
	ReturnTypeCode       *string                            `json:"returnTypeCode"`
	ReturnedQuantity     *float32                           `json:"returnedQuantity"`
	SerialNumber         *string                            `json:"serialNumber"`
	ShippingWeight       *float64                           `json:"shippingWeight"`
	StoreId              *int                               `json:"storeId"`
	TaxAmount            *float64                           `json:"taxAmount"`
	TaxModifiers         *[]TransactionLineItemTaxModifiers `json:"taxModifiers"`
	TransSeq             *int                               `json:"transSeq"`
	UnitCost             *float64                           `json:"unitCost"`
	UnitPrice            *float64                           `json:"unitPrice"`
	Uom                  *string                            `json:"uom"`
	VendorId             *string                            `json:"vendorId"`
}

type TransactionLineItemModifier struct {
	Amount                  *float64 `json:"amount"`
	BusinessDate            *int64   `json:"businessDate"`
	DealId                  *string  `json:"dealId"`
	Description             *string  `json:"description"`
	DiscountCode            *string  `json:"discountCode"`
	DiscountReasonCode      *string  `json:"discountReasonCode"`
	ExtendedAmount          *float64 `json:"extendedAmount"`
	GroupDiscountId         *string  `json:"groupDiscountId"`
	IsVoid                  *bool    `json:"isVoid"`
	LineItemModSeq          *int     `json:"lineItemModSeq"`
	LineItemSeq             *int     `json:"lineItemSeq"`
	Notes                   *string  `json:"notes"`
	Percent                 *int     `json:"percent"`
	PosId                   *int     `json:"posId"`
	PriceModifierReasonCode *string  `json:"priceModifierReasonCode"`
	PromotionId             *string  `json:"promotionId"`
	StoreId                 *int     `json:"storeId"`
	TransSeq                *int     `json:"transSeq"`
}

type TransactionLineItemTaxModifiers struct {
	AuthorityId           *string  `json:"authorityId"`
	AuthorityName         *string  `json:"authorityName"`
	AuthorityType         *string  `json:"authorityType"`
	LineItemSeq           *int     `json:"lineItemSeq"`
	OriginalTaxableAmount *float64 `json:"originalTaxableAmount"`
	RawTaxAmount          *float64 `json:"rawTaxAmount"`
	RawTaxPercentage      *float32 `json:"rawTaxPercentage"`
	TaxAmount             *float64 `json:"taxAmount"`
	TaxGroupId            *string  `json:"taxGroupId"`
	TaxLocationId         *string  `json:"taxLocationId"`
	TaxOverride           *bool    `json:"taxOverride"`
	TaxOverrideAmount     *string  `json:"taxOverrideAmount"`
	TaxOverridePercent    *string  `json:"taxOverridePercent"`
	TaxOverrideReasonCode *string  `json:"taxOverrideReasonCode"`
	TaxPercent            *float64 `json:"taxPercent"`
	TaxRuleId             *string  `json:"taxRuleId"`
	TaxRuleName           *string  `json:"taxRuleName"`
	TaxableAmount         *float64 `json:"taxableAmount"`
	TransSeq              *int     `json:"transSeq"`
}

type TransactionPaymentLine struct {
	Amount           *float64 `json:"amount"`
	AuthCode         *string  `json:"authCode"`
	BeginDate        *int64   `json:"beginDate"`
	CurrencyId       *string  `json:"currencyId"`
	EndDate          *int64   `json:"endDate"`
	IsVoid           *bool    `json:"isVoid"`
	PaymentSeq       *int     `json:"paymentSeq"`
	TenderId         *string  `json:"tenderId"`
	TenderStatusCode *string  `json:"tenderStatusCode"`
	Token            *string  `json:"token"`
	TransId          *int     `json:"transId"`
}
