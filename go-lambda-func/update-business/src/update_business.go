package src

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	svc "service"
)

func UpdateBusiness(businessId string, req *svc.CreateBusinessRequest) *svc.StoreError {
	db := svc.GetDynamoDbClient()

	pk := "STORE#" + businessId

	updateExp := "SET "
	expressionAttributeNames := make(map[string]*string)
	expressionAttributeValues := make(map[string]*dynamodb.AttributeValue)

	if req.Name != nil {
		updateExp += "#name = :name, "
		expressionAttributeNames["#name"] = aws.String("name")
		expressionAttributeValues[":name"] = &dynamodb.AttributeValue{S: aws.String(*req.Name)}
	}

	if req.Email != nil {
		updateExp += "#email = :email, "
		expressionAttributeNames["#email"] = aws.String("email")
		expressionAttributeValues[":email"] = &dynamodb.AttributeValue{S: aws.String(*req.Email)}
	}

	if req.Address1 != nil {
		updateExp += "#address1 = :address1, "
		expressionAttributeNames["#address1"] = aws.String("address1")
		expressionAttributeValues[":address1"] = &dynamodb.AttributeValue{S: aws.String(*req.Address1)}
	}

	if req.Address2 != nil {
		updateExp += "#address2 = :address2, "
		expressionAttributeNames["#address2"] = aws.String("address2")
		expressionAttributeValues[":address2"] = &dynamodb.AttributeValue{S: aws.String(*req.Address2)}
	}

	if req.City != nil {
		updateExp += "#city = :city, "
		expressionAttributeNames["#city"] = aws.String("city")
		expressionAttributeValues[":city"] = &dynamodb.AttributeValue{S: aws.String(*req.City)}
	}

	if req.State != nil {
		updateExp += "#state = :state, "
		expressionAttributeNames["#state"] = aws.String("state")
		expressionAttributeValues[":state"] = &dynamodb.AttributeValue{S: aws.String(*req.State)}
	}

	if req.PostalCode != nil {
		updateExp += "#postalCode = :postalCode, "
		expressionAttributeNames["#postalCode"] = aws.String("postalCode")
		expressionAttributeValues[":postalCode"] = &dynamodb.AttributeValue{S: aws.String(*req.PostalCode)}
	}

	if req.Country != nil {
		updateExp += "#country = :country, "
		expressionAttributeNames["#country"] = aws.String("country")
		expressionAttributeValues[":country"] = &dynamodb.AttributeValue{S: aws.String(*req.Country)}
	}

	if req.Phone != nil {
		updateExp += "#phone = :phone, "
		expressionAttributeNames["#phone"] = aws.String("phone")
		expressionAttributeValues[":phone"] = &dynamodb.AttributeValue{S: aws.String(*req.Phone)}
	}

	if req.Locale != nil {
		updateExp += "#locale = :locale, "
		expressionAttributeNames["#locale"] = aws.String("locale")
		expressionAttributeValues[":locale"] = &dynamodb.AttributeValue{S: aws.String(*req.Locale)}
	}

	if req.Gst != nil {
		updateExp += "custom_attribute.#gst = :gst, "
		expressionAttributeNames["#gst"] = aws.String("GST")
		expressionAttributeValues[":gst"] = &dynamodb.AttributeValue{S: aws.String(*req.Gst)}
	}

	if req.Pan != nil {
		updateExp += "custom_attribute.#pan = :pan, "
		expressionAttributeNames["#pan"] = aws.String("PAN")
		expressionAttributeValues[":pan"] = &dynamodb.AttributeValue{S: aws.String(*req.Pan)}
	}

	if len(updateExp) <= 4 {
		// Throw error
		return &svc.StoreError{
			ApiError: &svc.ApiError{
				CausedBy: "UpdateBusiness",
				Message:  "No data to update",
			},
		}
	}

	updateExp = updateExp[:len(updateExp)-2]

	tmp := dynamodb.UpdateItemInput{
		TableName: aws.String("XPOS_DEV"),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(pk),
			},
			"SK": {
				S: aws.String(pk),
			},
		},
		ConditionExpression:         aws.String("attribute_exists(PK)"),
		UpdateExpression:            &updateExp,
		ExpressionAttributeNames:    expressionAttributeNames,
		ExpressionAttributeValues:   expressionAttributeValues,
		ReturnConsumedCapacity:      nil,
		ReturnItemCollectionMetrics: nil,
		ReturnValues:                nil,
	}

	_, err := db.UpdateItem(&tmp)

	if err != nil {

		switch e := err.(type) {
		case *dynamodb.ConditionalCheckFailedException:
			// Do something with the path
			return &svc.StoreError{
				ApiError: &svc.ApiError{
					CausedBy: "UpdateBusiness",
					Message:  "Business not found",
				},
			}
		default:
			log.Println(e)
		}

		return &svc.StoreError{
			ApiError: &svc.ApiError{
				CausedBy: "UpdateBusiness",
				Message:  err.Error(),
			},
		}
	}

	return nil
}

func UpdateBusinessHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	eventJson, _ := json.Marshal(event)
	log.Printf("EVENT: %s", eventJson)

	pathParams := event.PathParameters
	bId := pathParams["businessId"]

	var req svc.CreateBusinessRequest
	err := json.Unmarshal([]byte(event.Body), &req)

	// Error Unmarshalling the data
	if err != nil {
		msg := &svc.XPOSApiError{
			ErrorMessage: "Invalid Input Data",
		}

		res, _ := json.Marshal(msg)

		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Headers:         nil,
			Body:            string(res),
			IsBase64Encoded: false,
		}, nil
	}

	er := UpdateBusiness(bId, &req)

	if er != nil {
		res, _ := json.Marshal(er)

		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Headers:         nil,
			Body:            string(res),
			IsBase64Encoded: false,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         nil,
		IsBase64Encoded: false,
	}, nil
}
