package src

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	svc "service"
)

func GetBusinessByIdHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	pathParams := event.PathParameters
	bId := pathParams["businessId"]
	db := svc.GetDynamoDbClient()

	inpReq := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String("STORE#" + bId),
			},
			"SK": {
				S: aws.String("STORE#" + bId),
			},
		},
		TableName: aws.String("XPOS_DEV"),
	}

	res, err := db.GetItem(inpReq)

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

	if res.Item == nil {
		msg := &svc.XPOSApiError{
			ErrorMessage: "No Business Found",
		}

		res, _ := json.Marshal(msg)
		return events.APIGatewayProxyResponse{
			StatusCode:      404,
			Headers:         nil,
			IsBase64Encoded: false,
			Body:            string(res),
		}, nil
	}

	var data *svc.Business
	_ = dynamodbattribute.UnmarshalMap(res.Item, &data)

	byteData, _ := json.Marshal(data)

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         nil,
		Body:            string(byteData),
		IsBase64Encoded: false,
	}, nil
}
