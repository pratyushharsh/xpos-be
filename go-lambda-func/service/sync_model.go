package src

import "model"

type UpdateSyncRequest struct {
	Transactions *[]*model.TransactionHeaderEntity `json:"transactions"`
}

type UpdateSyncResponse struct {
	Transactions *[]*model.TransactionHeaderEntity `json:"transactions"`
}
