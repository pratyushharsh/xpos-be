package src

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	svc "service"
)

var (
	taxRepo      svc.ITax      = &svc.TaxRepository{}
	businessRepo svc.IBusiness = &svc.BusinessRepository{}
)

func GetTaxGroupForStore(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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

	group, err := taxRepo.GetAllTaxGroupForStore(bId)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	res, _ := json.Marshal(group)
	return events.APIGatewayProxyResponse{Body: string(res), StatusCode: 200}, nil
}
