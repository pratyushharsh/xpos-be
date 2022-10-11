package src

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	svc "service"
)

var (
	empRepo svc.IEmployee = &svc.EmployeeRepository{}
)

func UpdateEmployeeHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req svc.EmployeeUpdateRequest
	err := json.Unmarshal([]byte(event.Body), &req)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Headers:         nil,
			Body:            err.Error(),
			IsBase64Encoded: false,
		}, nil
	}
	pathParams := event.PathParameters
	userId := pathParams["userid"]

	err = empRepo.UpdateEmployee(&req, &userId)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Headers:         nil,
			Body:            err.Error(),
			IsBase64Encoded: false,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
