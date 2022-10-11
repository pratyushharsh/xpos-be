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

func GetEmployeeFromBusinessHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	pathParams := event.PathParameters

	bId := pathParams["businessId"]

	db := svc.GetDynamoDbClient()

	// @TODO Check if business not exist

	//uId := "USER#6cde1410-2a42-443d-a64a-bb23ed5190a0"

	//inpReq := &dynamodb.QueryInput{
	//	ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
	//		":pk": {
	//			S: aws.String(uId),
	//		},
	//		":beginsWith": {
	//			S: aws.String("STORE#"),
	//		},
	//	},
	//	IndexName:                 aws.String("GPK1-GSK1-index"),
	//	KeyConditionExpression:    aws.String("GPK1 = :pk and begins_with(GSK1, :beginsWith)"),
	//	TableName:                 aws.String("XPOS_DEV"),
	//}

	inpReq := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String("STORE#" + bId),
			},
			":beginsWith": {
				S: aws.String("EMP#"),
			},
		},
		IndexName:              aws.String("GPK1-GSK1-index"),
		KeyConditionExpression: aws.String("GPK1 = :pk and begins_with(GSK1, :beginsWith)"),
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

	var data []*svc.StoreEmployeeRole
	_ = dynamodbattribute.UnmarshalListOfMaps(res.Items, &data)

	byteData, _ := json.Marshal(data)

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         nil,
		Body:            string(byteData),
		IsBase64Encoded: false,
	}, nil
}
