package src

import "model"

type ITax interface {
	GetAllTaxGroupForStore(storeId string) (*[]model.TaxGroupEntity, error)
	CreateTaxGroupForStore(storeId string, request *[]model.TaxGroupEntity) (*[]model.TaxGroupEntity, error)
}
