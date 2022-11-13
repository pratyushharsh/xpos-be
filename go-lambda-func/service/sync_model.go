package src

import "model"

type UpdateSyncRequest struct {
	Transactions *[]*model.TransactionHeaderEntity `json:"transactions"`
	Customers    *[]*model.CustomerEntity          `json:"customers"`
	Products     *[]*model.ProductEntity           `json:"products"`
	Config       *SyncConfig                       `json:"config"`
}

type SyncConfig struct {
	TaxConfig     *[]*model.TaxGroupEntity     `json:"taxConfig"`
	InvoiceConfig *[]*model.ReportConfigEntity `json:"invoiceConfig"`
}

type UpdateSyncResponse struct {
	ErrorKey     *[]*string `json:"errorKey"`
	LastSyncedAt *int64     `json:"lastSyncedAt"`
}

type GetSyncResponse struct {
	Transactions SyncData     `json:"transactions"`
	Customers    SyncData     `json:"customers"`
	Products     SyncData     `json:"products"`
	Config       ConfigOutput `json:"config"`
}

type ConfigOutput struct {
	TaxConfig     SyncData `json:"taxConfig"`
	InvoiceConfig SyncData `json:"invoiceConfig"`
}

type SyncData struct {
	Data interface{} `json:"data"`
	From *int64      `json:"from"`
	To   *int64      `json:"to"`
}
