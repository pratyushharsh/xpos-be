package src

import "model"

type ITax interface {
	GetAllTaxGroupForStore(storeId string) (*[]model.TaxGroup, error)
	CreateTaxGroupForStore(storeId string, request *[]model.TaxGroup) (*[]model.TaxGroup, error)
}
