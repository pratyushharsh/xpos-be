package src

type ServiceError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (e *ServiceError) Error() string {
	return e.ErrorCode + "->" + e.ErrorMessage
}
