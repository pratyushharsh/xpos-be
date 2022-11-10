package src

import "model"

type UpdateSyncRequest struct {
	Transactions *[]*model.TransactionHeaderEntity `json:"transactions"`
	Customers    *[]*model.CustomerEntity          `json:"customers"`
	Products     *[]*model.ProductEntity           `json:"products"`
}

type UpdateSyncResponse struct {
	ErrorKey     *[]*string `json:"errorKey"`
	LastSyncedAt *int64     `json:"lastSyncedAt"`
}

type GetSyncResponse struct {
	Transactions struct {
		Data *[]*model.TransactionHeaderEntity `json:"data"`
		From *int64                            `json:"from"`
		To   *int64                            `json:"to"`
	} `json:"transactions"`
	Customers struct {
		Data *[]*model.CustomerEntity `json:"data"`
		From *int64                   `json:"from"`
		To   *int64                   `json:"to"`
	} `json:"customers"`
	Products struct {
		Data *[]*model.ProductEntity `json:"data"`
		From *int64                  `json:"from"`
		To   *int64                  `json:"to"`
	} `json:"products"`
}
