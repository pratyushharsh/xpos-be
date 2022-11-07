package src

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"model"
	"strconv"
	"time"
)

type SyncRepository struct{}

func (sr *SyncRepository) UpdateSyncData(storeId string, req *UpdateSyncRequest) *ServiceError {
	db := GetDynamoDbClient()

	now := time.Now()
	utcTimeStamp := now.UnixMilli()
	//Convert the item to DynamoDB AttributeValues
	var writeRequest []*dynamodb.WriteRequest

	// For each transaction request create a write request
	for _, transaction := range *req.Transactions {

		pk := "STORE#" + storeId + "#TRAN#" + strconv.Itoa(*transaction.TransId)

		// Convert the item to DynamoDB AttributeValues
		dao := model.TransactionHeaderEntityDao{
			PK:                      &pk,
			SK:                      &pk,
			GPK1:                    &storeId,
			GSK1:                    &utcTimeStamp,
			TransactionHeaderEntity: transaction,
		}

		av, err := dynamodbattribute.MarshalMap(dao)
		if err != nil {
			return &ServiceError{
				ErrorCode:    "500",
				ErrorMessage: "Error Unparsing input request",
			}
		}
		writeRequest = append(writeRequest, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: av,
			},
		})
	}

	//Create a batch write request
	batchWriteRequest := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			DataTable: writeRequest,
		},
	}

	_, err := db.BatchWriteItem(batchWriteRequest)
	if err != nil {
		return &ServiceError{
			ErrorCode:    "500",
			ErrorMessage: err.Error(),
		}
	}

	return nil
}

func (sr *SyncRepository) GetSyncData(storeId string, timestamp *int64) (*UpdateSyncResponse, *ServiceError) {
	db := GetDynamoDbClient()

	// @TODO CHECK IF STORE NOT EXIST

	inpReq := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":storeId": {
				S: aws.String(storeId),
			},
		},
		IndexName:              aws.String("GPK1-GSK1-index"),
		KeyConditionExpression: aws.String("GPK1 = :storeId"),
		TableName:              aws.String(DataTable),
	}

	if timestamp != nil {
		inpReq.ExpressionAttributeValues[":timestamp"] = &dynamodb.AttributeValue{
			N: aws.String(strconv.FormatInt(*timestamp, 10)),
		}
		inpReq.KeyConditionExpression = aws.String("GPK1 = :pk and GSK1 > :timestamp")
	}

	res, err := db.Query(inpReq)

	// Error in input request
	if err != nil {
		return nil, &ServiceError{
			ErrorCode:    "500",
			ErrorMessage: err.Error(),
		}
	}

	var data []*model.TransactionHeaderEntity
	mErr := dynamodbattribute.UnmarshalListOfMaps(res.Items, &data)

	if mErr != nil {
		return nil, &ServiceError{
			ErrorCode:    "500",
			ErrorMessage: mErr.Error(),
		}
	}

	return &UpdateSyncResponse{
		Transactions: &data,
	}, nil
}
