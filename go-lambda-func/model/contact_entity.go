package model

type CustomerEntityDao struct {
	PK   *string `json:"PK"`
	SK   *string `json:"SK"`
	GPK1 *string `json:"GPK1"`
	GSK1 *int64  `json:"GSK1"`
	Type *string `json:"Type"`
	*CustomerEntity
}

type CustomerEntity struct {
	BillingAddress  *AddressEntity `json:"billingAddress"`
	ContactId       *string        `json:"contactId"`
	CreateTime      *int64         `json:"createTime"`
	Email           *string        `json:"email"`
	FirstName       *string        `json:"firstName"`
	Gstin           *string        `json:"gstin"`
	LastChangedAt   *int64         `json:"lastChangedAt"`
	LastName        *string        `json:"lastName"`
	LastSyncAt      *int64         `json:"lastSyncAt"`
	PanCard         *string        `json:"panCard"`
	PhoneNumber     *string        `json:"phoneNumber"`
	ShippingAddress *AddressEntity `json:"shippingAddress"`
	StoreId         *string        `json:"storeId"`
}
