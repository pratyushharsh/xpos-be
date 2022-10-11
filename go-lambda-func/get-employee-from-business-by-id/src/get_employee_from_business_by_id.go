package src

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
	svc "service"
)

var (
	empRepo         svc.IEmployee = &svc.EmployeeRepository{}
	CognitoUserPool               = os.Getenv("CognitoUserPool")
	CommonTable                   = os.Getenv("DBTable")
)

func GetUserFromCognito(username string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	cc := svc.GetCognitoClient()

	user, err := cc.AdminGetUser(&cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(CognitoUserPool),
		Username:   aws.String(username),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func GetEmployeeFromBusinessById(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	pathParams := event.PathParameters
	bId := pathParams["businessId"]
	userId := pathParams["userid"]

	db := svc.GetDynamoDbClient()

	inpReq := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String("EMP#" + userId),
			},
			"SK": {
				S: aws.String("STORE#" + bId),
			},
		},
		TableName: aws.String(CommonTable),
	}

	res, err := db.GetItem(inpReq)

	// Error in input request
	if err != nil {
		log.Println(err)
		msg := &svc.XPOSApiError{
			ErrorMessage: "Error getting the data from db",
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
		return events.APIGatewayProxyResponse{
			StatusCode:      404,
			Headers:         nil,
			IsBase64Encoded: false,
		}, nil
	}

	var data *svc.StoreEmployeeRole
	_ = dynamodbattribute.UnmarshalMap(res.Item, &data)

	user, _ := empRepo.GetEmployeeId(&userId)

	resp := &svc.StoreEmployeeResponse{
		Employee:  user,
		StoreData: data,
	}

	byteData, _ := json.Marshal(resp)

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         nil,
		Body:            string(byteData),
		IsBase64Encoded: false,
	}, nil
}
