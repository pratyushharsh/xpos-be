package src

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	svc "service"
)

var (
	syncRepo svc.SyncRepository = svc.SyncRepository{}
)

func UpdateSyncHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	eventJson, _ := json.Marshal(event)
	log.Printf("EVENT: %s", eventJson)

	pathParams := event.PathParameters
	bId := pathParams["businessId"]

	var req svc.UpdateSyncRequest
	err := json.Unmarshal([]byte(event.Body), &req)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Headers:         nil,
			Body:            err.Error(),
			IsBase64Encoded: false,
		}, nil
	}

	dbErr := syncRepo.UpdateSyncData(bId, &req)
	if dbErr != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Headers:         nil,
			Body:            dbErr.Error(),
			IsBase64Encoded: false,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
