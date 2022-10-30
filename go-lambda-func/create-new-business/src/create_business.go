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
	"strconv"
	"time"
)

var (
	empRepo     svc.IEmployee = &svc.EmployeeRepository{}
	CommonTable               = os.Getenv("DBTable")
)

func GetNewBusinessID() (*int, error) {
	db := svc.GetDynamoDbClient()

	input := &dynamodb.UpdateItemInput{
		ConditionExpression: aws.String("attribute_exists(sequenceType)"),
		ExpressionAttributeNames: map[string]*string{
			"#sequenceValue": aws.String("sequenceValue"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":sequenceValue": {
				N: aws.String("1"),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"sequenceType": {
				S: aws.String("XPOS_BUSINESS"),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String("THELAWALA_SEQUENCE"),
		UpdateExpression: aws.String("SET #sequenceValue = #sequenceValue + :sequenceValue"),
	}

	res, err := db.UpdateItem(input)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var seq svc.Sequence
	_ = dynamodbattribute.UnmarshalMap(res.Attributes, &seq)

	return seq.SequenceValue, nil
}

func createNewBusinessInputRequest(businessId *int, businessDetail *svc.CreateBusinessRequest, createUser *string) *svc.CreateBusinessResponse {
	now := time.Now()
	return &svc.CreateBusinessResponse{
		BusinessId: businessId,
		Name:       businessDetail.Name,
		LegalName:  businessDetail.LegalName,
		Email:      businessDetail.Email,
		Address1:   businessDetail.Address1,
		Address2:   businessDetail.Address2,
		City:       businessDetail.City,
		State:      businessDetail.State,
		PostalCode: businessDetail.PostalCode,
		Country:    businessDetail.Country,
		Currency:   businessDetail.Currency,
		Phone:      businessDetail.Phone,
		Locale:     businessDetail.Locale,
		Gst:        businessDetail.Gst,
		Pan:        businessDetail.Pan,
		CreatedBy:  createUser,
		CreatedAt:  &now,
	}
}

func CreateBusinessInDb(req *svc.CreateBusinessResponse) *svc.CreateBusinessResponse {

	businessId := strconv.Itoa(*req.BusinessId)

	pk := "STORE#" + businessId
	//gpk1 := "USER#" + *req.CreatedBy
	//gsk1 := "STORE#" + businessId

	dbReq := &svc.BusinessDao{
		PK: &pk,
		SK: &pk,
		Business: &svc.Business{
			Type:       aws.String("STORE"),
			BusinessId: req.BusinessId,
			Name:       req.Name,
			LegalName:  req.LegalName,
			Email:      req.Email,
			Address1:   req.Address1,
			Address2:   req.Address2,
			City:       req.City,
			State:      req.State,
			PostalCode: req.PostalCode,
			Country:    req.Country,
			Currency:   req.Currency,
			Phone:      req.Phone,
			Locale:     req.Locale,
			CreatedBy:  req.CreatedBy,
			CreatedAt:  req.CreatedAt,
			CustomAttribute: &map[string]interface{}{
				"GST": req.Gst,
				"PAN": req.Pan,
			},
		},
	}
	val, _ := dynamodbattribute.Marshal(dbReq)

	//userReq := &EmployeeDao{
	//	PK: &pk,
	//	SK: &sk,
	//	//GPK1: &gpk1,
	//	//GSK1: &gsk1,
	//	Employee: &Employee{
	//		Type:      aws.String("STR_EMP"),
	//		JoinedAt:  req.CreatedAt,
	//		CreatedBy: req.CreatedBy,
	//		CreatedAt: req.CreatedAt,
	//	},
	//}
	//usrInp, _ := dynamodbattribute.Marshal(userReq)

	db := svc.GetDynamoDbClient()

	_, err := db.PutItem(&dynamodb.PutItemInput{
		Item:                   val.M,
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(CommonTable),
	})
	if err != nil {
		log.Println(err)
	}
	return req
}

func CreateBusinessHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	eventJson, _ := json.Marshal(event)
	log.Printf("EVENT: %s", eventJson)

	var req svc.CreateBusinessRequest
	err := json.Unmarshal([]byte(event.Body), &req)

	// Error Unmarshalling the data
	if err != nil {
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

	// Step 2: Validate Valid User

	// Step 4: Create New Business Id
	businessId, err := GetNewBusinessID()
	if err != nil {
		log.Println(err)
		msg := &svc.XPOSApiError{
			ErrorMessage: "Error Creating New Business ID",
		}

		res, _ := json.Marshal(msg)

		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Headers:         nil,
			Body:            string(res),
			IsBase64Encoded: false,
		}, nil
	}

	// Save the business in the DB
	dbReq := createNewBusinessInputRequest(businessId, &req, req.CreatedBy)
	dbResp := CreateBusinessInDb(dbReq)

	stId := strconv.Itoa(*businessId)

	// Associate Admin Role with this user.
	storeId := "STORE#" + stId
	empId := "EMP#" + *req.CreatedBy

	empRole := &svc.StoreEmployeeRoleDao{
		PK:   &empId,
		SK:   &storeId,
		GPK1: &storeId,
		GSK1: &empId,
		StoreEmployeeRole: &svc.StoreEmployeeRole{
			EmployeeId: req.CreatedBy,
			StoreId:    &stId,
			Roles:      nil,
			Locale:     aws.String("en_US"),
			JoinedAt:   aws.Time(time.Now()),
			CreatedBy:  req.CreatedBy,
			CreatedAt:  aws.Time(time.Now()),
		},
	}

	employee, err := empRepo.CreateNewEmployeeForStore(empRole)

	log.Println(employee)

	// Return data to the user
	data, _ := json.Marshal(dbResp)
	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         nil,
		Body:            string(data),
		IsBase64Encoded: false,
	}, nil
}
