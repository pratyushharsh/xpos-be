package src

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"model"
	svc "service"
)

var (
	taxRepo      svc.ITax      = &svc.TaxRepository{}
	businessRepo svc.IBusiness = &svc.BusinessRepository{}
)

func CreateNewTaxGroupHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	pathParams := event.PathParameters
	bId := pathParams["businessId"]

	_, svcErr := businessRepo.GetBusinessById(bId)

	if svcErr != nil {

		msg := &svc.XPOSApiError{
			ErrorMessage: svcErr.Error(),
		}

		res, _ := json.Marshal(msg)

		return events.APIGatewayProxyResponse{Body: string(res), StatusCode: 500}, nil
	}

	// Check if business exist and get business name

	var req *[]model.TaxGroupEntity
	err := json.Unmarshal([]byte(event.Body), &req)

	// Check if store exists

	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	// Create New Tax Group
	_, err = taxRepo.CreateTaxGroupForStore(bId, req)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	return events.APIGatewayProxyResponse{Body: "Success", StatusCode: 200}, nil
}
