package src

type IBusiness interface {
	GetBusinessById(businessId string) (*Business, *ServiceError)
}
