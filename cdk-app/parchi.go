package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	awscdkrest "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"

	awscdklambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ParchiStackProps struct {
	awscdk.StackProps
}

func NewParchiStack(scope constructs.Construct, id string, props *ParchiStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	role := awsiam.Role_FromRoleName(stack, jsii.String("Role"), jsii.String("thelawala-admin-dev"), &awsiam.FromRoleNameOptions{
		Mutable: jsii.Bool(false),
	})

	//role := awsiam.Role_FromRoleArn(stack, jsii.String("ParchiRole"), jsii.String("arn:aws:iam::${AWS::AccountId}:role/thelawala-admin-dev"), &awsiam.FromRoleArnOptions{
	//	Mutable: jsii.Bool(false),
	//})

	getBusinessByIdLambda := awscdklambdago.NewGoFunction(stack, jsii.String("DevXPOSGetBusinessById"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String("DevXPOSGetBusinessById"),
		Description:  jsii.String("Get Business Detail By ID"),
		Entry:        jsii.String("../go-lambda-func/get-business-by-id"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  &map[string]*string{},
	})

	createNewBusinessLambda := awscdklambdago.NewGoFunction(stack, jsii.String("DevXPOSCreateNewBusinessLambda"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String("DevXPOSCreateNewBusinessLambda"),
		Description:  jsii.String("Create New Business for the application"),
		Entry:        jsii.String("../go-lambda-func/create-new-business"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  &map[string]*string{},
	})

	updateBusinessLambda := awscdklambdago.NewGoFunction(stack, jsii.String("DevXPOSUpdateBusinessLambda"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String("DevXPOSUpdateBusinessLambda"),
		Description:  jsii.String("Update Business Detail."),
		Entry:        jsii.String("../go-lambda-func/update-business"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  &map[string]*string{},
	})

	restApi := awscdkrest.NewRestApi(stack, jsii.String("XPOS-POC"), &awscdkrest.RestApiProps{
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

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewParchiStack(app, "ParchiStack", &ParchiStackProps{
		awscdk.StackProps{
			Env: env(),
		},
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
