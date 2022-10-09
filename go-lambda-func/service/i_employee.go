package src

import "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

type IEmployee interface {
	GetEmployeeId(employeeId *string) (*Employee, error)
	GetAllEmployeeFromStore(storeId *string) (*[]Employee, error)
	CreateEmployee(request *EmployeeDao) (*Employee, error)
	CreateNewEmployeeForStore(req *StoreEmployeeRoleDao) (*StoreEmployeeRole, error)
	UpdateEmployee(request *EmployeeUpdateRequest, username *string) error
	GetStoreAssignedToEmployee(employeeId *string) (*[]Business, error)
	CreateEmployeeOnCognito(req *CreateStoreEmployeeRequest) (*cognitoidentityprovider.UserType, error)
	IsUserAlreadyExistForStore(storeId *string, userId *string) (*bool, error)
}

type IEmployeeRelation interface {
	AssignRoleToEmployeeForStore()
}

type EmployeeError struct {
	CausedBy string `json:"caused_by"`
	Message  string `json:"message"`
}

func (ee *EmployeeError) Error() string {
	return ee.Message
}
