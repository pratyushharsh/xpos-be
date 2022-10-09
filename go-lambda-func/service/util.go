package src

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var (
	sess          *session.Session
	lambdaClient  *lambda.Lambda
	cognitoClient *cognitoidentityprovider.CognitoIdentityProvider
	dynamoClient  *dynamodb.DynamoDB
)

func GetSession() *session.Session {
	if sess == nil {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}
	return sess
}

func GetLambdaClient() *lambda.Lambda {
	if lambdaClient == nil {
		lambdaClient = lambda.New(GetSession(), &aws.Config{Region: aws.String("ap-south-1")})
	}
	return lambdaClient
}

func GetCognitoClient() *cognitoidentityprovider.CognitoIdentityProvider {
	if cognitoClient == nil {
		cognitoClient = cognitoidentityprovider.New(GetSession())
	}
	return cognitoClient
}

func GetDynamoDbClient() *dynamodb.DynamoDB {
	if dynamoClient == nil {
		dynamoClient = dynamodb.New(GetSession())
	}
	return dynamoClient
}

func UnmarshalStreamImage(attribute map[string]events.DynamoDBAttributeValue, out interface{}) error {

	dbAttrMap := make(map[string]*dynamodb.AttributeValue)

	for k, v := range attribute {

		var dbAttr dynamodb.AttributeValue

		bytes, marshalErr := v.MarshalJSON()
		if marshalErr != nil {
			return marshalErr
		}

		json.Unmarshal(bytes, &dbAttr)
		dbAttrMap[k] = &dbAttr
	}

	return dynamodbattribute.UnmarshalMap(dbAttrMap, out)

}
