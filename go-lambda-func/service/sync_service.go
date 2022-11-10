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

var (
	transactionType = "TRN"
	productType     = "PROD"
	customerType    = "CUS"
)

func (sr *SyncRepository) UpdateSyncData(storeId string, req *UpdateSyncRequest) (*UpdateSyncResponse, *ServiceError) {
	db := GetDynamoDbClient()

	now := time.Now()
	utcTimeStamp := now.UnixMicro()
	//Convert the item to DynamoDB AttributeValues
	var writeRequest []*dynamodb.WriteRequest

	// For each transaction request create a write request
	for _, transaction := range *req.Transactions {

		pk := "STORE#" + storeId + "#TRAN#" + strconv.Itoa(*transaction.TransId)

		transaction.LastSyncedAt = &utcTimeStamp
		// Convert the item to DynamoDB AttributeValues
		dao := model.TransactionHeaderEntityDao{
			PK:                      &pk,
			SK:                      &pk,
			GPK1:                    &storeId,
			GSK1:                    &utcTimeStamp,
			Type:                    &transactionType,
			TransactionHeaderEntity: transaction,
		}

		av, err := dynamodbattribute.MarshalMap(dao)
		if err != nil {
			return nil, &ServiceError{
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

	// For each product request create a write request
	if req.Products != nil {
		for _, product := range *req.Products {

			pk := "STORE#" + storeId + "#PROD#" + *product.ProductId

			product.LastSyncAt = &utcTimeStamp
			// Convert the item to DynamoDB AttributeValues
			dao := model.ProductEntityDao{
				PK:            &pk,
				SK:            &pk,
				GPK1:          &storeId,
				GSK1:          &utcTimeStamp,
				Type:          &productType,
				ProductEntity: product,
			}

			av, err := dynamodbattribute.MarshalMap(dao)
			if err != nil {
				return nil, &ServiceError{
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
	}

	// For each customer request create a write request
	if req.Customers != nil {
		for _, customer := range *req.Customers {

			pk := "STORE#" + storeId + "#CUS#" + *customer.ContactId

			customer.LastSyncAt = &utcTimeStamp
			// Convert the item to DynamoDB AttributeValues
			dao := model.CustomerEntityDao{
				PK:             &pk,
				SK:             &pk,
				GPK1:           &storeId,
				GSK1:           &utcTimeStamp,
				Type:           &customerType,
				CustomerEntity: customer,
			}

			av, err := dynamodbattribute.MarshalMap(dao)
			if err != nil {
				return nil, &ServiceError{
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
	}

	//Create a batch write request
	batchWriteRequest := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			DataTable: writeRequest,
		},
	}

	_, err := db.BatchWriteItem(batchWriteRequest)
	if err != nil {
		return nil, &ServiceError{
			ErrorCode:    "500",
			ErrorMessage: err.Error(),
		}
	}

	return &UpdateSyncResponse{
		LastSyncedAt: &utcTimeStamp,
	}, nil
}

func (sr *SyncRepository) GetSyncData(storeId string, from *int64, to *int64) (*GetSyncResponse, *ServiceError) {
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

	if from != nil && to != nil {
		inpReq.ExpressionAttributeValues[":from"] = &dynamodb.AttributeValue{
			N: aws.String(strconv.FormatInt(*from, 10)),
		}
		inpReq.ExpressionAttributeValues[":to"] = &dynamodb.AttributeValue{
			N: aws.String(strconv.FormatInt(*to, 10)),
		}
		inpReq.KeyConditionExpression = aws.String("GPK1 = :storeId AND GSK1 BETWEEN :from AND :to")
	} else if from != nil {
		inpReq.ExpressionAttributeValues[":from"] = &dynamodb.AttributeValue{N: aws.String(strconv.FormatInt(*from, 10))}
		inpReq.KeyConditionExpression = aws.String(*inpReq.KeyConditionExpression + " AND GSK1 >= :from")
	} else if to != nil {
		inpReq.ExpressionAttributeValues[":to"] = &dynamodb.AttributeValue{N: aws.String(strconv.FormatInt(*to, 10))}
		inpReq.KeyConditionExpression = aws.String(*inpReq.KeyConditionExpression + " AND GSK1 < :to")
	}

	res, err := db.Query(inpReq)

	// Error in input request
	if err != nil {
		return nil, &ServiceError{
			ErrorCode:    "500",
			ErrorMessage: err.Error(),
		}
	}

	// For each item parse check the type and add to the response
	var products []*model.ProductEntity
	var customers []*model.CustomerEntity
	var transactions []*model.TransactionHeaderEntity

	for _, item := range res.Items {
		typ := item["Type"].S
		if *typ == productType {
			var tp model.ProductEntity
			_ = dynamodbattribute.UnmarshalMap(item, &tp)
			products = append(products, &tp)
		} else if *typ == customerType {
			var tp model.CustomerEntity
			_ = dynamodbattribute.UnmarshalMap(item, &tp)
			customers = append(customers, &tp)
		} else if *typ == transactionType {
			var tp model.TransactionHeaderEntity
			_ = dynamodbattribute.UnmarshalMap(item, &tp)
			transactions = append(transactions, &tp)
		}
	}

	return &GetSyncResponse{
		Transactions: struct {
			Data *[]*model.TransactionHeaderEntity `json:"data"`
			From *int64                            `json:"from"`
			To   *int64                            `json:"to"`
		}{
			Data: &transactions,
			From: from,
			To:   to,
		},
		Customers: struct {
			Data *[]*model.CustomerEntity `json:"data"`
			From *int64                   `json:"from"`
			To   *int64                   `json:"to"`
		}{
			Data: &customers,
			From: from,
			To:   to,
		},
		Products: struct {
			Data *[]*model.ProductEntity `json:"data"`
			From *int64                  `json:"from"`
			To   *int64                  `json:"to"`
		}{
			Data: &products,
			From: from,
			To:   to,
		},
	}, nil
}
