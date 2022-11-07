package src

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type BusinessRepository struct{}

func (ts *BusinessRepository) GetBusinessById(storeId string) (*Business, *ServiceError) {

	db := GetDynamoDbClient()

	inpReq := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String("STORE#" + storeId),
			},
			"SK": {
				S: aws.String("STORE#" + storeId),
			},
		},
		TableName: aws.String(CommonTable),
	}

	res, err := db.GetItem(inpReq)

	if err != nil {
		return nil, &ServiceError{
			ErrorCode:    "100000",
			ErrorMessage: "Error Unparsing input request",
		}
	}

	if res.Item == nil {
		return nil, &ServiceError{
			ErrorCode:    "100001",
			ErrorMessage: "No Business Found",
		}
	}

	var data *Business
	err = dynamodbattribute.UnmarshalMap(res.Item, &data)

	if err != nil {
		return nil, &ServiceError{
			ErrorCode:    "100002",
			ErrorMessage: "Error Unparsing input request",
		}
	}

	return data, nil
}
