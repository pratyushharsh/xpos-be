package src

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
	svc "service"
)

var (
	CommonTable = os.Getenv("DBTable")
)

func GetBusinessForUserHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	pathParams := event.PathParameters
	bId := pathParams["userid"]
	db := svc.GetDynamoDbClient()

	inpReq := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String("EMP#" + bId),
			},
			":beginsWith": {
				S: aws.String("STORE#"),
			},
		},
		KeyConditionExpression: aws.String("PK = :pk and begins_with(SK, :beginsWith)"),
		TableName:              aws.String(CommonTable),
	}

	res, err := db.Query(inpReq)

	// Error in input request
	if err != nil {
		log.Println(err)
		msg := &svc.XPOSApiError{
			ErrorMessage: "Error Unparsing input request",
		}
		res, _ := json.Marshal(msg)

		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Headers:         nil,
			Body:            string(res),
			IsBase64Encoded: false,
		}, nil
	}

	if res.Items == nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      404,
			Headers:         nil,
			IsBase64Encoded: false,
		}, nil
	}

	var data *[]svc.StoreEmployeeRole
	_ = dynamodbattribute.UnmarshalListOfMaps(res.Items, &data)

	byteData, _ := json.Marshal(data)

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         nil,
		Body:            string(byteData),
		IsBase64Encoded: false,
	}, nil
}
