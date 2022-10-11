package src

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"log"
	svc "service"
	"time"
)

var (
	empRepo svc.IEmployee = &svc.EmployeeRepository{}
)

func CreateNewEmployee(req *svc.CreateStoreEmployeeRequest) (*svc.Employee, error) {

	_, err := empRepo.GetEmployeeId(req.Phone)
	if err != nil && err.(*svc.EmployeeError).Message != "User not found." {
		log.Printf("Error while fetching employee %v", err)
		return nil, err
	}

	// Create New Employee in Amazon Cognito
	ccuser, err := empRepo.CreateEmployeeOnCognito(req)
	if err != nil {
		log.Printf("Error while creating employee in cognito %v", err)
		return nil, err
	}

	userId := "EMP#" + *req.Phone
	// Save in the dynamodb
	strEmp := &svc.EmployeeDao{
		PK: &userId,
		SK: &userId,
		Employee: &svc.Employee{
			CreatedAt:  ccuser.UserCreateDate,
			EmployeeId: ccuser.Username,
			Phone:      req.Phone,
			Locale:     req.Locale,
			Type:       aws.String("EMP"),
		},
	}

	_, err = empRepo.CreateEmployee(strEmp)
	if err != nil {
		log.Printf("Error while creating employee in dynamodb %v", err)
		return nil, err
	}

	return strEmp.Employee, nil
}

func AssignRoleToUserForParticularStore(req *svc.CreateStoreEmployeeRequest) (*svc.StoreEmployeeRole, error) {

	storeId := "STORE#" + *req.BusinessId
	empId := "EMP#" + *req.Phone

	empRole := &svc.StoreEmployeeRoleDao{
		PK:   &empId,
		SK:   &storeId,
		GPK1: &storeId,
		GSK1: &empId,
		StoreEmployeeRole: &svc.StoreEmployeeRole{
			EmployeeId: req.Phone,
			StoreId:    req.BusinessId,
			Roles:      nil,
			Locale:     aws.String("en_US"),
			JoinedAt:   aws.Time(time.Now()),
			CreatedBy:  nil,
			CreatedAt:  aws.Time(time.Now()),
		},
	}

	_, err := empRepo.CreateNewEmployeeForStore(empRole)

	if err != nil {
		return nil, err
	}

	return empRole.StoreEmployeeRole, nil
}

func CreateNewBusinessEmployee(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req svc.CreateStoreEmployeeRequest
	err := json.Unmarshal([]byte(event.Body), &req)
	//eventJson, _ := json.MarshalIndent(event, "", "  ")
	//log.Printf("EVENT: %s", eventJson)

	_, err = empRepo.IsUserAlreadyExistForStore(req.BusinessId, req.Phone)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Headers:         nil,
			Body:            err.Error(),
			IsBase64Encoded: false,
		}, nil
	}

	// Create or get employee present in the user database
	_, err = CreateNewEmployee(&req)

	// Using the user fetched create user and store role assignation.
	_, err = AssignRoleToUserForParticularStore(&req)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Headers:         nil,
			Body:            err.Error(),
			IsBase64Encoded: false,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      201,
		Headers:         nil,
		IsBase64Encoded: false,
	}, nil
}
