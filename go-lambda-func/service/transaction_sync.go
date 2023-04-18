package src

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"model"
)

func BuildWriteRequestForTransaction(storeId *string, utcTimeStamp *int64, transactions *[]*model.TransactionHeaderEntity) (*[]*model.TransactionHeaderEntity, *[]*TransactionServiceError) {
	var res []*model.TransactionHeaderEntity
	var trnErr []*TransactionServiceError
	for _, transaction := range *transactions {
		data, err := InsertUpdateTransaction(storeId, utcTimeStamp, transaction)
		res = append(res, data)
		if err != nil {
			trnErr = append(trnErr, err)
		}
	}
	return &res, &trnErr
}

func InsertUpdateTransaction(storeId *string, utcTimeStamp *int64, transaction *model.TransactionHeaderEntity) (*model.TransactionHeaderEntity, *TransactionServiceError) {
	db := GetDynamoDbClient()
	pk := "STORE#" + *storeId
	sk := "TRAN#" + *transaction.TransId

	transaction.LastSyncedAt = utcTimeStamp
	// Convert the item to DynamoDB AttributeValues
	dao := model.TransactionHeaderEntityDao{
		PK:                      &pk,
		SK:                      &sk,
		GPK1:                    storeId,
		GSK1:                    utcTimeStamp,
		Type:                    &transactionType,
		TransactionHeaderEntity: transaction,
	}

	av, err := dynamodbattribute.MarshalMap(dao)
	if err != nil {
		return nil, &TransactionServiceError{
			&ServiceError{
				ErrorCode:    "100000",
				ErrorMessage: "Error Unparsing input request",
			},
		}
	}

	// Try inserting the data in the dynamodb first if pk exception is thrown then merge the database
	_, err = db.PutItem(&dynamodb.PutItemInput{
		ConditionExpression:         aws.String("attribute_not_exists(PK)"),
		ConditionalOperator:         nil,
		Expected:                    nil,
		Item:                        av,
		ReturnConsumedCapacity:      aws.String("TOTAL"),
		ReturnItemCollectionMetrics: nil,
		ReturnValues:                aws.String("NONE"),
		TableName:                   aws.String(DataTable),
	})

	if err == nil {
		return transaction, nil
	}

	switch err.(type) {
	case *dynamodb.ConditionalCheckFailedException:
		// If the item already exists then update the item
		log.Printf("Item already exists, updating the item")
		break
	default:
		return nil, &TransactionServiceError{
			&ServiceError{
				ErrorCode:    "100001",
				ErrorMessage: "Error inserting the item",
			},
		}
	}

	inpReq := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(pk),
			},
			"SK": {
				S: aws.String(sk),
			},
		},
		TableName: aws.String(DataTable),
	}

	// Get the transaction from the dynamo db
	result, err := db.GetItem(inpReq)

	// If the transaction is not present in the dynamo db then return the error
	// If the transaction is present in the dynamo db then merge the data
	var t model.TransactionHeaderEntityDao
	err = dynamodbattribute.UnmarshalMap(result.Item, &t)

	// Check the lastTimestamp of the transaction in the dynamo db
	// If the lastTimestamp is greater than the current timestamp then return the error
	if *t.LastChangedAt < *transaction.LastChangedAt {
		// Insert into the dynamo db
		dao := model.TransactionHeaderEntityDao{
			PK:                      &pk,
			SK:                      &sk,
			GPK1:                    storeId,
			GSK1:                    utcTimeStamp,
			Type:                    &transactionType,
			TransactionHeaderEntity: transaction,
		}

		av, err := dynamodbattribute.MarshalMap(dao)
		if err != nil {
			return nil, &TransactionServiceError{
				&ServiceError{
					ErrorCode:    "100000",
					ErrorMessage: "Error Unparsing input request",
				},
			}
		}

		// Try inserting the data in the dynamodb first if pk exception is thrown then merge the database
		_, err = db.PutItem(&dynamodb.PutItemInput{
			ConditionalOperator:         nil,
			Expected:                    nil,
			Item:                        av,
			ReturnConsumedCapacity:      aws.String("TOTAL"),
			ReturnItemCollectionMetrics: nil,
			ReturnValues:                aws.String("NONE"),
			TableName:                   aws.String(DataTable),
		})

		if err != nil {
			return nil, &TransactionServiceError{
				&ServiceError{
					ErrorCode:    "100001",
					ErrorMessage: "Error inserting the item",
				},
			}
		}
	}

	return transaction, nil
}
