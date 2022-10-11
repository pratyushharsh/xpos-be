package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	awscdkrest "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awscdklambdanodejs "github.com/aws/aws-cdk-go/awscdk/v2/awslambdanodejs"
	awscdklambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const ProjectName = "XPOS"       // USE 3-5 Digit Capital Letters
const ProjectEnvironment = "DEV" // USE Capital letters
const ProjectVersion = "0.0.1"
const ProjectDescription = "XPOS"

const ProjectPrefix = ProjectName + "-" + ProjectEnvironment + "-"
const StackPrefix = ProjectName + "-"

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	//st := NewParchiStack(app, "ParchiStack", &ParchiStackProps{
	//	awscdk.StackProps{
	//		Env: env(),
	//	},
	//})

	globalRole := jsii.String("thelawala-admin-dev")

	NewUserStack(app, StackPrefix+"UserStack", &UserStackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		RoleName: globalRole,
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}

type UserStackProps struct {
	awscdk.StackProps
	RoleName *string
}

func NewUserStack(scope constructs.Construct, id string, props *UserStackProps) awscdk.Stack {
	var stackProps awscdk.StackProps
	if props != nil {
		stackProps = props.StackProps
	}
	userStack := awscdk.NewStack(scope, &id, &stackProps)

	role := awsiam.Role_FromRoleName(userStack, jsii.String("Role"), props.RoleName, &awsiam.FromRoleNameOptions{
		Mutable: jsii.Bool(false),
	})
	/* Lambda Function to handle the custom authentication for
	   OTP Login using aws
	*/
	preSignUpLambda := awscdklambdanodejs.NewNodejsFunction(userStack, jsii.String(ProjectPrefix+"PreSignUp"), &awscdklambdanodejs.NodejsFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "PreSignUp"),
		Entry:        jsii.String("../nodejs-lambda-func/auth/pre-signup/index.ts"),
		Handler:      jsii.String("handler"),
		Runtime:      awslambda.Runtime_NODEJS_16_X(),
		MemorySize:   jsii.Number(128),
		Role:         role,
	})

	createAuthChallangeLambda := awscdklambdanodejs.NewNodejsFunction(userStack, jsii.String(ProjectPrefix+"CreateAuthChallenge"), &awscdklambdanodejs.NodejsFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "CreateAuthChallenge"),
		Entry:        jsii.String("../nodejs-lambda-func/auth/create-auth-challenge/index.ts"),
		Handler:      jsii.String("handler"),
		Runtime:      awslambda.Runtime_NODEJS_16_X(),
		MemorySize:   jsii.Number(128),
		Role:         role,
	})

	defineAuthChallange := awscdklambdanodejs.NewNodejsFunction(userStack, jsii.String(ProjectPrefix+"DefineAuthChallenge"), &awscdklambdanodejs.NodejsFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "DefineAuthChallenge"),
		Entry:        jsii.String("../nodejs-lambda-func/auth/define-auth-challenge/index.ts"),
		Handler:      jsii.String("handler"),
		Runtime:      awslambda.Runtime_NODEJS_16_X(),
		MemorySize:   jsii.Number(128),
		Role:         role,
	})

	verifyAuthChallenge := awscdklambdanodejs.NewNodejsFunction(userStack, jsii.String(ProjectPrefix+"VerifyAuthChallenge"), &awscdklambdanodejs.NodejsFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "VerifyAuthChallenge"),
		Entry:        jsii.String("../nodejs-lambda-func/auth/verify-auth-challenge/index.ts"),
		Handler:      jsii.String("handler"),
		Runtime:      awslambda.Runtime_NODEJS_16_X(),
		MemorySize:   jsii.Number(128),
		Role:         role,
	})

	userPool := awscognito.NewUserPool(userStack, jsii.String(ProjectPrefix+"UserPool"), &awscognito.UserPoolProps{
		UserPoolName:      jsii.String(ProjectPrefix + "UserPool"),
		SelfSignUpEnabled: jsii.Bool(true),
		SignInAliases: &awscognito.SignInAliases{
			Phone: jsii.Bool(true),
			Email: jsii.Bool(false),
		},
		StandardAttributes: &awscognito.StandardAttributes{
			PhoneNumber: &awscognito.StandardAttribute{
				Required: jsii.Bool(true),
				Mutable:  jsii.Bool(false),
			},
		},
	})

	userPool.AddTrigger(awscognito.UserPoolOperation_PRE_SIGN_UP(), preSignUpLambda)
	userPool.AddTrigger(awscognito.UserPoolOperation_CREATE_AUTH_CHALLENGE(), createAuthChallangeLambda)
	userPool.AddTrigger(awscognito.UserPoolOperation_DEFINE_AUTH_CHALLENGE(), defineAuthChallange)
	userPool.AddTrigger(awscognito.UserPoolOperation_VERIFY_AUTH_CHALLENGE_RESPONSE(), verifyAuthChallenge)

	return userStack
}

type BusinessStackProps struct {
	awscdk.StackProps
	Role awsiam.IRole
}

func NewBusinessStack(scope constructs.Construct, id string, props *BusinessStackProps) awscdk.Stack {
	var stackProps awscdk.StackProps
	if props != nil {
		stackProps = props.StackProps
	}
	businessStack := awscdk.NewStack(scope, &id, &stackProps)

	getBusinessByIdLambda := awscdklambdago.NewGoFunction(businessStack, jsii.String("DevXPOSGetBusinessById"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String("DevXPOSGetBusinessById"),
		Description:  jsii.String("Get Business Detail By ID"),
		Entry:        jsii.String("../go-lambda-func/get-business-by-id"),
		Role:         props.Role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  &map[string]*string{},
	})

	createNewBusinessLambda := awscdklambdago.NewGoFunction(businessStack, jsii.String("DevXPOSCreateNewBusinessLambda"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String("DevXPOSCreateNewBusinessLambda"),
		Description:  jsii.String("Create New Business for the application"),
		Entry:        jsii.String("../go-lambda-func/create-new-business"),
		Role:         props.Role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  &map[string]*string{},
	})

	updateBusinessLambda := awscdklambdago.NewGoFunction(businessStack, jsii.String("DevXPOSUpdateBusinessLambda"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String("DevXPOSUpdateBusinessLambda"),
		Description:  jsii.String("Update Business Detail."),
		Entry:        jsii.String("../go-lambda-func/update-business"),
		Role:         props.Role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  &map[string]*string{},
	})

	// @TODO If API ID is present then use the existing api or create a new api.
	restApi := awscdkrest.NewRestApi(businessStack, jsii.String("XPOS-POC"), &awscdkrest.RestApiProps{
		Deploy: jsii.Bool(false),
	})

	businessPath := restApi.Root().AddResource(jsii.String("business"), &awscdkrest.ResourceOptions{})

	// Handler for: POST /business/
	businessPath.AddMethod(jsii.String("POST"), awscdkrest.NewLambdaIntegration(createNewBusinessLambda, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	businessId := businessPath.AddResource(jsii.String("{businessId}"), &awscdkrest.ResourceOptions{
		DefaultIntegration:   nil,
		DefaultMethodOptions: nil,
	})

	// Handler for: GET /business/{businessId}
	businessId.AddMethod(jsii.String("GET"), awscdkrest.NewLambdaIntegration(getBusinessByIdLambda, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	// Handler for updating business detail: PUT /business/{businessId}
	businessId.AddMethod(jsii.String("PUT"), awscdkrest.NewLambdaIntegration(updateBusinessLambda, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	return businessStack
}
