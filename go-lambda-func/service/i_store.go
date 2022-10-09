package src

type IStore interface {
	GetBusiness(businessId *string) (*Business, error)
	CreateBusiness(request *BusinessDao) (*Business, error)
	UpdateBusiness(request *BusinessDao) (*Business, error)
}

type ApiError struct {
	CausedBy string `json:"caused_by"`
	Message  string `json:"message"`
}

type StoreError struct {
	*ApiError
}

func (ee *StoreError) Error() string {
	return ee.CausedBy + " " + ee.Message
}
