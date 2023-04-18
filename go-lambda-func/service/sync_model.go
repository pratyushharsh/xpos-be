package src

import "model"

type UpdateSyncRequest struct {
	Transactions *[]*model.TransactionHeaderEntity `json:"transactions"`
	Customers    *[]*model.CustomerEntity          `json:"customers"`
	Products     *[]*model.ProductEntity           `json:"products"`
	Config       *SyncConfig                       `json:"config"`
}

type SyncConfig struct {
	TaxConfig      *[]*model.TaxGroupEntity     `json:"taxConfig"`
	InvoiceConfig  *[]*model.ReportConfigEntity `json:"invoiceConfig"`
	SequenceConfig *[]*model.SequenceEntity     `json:"sequenceConfig"`
}

type UpdateSyncResponse struct {
	Transactions *SyncResponse `json:"transactions"`
	Customers    *SyncResponse `json:"customers"`
	Products     *SyncResponse `json:"products"`
	ErrorKey     *[]*string   `json:"errorKey"`
	LastSyncedAt *int64       `json:"lastSyncedAt"`
}

type GetSyncResponse struct {
	Transactions SyncData     `json:"transactions"`
	Customers    SyncData     `json:"customers"`
	Products     SyncData     `json:"products"`
	Config       ConfigOutput `json:"config"`
}

type ConfigOutput struct {
	TaxConfig      SyncData `json:"taxConfig"`
	InvoiceConfig  SyncData `json:"invoiceConfig"`
	SequenceConfig SyncData `json:"sequenceConfig"`
}

type SyncResponse struct {
	Data         interface{} `json:"data"`
	Error        interface{} `json:"error"`
	LastSyncedAt *int64      `json:"lastSyncedAt"`
}

type SyncData struct {
	Data interface{} `json:"data"`
	From *int64      `json:"from"`
	To   *int64      `json:"to"`
}
