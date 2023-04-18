package src

type ErrorCode string

type TransactionServiceError struct {
	*ServiceError
}

type ServiceError struct {
	ErrorCode    ErrorCode `json:"errorCode"`
	ErrorMessage string    `json:"errorMessage"`
}

func (e *ServiceError) Error() string {
	return string(e.ErrorCode) + "->" + e.ErrorMessage
}
