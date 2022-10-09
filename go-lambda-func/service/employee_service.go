package src

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

type EmployeeRepository struct{}

func GetUserFromCognito(username string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	cc := GetCognitoClient()

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

func FindUserAttributeFromCognito(attribute []*cognitoidentityprovider.AttributeType, key string) *string {
	data := ""

	for _, v := range attribute {
		if *v.Name == key {
			data = *v.Value
			break
		}
	}

	return &data
}

func (e *EmployeeRepository) GetEmployeeId(username *string) (*Employee, error) {

	user, err := GetUserFromCognito(*username)

	if err != nil {
		return nil, &EmployeeError{
			CausedBy: "User not found.",
			Message:  "User not found.",
		}
	}

	emp := &Employee{
		Type:          nil,
		EmployeeId:    user.Username,
		FirstName:     FindUserAttributeFromCognito(user.UserAttributes, "given_name"),
		MiddleName:    FindUserAttributeFromCognito(user.UserAttributes, "middle_name"),
		LastName:      FindUserAttributeFromCognito(user.UserAttributes, "family_name"),
		Locale:        FindUserAttributeFromCognito(user.UserAttributes, "locale"),
		Email:         FindUserAttributeFromCognito(user.UserAttributes, "email"),
		Phone:         FindUserAttributeFromCognito(user.UserAttributes, "phone_number"),
		Gender:        FindUserAttributeFromCognito(user.UserAttributes, "gender"),
		Picture:       FindUserAttributeFromCognito(user.UserAttributes, "picture"),
		EmailVerified: nil,
		PhoneVerified: nil,
		Dob:           nil,
		JoinedAt:      user.UserCreateDate,
		CreatedBy:     nil,
		CreatedAt:     user.UserCreateDate,
	}

	return emp, nil
}

func (e *EmployeeRepository) GetAllEmployeeFromStore(storeId *string) (*[]Employee, error) {
	return nil, nil
}

func (e *EmployeeRepository) CreateEmployee(req *EmployeeDao) (*Employee, error) {

	usrInp, _ := dynamodbattribute.Marshal(req)

	db := GetDynamoDbClient()
	item, err := db.PutItem(&dynamodb.PutItemInput{
		ConditionExpression:    nil,
		ConditionalOperator:    nil,
		Item:                   usrInp.M,
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(CommonTable),
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Item: %v", item)
	return nil, nil
}

func (e *EmployeeRepository) CreateNewEmployeeForStore(req *StoreEmployeeRoleDao) (*StoreEmployeeRole, error) {

	usrInp, _ := dynamodbattribute.Marshal(req)

	db := GetDynamoDbClient()
	item, err := db.PutItem(&dynamodb.PutItemInput{
		ConditionExpression:    nil,
		ConditionalOperator:    nil,
		Item:                   usrInp.M,
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(CommonTable),
	})
	if err != nil {
		return nil, err
	}

	log.Println(item)

	return req.StoreEmployeeRole, nil
}

func (e *EmployeeRepository) UpdateEmployee(emp *EmployeeUpdateRequest, username *string) error {
	cc := GetCognitoClient()

	var userAttributeInput []*cognitoidentityprovider.AttributeType

	if emp.FirstName != nil {
		userAttributeInput = append(userAttributeInput, &cognitoidentityprovider.AttributeType{
			Name:  aws.String("given_name"),
			Value: emp.FirstName,
		})
	}

	if emp.MiddleName != nil {
		userAttributeInput = append(userAttributeInput, &cognitoidentityprovider.AttributeType{
			Name:  aws.String("middle_name"),
			Value: emp.MiddleName,
		})
	}

	if emp.LastName != nil {
		userAttributeInput = append(userAttributeInput, &cognitoidentityprovider.AttributeType{
			Name:  aws.String("family_name"),
			Value: emp.LastName,
		})
	}

	if emp.Locale != nil {
		userAttributeInput = append(userAttributeInput, &cognitoidentityprovider.AttributeType{
			Name:  aws.String("locale"),
			Value: emp.Locale,
		})
	}

	if emp.Email != nil {
		userAttributeInput = append(userAttributeInput, &cognitoidentityprovider.AttributeType{
			Name:  aws.String("email"),
			Value: emp.Email,
		})
	}

	if emp.Gender != nil {
		userAttributeInput = append(userAttributeInput, &cognitoidentityprovider.AttributeType{
			Name:  aws.String("gender"),
			Value: emp.Gender,
		})
	}

	if emp.Dob != nil {
		userAttributeInput = append(userAttributeInput, &cognitoidentityprovider.AttributeType{
			Name:  aws.String("birthdate"),
			Value: aws.String(emp.Dob.UTC().Format("2006-01-02")),
		})
	}

	_, err := cc.AdminUpdateUserAttributes(&cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserAttributes: userAttributeInput,
		Username:       username,
		UserPoolId:     aws.String(CognitoUserPool),
	})

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (e *EmployeeRepository) GetStoreAssignedToEmployee(employeeId *string) (*[]Business, error) {
	return nil, nil
}

func (e *EmployeeRepository) CreateEmployeeOnCognito(req *CreateStoreEmployeeRequest) (*cognitoidentityprovider.UserType, error) {
	cc := GetCognitoClient()

	user, err := cc.AdminCreateUser(&cognitoidentityprovider.AdminCreateUserInput{
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("phone_number"),
				Value: req.Phone,
			},
			{
				Name:  aws.String("phone_number_verified"),
				Value: aws.String("true"),
			},
		},
		UserPoolId: aws.String(CognitoUserPool),
		Username:   req.Phone,
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	return user.User, nil
}

func (e *EmployeeRepository) IsUserAlreadyExistForStore(storeId *string, userId *string) (*bool, error) {
	db := GetDynamoDbClient()

	res := true

	dbBatchExecuteRequest := &dynamodb.BatchExecuteStatementInput{
		ReturnConsumedCapacity: aws.String("TOTAL"),
		Statements: []*dynamodb.BatchStatementRequest{
			{
				ConsistentRead: nil,
				Parameters:     nil,
				Statement:      aws.String(fmt.Sprintf("SELECT PK, SK FROM \"%s\" WHERE PK = '%s' AND SK = '%s'", CommonTable, fmt.Sprintf("STORE#%s", *storeId), fmt.Sprintf("STORE#%s", *storeId))),
			},
			{
				ConsistentRead: nil,
				Parameters:     nil,
				Statement:      aws.String(fmt.Sprintf("SELECT PK, SK FROM \"%s\" WHERE PK = '%s' AND SK = '%s'", CommonTable, fmt.Sprintf("EMP#%s", *userId), fmt.Sprintf("STORE#%s", *storeId))),
			},
		},
	}

	item, err := db.BatchExecuteStatement(dbBatchExecuteRequest)

	log.Printf("Item: %v", item.ConsumedCapacity)

	if err != nil {
		return nil, err
	}

	if item.Responses[0].Item == nil {
		return nil, &EmployeeError{
			CausedBy: "Store Not Found",
			Message:  "Store Not Found",
		}
	}

	if item.Responses[1].Item != nil {
		return nil, &EmployeeError{
			CausedBy: "User Already Exist",
			Message:  "User Already Exist",
		}
	}

	return &res, nil
}
