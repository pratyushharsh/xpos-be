package src

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	svc "service"
	"strconv"
	"time"
)

var (
	syncRepo svc.SyncRepository = svc.SyncRepository{}
)

func GetSyncDataHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	pathParams := event.PathParameters
	bId := pathParams["businessId"]

	from := event.QueryStringParameters["startTime"]

	i, _ := strconv.ParseInt(from, 10, 64)

	now := time.Now()
	to := now.UnixMicro()

	syncData, err := syncRepo.GetSyncData(bId, &i, &to)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	res, _ := json.Marshal(syncData)

	return events.APIGatewayProxyResponse{Body: string(res), StatusCode: 200}, nil
}
