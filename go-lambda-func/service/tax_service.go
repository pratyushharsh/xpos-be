package src

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"model"
)

//var (
//	CommonTable     = os.Getenv("DBTable")
//	CognitoUserPool = os.Getenv("CognitoUserPool")
//)

type TaxRepository struct{}

func (ts *TaxRepository) GetAllTaxGroupForStore(storeId string) (*[]model.TaxGroupEntity, error) {
	db := GetDynamoDbClient()

	// Get all the data for a store.
	inpReq := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String("STORE#" + storeId),
			},
			":beginsWith": {
				S: aws.String("TAX#"),
			},
		},
		KeyConditionExpression: aws.String("PK = :pk and begins_with(SK, :beginsWith)"),
		TableName:              aws.String(CommonTable),
	}

	// Parse the input request

	res, err := db.Query(inpReq)

	// Error in input request
	if err != nil {
		return nil, &ServiceError{
			ErrorCode:    "500",
			ErrorMessage: "Error Unparsing input request",
		}
	}

	var data *[]model.TaxGroupEntity
	_ = dynamodbattribute.UnmarshalListOfMaps(res.Items, &data)

	return data, nil
}

func (ts *TaxRepository) CreateTaxGroupForStore(storeId string, request *[]model.TaxGroupEntity) (*[]model.TaxGroupEntity, error) {
	//db := GetDynamoDbClient()
	//
	//// Convert the item to DynamoDB AttributeValues
	//var writeRequest []*dynamodb.WriteRequest

	//// For each tax group, create a write request
	//for _, taxGroup := range *request {
	//
	//	pk := "STORE#" + storeId
	//	sk := "TAX#" + *taxGroup.GroupId
	//
	//	taxGroup := &model.TaxGroupDao{
	//		PK:       &pk,
	//		SK:       &sk,
	//		TaxGroupEntity: &taxGroup,
	//	}
	//
	//	dbMap, _ := dynamodbattribute.MarshalMap(taxGroup)
	//
	//	writeRequest = append(writeRequest, &dynamodb.WriteRequest{
	//		PutRequest: &dynamodb.PutRequest{
	//			Item: dbMap,
	//		},
	//	})
	//}
	//
	//req := dynamodb.BatchWriteItemInput{
	//	RequestItems: map[string][]*dynamodb.WriteRequest{
	//		CommonTable: writeRequest,
	//	},
	//}
	//
	//_, err := db.BatchWriteItem(&req)
	//if err != nil {
	//	return nil, err
	//}
	//
	return nil, nil
}
