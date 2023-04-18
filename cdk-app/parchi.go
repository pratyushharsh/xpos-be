package main

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	awscdkrest "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awscdklambdanodejs "github.com/aws/aws-cdk-go/awscdk/v2/awslambdanodejs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	awscdklambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const ProjectName = "XPOS"       // USE 3-5 Digit Capital Letters
const ProjectEnvironment = "DEV" // USE Capital letters
const ProjectVersion = "0.0.1"
const ProjectDescription = "XPOS"
const ProjectTag = ProjectName + "-" + ProjectEnvironment

const ProjectPrefix = ProjectName + "-" + ProjectEnvironment + "-"
const StackPrefix = ProjectName + "-"

func mergeTwoMaps(a *map[string]*string, b *map[string]*string) *map[string]*string {
	merged := make(map[string]*string)
	for k, v := range *a {
		merged[k] = v
	}
	for k, v := range *b {
		merged[k] = v
	}
	return &merged
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	//st := NewParchiStack(app, "ParchiStack", &ParchiStackProps{
	//	awscdk.StackProps{
	//		Env: env(),
	//	},
	//})

	globalRole := jsii.String("thelawala-admin-dev")

	userStack := NewUserStack(app, StackPrefix+"UserStack", &UserStackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
			Tags: &map[string]*string{
				"Project": jsii.String(ProjectTag),
			},
		},
		RoleName: globalRole,
	})

	fmt.Printf("%s", userStack)

	NewBusinessStack(app, StackPrefix+"BusinessStack", &BusinessStackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
			Tags: &map[string]*string{
				"Project": jsii.String(ProjectTag),
			},
		},
		RoleName: globalRole,
	})

	//NewImageStack(app, StackPrefix+"ImageStack", &ImageStackProps{
	//	StackProps: awscdk.StackProps{
	//		Env: env(),
	//		Tags: &map[string]*string{
	//			"Project": jsii.String(ProjectTag),
	//		},
	//	},
	//	RoleName: globalRole,
	//})

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
	//
	//preTokenGenerationLambda := awscdklambdanodejs.NewNodejsFunction(userStack, jsii.String(ProjectPrefix+"PreTokenGeneration"), &awscdklambdanodejs.NodejsFunctionProps{
	//	FunctionName: jsii.String(ProjectPrefix + "PreTokenGeneration"),
	//	Entry:        jsii.String("../nodejs-lambda-func/auth/pre-token-generation-trigger/index.ts"),
	//	Handler:      jsii.String("handler"),
	//	Runtime:      awslambda.Runtime_NODEJS_16_X(),
	//	MemorySize:   jsii.Number(128),
	//	Role:         role,
	//})

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
		PasswordPolicy: &awscognito.PasswordPolicy{
			RequireDigits:    jsii.Bool(false),
			RequireLowercase: jsii.Bool(false),
			RequireSymbols:   jsii.Bool(false),
			RequireUppercase: jsii.Bool(false),
			MinLength:        jsii.Number(8),
		},
		DeviceTracking: &awscognito.DeviceTracking{
			DeviceOnlyRememberedOnUserPrompt: jsii.Bool(false),
			ChallengeRequiredOnNewDevice:     jsii.Bool(false),
		},
	})


	userPool.AddTrigger(awscognito.UserPoolOperation_PRE_SIGN_UP(), preSignUpLambda)
	userPool.AddTrigger(awscognito.UserPoolOperation_CREATE_AUTH_CHALLENGE(), createAuthChallangeLambda)
	userPool.AddTrigger(awscognito.UserPoolOperation_DEFINE_AUTH_CHALLENGE(), defineAuthChallange)
	userPool.AddTrigger(awscognito.UserPoolOperation_VERIFY_AUTH_CHALLENGE_RESPONSE(), verifyAuthChallenge)
	//userPool.AddTrigger(awscognito.UserPoolOperation_PRE_TOKEN_GENERATION(), preTokenGenerationLambda)

	resourceScope := awscognito.NewResourceServerScope(&awscognito.ResourceServerScopeProps{
		ScopeName: jsii.String("xpos"),
		ScopeDescription: jsii.String("It will have access to all the api related to xpos"),
	})

	resourceServer := awscognito.NewUserPoolResourceServer(userStack, jsii.String(ProjectPrefix+"ResourceServer"), &awscognito.UserPoolResourceServerProps{
		Identifier:                 jsii.String(ProjectPrefix+"AppClientResourceServer"),
		Scopes:                     &[]awscognito.ResourceServerScope{
			resourceScope,
		},
		UserPoolResourceServerName: jsii.String("xpos-resource-server"),
		UserPool:                   userPool,
	})

	userPool.AddClient(jsii.String(ProjectPrefix+"MobileApp"), &awscognito.UserPoolClientOptions{
		GenerateSecret:     jsii.Bool(false),
		UserPoolClientName: jsii.String(ProjectPrefix + "MobileApp"),
		OAuth: &awscognito.OAuthSettings{
			Scopes: &[]awscognito.OAuthScope{
				awscognito.OAuthScope_ResourceServer(resourceServer, resourceScope),
			},
		},
	})

	userPool.AddClient(jsii.String(ProjectPrefix+"IOSMobileApp"), &awscognito.UserPoolClientOptions{
		GenerateSecret:     jsii.Bool(false),
		UserPoolClientName: jsii.String(ProjectPrefix + "IOSMobileApp"),
		OAuth: &awscognito.OAuthSettings{
			Scopes: &[]awscognito.OAuthScope{
				awscognito.OAuthScope_ResourceServer(resourceServer, resourceScope),
			},
		},
	})

	userStack.ExportValue(userPool.UserPoolId(), &awscdk.ExportValueOptions{
		Name: jsii.String(ProjectPrefix + "UserPool"),
	})

	return userStack
}

type BusinessStackProps struct {
	awscdk.StackProps
	RoleName *string
}

func NewBusinessStack(scope constructs.Construct, id string, props *BusinessStackProps) awscdk.Stack {
	var stackProps awscdk.StackProps
	if props != nil {
		stackProps = props.StackProps
	}
	businessStack := awscdk.NewStack(scope, &id, &stackProps)

	// Role that needs to be attached to the lambda function
	role := awsiam.Role_FromRoleName(businessStack, jsii.String("Role"), props.RoleName, &awsiam.FromRoleNameOptions{
		Mutable: jsii.Bool(false),
	})

	// @TODO If API ID is present then use the existing api or create a new api.
	restApi := awscdkrest.NewRestApi(businessStack, jsii.String("XPOS-POC"), &awscdkrest.RestApiProps{
		Deploy: jsii.Bool(false),
	})

	authorizer := awscdkrest.NewCognitoUserPoolsAuthorizer(businessStack, jsii.String("XPOS-POC-Authorizer"), &awscdkrest.CognitoUserPoolsAuthorizerProps{
		IdentitySource: jsii.String("method.request.header.Authorizer"),
		CognitoUserPools: &[]awscognito.IUserPool{
			awscognito.UserPool_FromUserPoolId(businessStack, jsii.String("XPOS-POC-UserPool-Authorizer"), awscdk.Fn_ImportValue(jsii.String(ProjectPrefix+"UserPool"))),
		},
	})

	//lambdaLayers := awslambda.LayerVersion_FromLayerVersionArn(businessStack, jsii.String(ProjectPrefix+"DependencyLayer"), jsii.String("arn:aws:lambda:ap-south-1:189468856814:layer:XPOSDependencyLayer:3"))
	//lambdaLayers := awslambda.NewLayerVersion(businessStack, jsii.String("MyLayer"), &awslambda.LayerVersionProps{
	//  	RemovalPolicy: awscdk.RemovalPolicy_RETAIN,
	//  	Code: awslambda.AssetCode,
	//  	compatibleArchitectures: []architecture{
	//  		lambda.*architecture_X86_64(),
	//  		lambda.*architecture_ARM_64(),
	//  	},
	//  })

	lambdaEnvironmentVariable := &map[string]*string{
		"DBTable":   jsii.String("XPOS_DEV"),
		"DataTable": jsii.String("XPOS_DATA"),
		//"CognitoUserPool": jsii.String("ap-south-1_gXgaeT7lu"),
		"CognitoUserPool": awscdk.Fn_ImportValue(jsii.String(ProjectPrefix + "UserPool")),
	}

	getBusinessByIdLambda := awscdklambdago.NewGoFunction(businessStack, jsii.String(ProjectPrefix+"GetBusinessById"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "GetBusinessById"),
		Description:  jsii.String("Get Business Detail By ID"),
		Entry:        jsii.String("../go-lambda-func/get-business-by-id"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  lambdaEnvironmentVariable,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
	})

	createNewBusinessLambda := awscdklambdago.NewGoFunction(businessStack, jsii.String(ProjectPrefix+"CreateNewBusinessLambda"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "CreateNewBusinessLambda"),
		Description:  jsii.String("Create New Business for the application"),
		Entry:        jsii.String("../go-lambda-func/create-new-business"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  lambdaEnvironmentVariable,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
	})

	updateBusinessLambda := awscdklambdago.NewGoFunction(businessStack, jsii.String(ProjectPrefix+"UpdateBusinessLambda"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "UpdateBusinessLambda"),
		Description:  jsii.String("Update Business Detail."),
		Entry:        jsii.String("../go-lambda-func/update-business"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  lambdaEnvironmentVariable,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
	})

	businessPath := restApi.Root().AddResource(jsii.String("business"), &awscdkrest.ResourceOptions{})

	storePath := restApi.Root().AddResource(jsii.String("store"), &awscdkrest.ResourceOptions{})

	// Handler for: POST /business/
	businessPath.AddMethod(jsii.String("POST"), awscdkrest.NewLambdaIntegration(createNewBusinessLambda, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	businessId := businessPath.AddResource(jsii.String("{businessId}"), &awscdkrest.ResourceOptions{
		DefaultIntegration:   nil,
		DefaultMethodOptions: nil,
	})

	storeId := storePath.AddResource(jsii.String("{businessId}"), &awscdkrest.ResourceOptions{
		DefaultIntegration:   nil,
		DefaultMethodOptions: nil,
	})

	// Handler for: GET /business/{businessId}
	businessId.AddMethod(jsii.String("GET"), awscdkrest.NewLambdaIntegration(getBusinessByIdLambda, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{
		Authorizer: authorizer,
	})

	// Handler for: GET /store/{businessId}
	storeId.AddMethod(jsii.String("GET"), awscdkrest.NewLambdaIntegration(getBusinessByIdLambda, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	// Handler for updating business detail: PUT /business/{businessId}
	businessId.AddMethod(jsii.String("PUT"), awscdkrest.NewLambdaIntegration(updateBusinessLambda, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	// Create signature for uploading the image url.
	businessImage := businessId.AddResource(jsii.String("image"), &awscdkrest.ResourceOptions{
		DefaultIntegration:   nil,
		DefaultMethodOptions: nil,
	})

	businessImageToken := businessImage.AddResource(jsii.String("token"), &awscdkrest.ResourceOptions{
		DefaultIntegration:   nil,
		DefaultMethodOptions: nil,
	})

	businessImageTokenLambda := awscdklambdago.NewGoFunction(businessStack, jsii.String(ProjectPrefix+"BusinessImageTokenLambda"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "BusinessImageTokenLambda"),
		Description:  jsii.String("Generate Token to upload image."),
		Entry:        jsii.String("../go-lambda-func/imagekit-token"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  lambdaEnvironmentVariable,
		LogRetention: awslogs.RetentionDays_THREE_DAYS,
	})

	// Get token to upload image url: GET /business/{businessId}/image/token
	businessImageToken.AddMethod(jsii.String("GET"), awscdkrest.NewLambdaIntegration(businessImageTokenLambda, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	// Create URL For deployment of Business Logo Upload
	businessLogoUploadUrlLambda := awscdklambdanodejs.NewNodejsFunction(businessStack, jsii.String(ProjectPrefix+"BusinessLogoUploadUrl"), &awscdklambdanodejs.NodejsFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "BusinessLogoUploadUrl"),
		Entry:        jsii.String("../nodejs-lambda-func/image-processing/business-upload-url/index.ts"),
		Handler:      jsii.String("handler"),
		Runtime:      awslambda.Runtime_NODEJS_16_X(),
		MemorySize:   jsii.Number(128),
		Role:         role,
		Environment: mergeTwoMaps(lambdaEnvironmentVariable, &map[string]*string{
			"URL_EXPIRATION_SECONDS": jsii.String("300"),
			"IMAGE_IMAGE_BUCKET":     jsii.String("xpos-image-stage"),
			"BUCKET_PREFIX":          jsii.String("parchi"),
		}),
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
	})

	// Create Sync Service
	syncService := businessId.AddResource(jsii.String("sync"), &awscdkrest.ResourceOptions{
		DefaultIntegration:   nil,
		DefaultMethodOptions: nil,
	})

	// Handler for: POST /business/{businessId}/sync
	updateSyncLambda := awscdklambdago.NewGoFunction(businessStack, jsii.String(ProjectPrefix+"UpdateSyncService"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "UpdateSyncService"),
		Description:  jsii.String("Sync Data from local to cloud."),
		Entry:        jsii.String("../go-lambda-func/update-sync-data"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  lambdaEnvironmentVariable,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
	})

	syncService.AddMethod(jsii.String("POST"), awscdkrest.NewLambdaIntegration(updateSyncLambda, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	// Handler for: GET /business/{businessId}/sync
	getSyncLambda := awscdklambdago.NewGoFunction(businessStack, jsii.String(ProjectPrefix+"GetSyncService"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "GetSyncService"),
		Description:  jsii.String("Get Sync Data from cloud."),
		Entry:        jsii.String("../go-lambda-func/get-sync-data"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  lambdaEnvironmentVariable,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
	})

	syncService.AddMethod(jsii.String("GET"), awscdkrest.NewLambdaIntegration(getSyncLambda, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	// Settings For a store
	settings := businessId.AddResource(jsii.String("settings"), &awscdkrest.ResourceOptions{
		DefaultIntegration:   nil,
		DefaultMethodOptions: nil,
	})

	tax := settings.AddResource(jsii.String("tax"), &awscdkrest.ResourceOptions{
		DefaultIntegration:   nil,
		DefaultMethodOptions: nil,
	})

	// Handler for: GET /business/{businessId}/settings/tax
	getTaxForGroupForStoreLambda := awscdklambdago.NewGoFunction(businessStack, jsii.String(ProjectPrefix+"GetTaxGroupForStore"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "GetTaxGroupForStore"),
		Description:  jsii.String("Get Tax Group For Store"),
		Entry:        jsii.String("../go-lambda-func/get-tax-group-for-store"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  lambdaEnvironmentVariable,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
	})

	// Handler for: POST /business/{businessId}/settings/tax
	createTaxForGroupForStoreLambda := awscdklambdago.NewGoFunction(businessStack, jsii.String(ProjectPrefix+"CreateTaxGroupForStore"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "CreateTaxGroupForStore"),
		Description:  jsii.String("Create Tax Group For Store"),
		Entry:        jsii.String("../go-lambda-func/create-tax-group"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  lambdaEnvironmentVariable,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
	})

	tax.AddMethod(jsii.String("GET"), awscdkrest.NewLambdaIntegration(getTaxForGroupForStoreLambda, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	tax.AddMethod(jsii.String("POST"), awscdkrest.NewLambdaIntegration(createTaxForGroupForStoreLambda, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	logoApi := businessId.AddResource(jsii.String("logo"), &awscdkrest.ResourceOptions{
		DefaultIntegration:   nil,
		DefaultMethodOptions: nil,
	})

	// Handler for: GET /business/{businessId}/logo
	logoApi.AddMethod(jsii.String("GET"), awscdkrest.NewLambdaIntegration(businessLogoUploadUrlLambda, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	// Lambda for compressing the image and store in the S3 bucket
	//compressImageLambda :=
	//awscdklambdanodejs.NewNodejsFunction(businessStack, jsii.String(ProjectPrefix+"CompressImageLambda"), &awscdklambdanodejs.NodejsFunctionProps{
	//	FunctionName: jsii.String(ProjectPrefix + "CompressImageLambda"),
	//	Entry:        jsii.String("../nodejs-lambda-func/image-processing/compress-image/compressImageHandler.js"),
	//	Handler:      jsii.String("main"),
	//	Runtime:      awslambda.Runtime_NODEJS_16_X(),
	//	MemorySize:   jsii.Number(256),
	//	Role:         role,
	//	Environment: mergeTwoMaps(lambdaEnvironmentVariable, &map[string]*string{
	//		"INPUT_BUCKET":  jsii.String("xpos-image-stage"),
	//		"OUTPUT_BUCKET": jsii.String("xpos-image"),
	//	}),
	//	Bundling: &awscdklambdanodejs.BundlingOptions{
	//		ExternalModules: &[]*string{
	//			jsii.String("aws-sdk"),
	//			jsii.String("sharp"),
	//		},
	//	},
	//	Layers: &[]awslambda.ILayerVersion{
	//		lambdaLayers,
	//	},
	//	LogRetention: awslogs.RetentionDays_ONE_WEEK,
	//})

	//inputImageBucket := awss3.Bucket_FromBucketName(businessStack, jsii.String("xpos-image-stage"), jsii.String("xpos-image-stage"))
	//////
	////// Create listener for the image upload notification
	//inputImageBucket.AddEventNotification(awss3.EventType_OBJECT_CREATED_PUT, awss3notifications.NewLambdaDestination(compressImageLambda), &awss3.NotificationKeyFilter{
	//	Prefix: jsii.String("parchi/"),
	//})

	// User Stack Creating new business and new user
	getEmployeeFromBusiness := awscdklambdago.NewGoFunction(businessStack, jsii.String(ProjectPrefix+"GetEmployeeFromBusiness"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "GetEmployeeFromBusiness"),
		Description:  jsii.String("Get Employees from business."),
		Entry:        jsii.String("../go-lambda-func/get-employee-from-business"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  lambdaEnvironmentVariable,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
	})

	businessEmpApi := businessId.AddResource(jsii.String("employee"), &awscdkrest.ResourceOptions{
		DefaultIntegration:   nil,
		DefaultMethodOptions: nil,
	})

	// Handler for: GET /business/{businessId}/employee
	businessEmpApi.AddMethod(jsii.String("GET"), awscdkrest.NewLambdaIntegration(getEmployeeFromBusiness, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	// Create new employee for the store.
	createCreateStoreEmployeeLambda := awscdklambdago.NewGoFunction(businessStack, jsii.String(ProjectPrefix+"CreateStoreEmployee"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "CreateStoreEmployee"),
		Description:  jsii.String("Get Employees from business."),
		Entry:        jsii.String("../go-lambda-func/create-new-employee"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  lambdaEnvironmentVariable,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
	})

	// Handler for: POST /business/{businessId}/employee
	businessEmpApi.AddMethod(jsii.String("POST"), awscdkrest.NewLambdaIntegration(createCreateStoreEmployeeLambda, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	// Get Employee from store by id
	getEmployeeFromBusinessById := awscdklambdago.NewGoFunction(businessStack, jsii.String(ProjectPrefix+"GetEmployeeFromBusinessById"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "GetEmployeeFromBusinessById"),
		Description:  jsii.String("Get Employees from business."),
		Entry:        jsii.String("../go-lambda-func/get-employee-from-business-by-id"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  lambdaEnvironmentVariable,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
	})

	businessUserApi := businessEmpApi.AddResource(jsii.String("{userid}"), &awscdkrest.ResourceOptions{
		DefaultIntegration:   nil,
		DefaultMethodOptions: nil,
	})

	// Handler for: GET /business/{businessId}/employee/{userid}
	businessUserApi.AddMethod(jsii.String("GET"), awscdkrest.NewLambdaIntegration(getEmployeeFromBusinessById, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	userPath := restApi.Root().AddResource(jsii.String("user"), &awscdkrest.ResourceOptions{})

	// Get Employee from store by id
	updateEmployeeDetail := awscdklambdago.NewGoFunction(businessStack, jsii.String(ProjectPrefix+"UpdateEmployee"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "UpdateEmployee"),
		Description:  jsii.String("Employee can update their detail like email language name."),
		Entry:        jsii.String("../go-lambda-func/update-employee"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  lambdaEnvironmentVariable,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
	})

	userIdPath := userPath.AddResource(jsii.String("{userid}"), &awscdkrest.ResourceOptions{})

	// Handler for: PUT /user/{userid}
	userIdPath.AddMethod(jsii.String("PUT"), awscdkrest.NewLambdaIntegration(updateEmployeeDetail, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	// Get Business for the user assigned
	getBusinessAssignedToUser := awscdklambdago.NewGoFunction(businessStack, jsii.String(ProjectPrefix+"GetBusinessForUser"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String(ProjectPrefix + "GetBusinessForUser"),
		Description:  jsii.String("List all the business assigned to user."),
		Entry:        jsii.String("../go-lambda-func/get-business-for-user"),
		Role:         role,
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Environment:  lambdaEnvironmentVariable,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
	})

	userBusiness := userIdPath.AddResource(jsii.String("business"), &awscdkrest.ResourceOptions{})

	// Handler for: GET /user/{userid}/business
	userBusiness.AddMethod(jsii.String("GET"), awscdkrest.NewLambdaIntegration(getBusinessAssignedToUser, &awscdkrest.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}), &awscdkrest.MethodOptions{})

	return businessStack
}

type ImageStackProps struct {
	awscdk.StackProps
	RoleName *string
}

func NewImageStack(scope constructs.Construct, id string, props *ImageStackProps) awscdk.Stack {
	var stackProps awscdk.StackProps
	if props != nil {
		stackProps = props.StackProps
	}
	imageStack := awscdk.NewStack(scope, &id, &stackProps)

	//lambdaEnvironmentVariable := &map[string]*string{
	//	"DBTable":   jsii.String("XPOS_DEV"),
	//	"DataTable": jsii.String("XPOS_DATA"),
	//	//"CognitoUserPool": jsii.String("ap-south-1_gXgaeT7lu"),
	//	"CognitoUserPool": awscdk.Fn_ImportValue(jsii.String(ProjectPrefix + "UserPool")),
	//}

	//lambdaLayers := awslambda.LayerVersion_FromLayerVersionArn(imageStack, jsii.String(ProjectPrefix+"DependencyLayer"), jsii.String("arn:aws:lambda:ap-south-1:189468856814:layer:XPOSDependencyLayer:3"))

	// Role that needs to be attached to the lambda function
	//role := awsiam.Role_FromRoleName(imageStack, jsii.String("Role"), props.RoleName, &awsiam.FromRoleNameOptions{
	//	Mutable: jsii.Bool(false),
	//})

	//imageDlq := awscdksqs.NewQueue(imageStack, jsii.String(ProjectPrefix+"ZipImageDLQ"), &awscdksqs.QueueProps{
	//	QueueName: jsii.String(ProjectPrefix + "ZipImageDLQ"),
	//})
	//
	//imageSqs := awscdksqs.NewQueue(imageStack, jsii.String(ProjectPrefix+"ZipImageQ"), &awscdksqs.QueueProps{
	//	QueueName: jsii.String(ProjectPrefix + "ZipImageQ"),
	//	DeadLetterQueue: &awscdksqs.DeadLetterQueue{
	//		MaxReceiveCount: jsii.Number(2),
	//		Queue:           imageDlq,
	//	},
	//})

	//inputImageBucket := awss3.Bucket_FromBucketName(imageStack, jsii.String("xpos-image-stage"), jsii.String("xpos-image-stage"))

	//inputImageBucket.AddEventNotification(awss3.EventType_OBJECT_CREATED_PUT, awss3notifications.NewSqsDestination(imageSqs), &awss3.NotificationKeyFilter{
	//	Prefix: jsii.String("fileImport/"),
	//	Suffix: jsii.String(".zip"),
	//})

	// Add the lambda function to read the s3 file and extract it
	//extractS3ImageLambda := awscdklambdago.NewGoFunction(imageStack, jsii.String(ProjectPrefix+"ExtractZipFileLambda"), &awscdklambdago.GoFunctionProps{
	//	FunctionName: jsii.String(ProjectPrefix + "ExtractZipFileLambda"),
	//	Description:  jsii.String("Extract zip files from s3 bucket."),
	//	Entry:        jsii.String("../go-lambda-func/s3-extract-zip"),
	//	Role:         role,
	//	Runtime:      awslambda.Runtime_GO_1_X(),
	//	MemorySize:   jsii.Number(128),
	//	Timeout:      awscdk.Duration_Seconds(jsii.Number(300)),
	//	LogRetention: awslogs.RetentionDays_ONE_WEEK,
	//})

	//sqsEventSource := awslambdaeventsources.NewSqsEventSource(imageSqs, &awslambdaeventsources.SqsEventSourceProps{
	//	BatchSize: jsii.Number(1),
	//	Enabled:   jsii.Bool(true),
	//})
	//
	//extractS3ImageLambda.AddEventSource(sqsEventSource)

	//compressImageLambda := awscdklambdanodejs.NewNodejsFunction(imageStack, jsii.String(ProjectPrefix+"ImageCompressorLambda"), &awscdklambdanodejs.NodejsFunctionProps{
	//	FunctionName: jsii.String(ProjectPrefix + "ImageCompressorLambda"),
	//	Entry:        jsii.String("../nodejs-lambda-func/image-processing/compress-image/compress_images.ts"),
	//	Handler:      jsii.String("handler"),
	//	Runtime:      awslambda.Runtime_NODEJS_16_X(),
	//	MemorySize:   jsii.Number(256),
	//	Timeout:      awscdk.Duration_Seconds(jsii.Number(120)),
	//	Role:         role,
	//	Environment: &map[string]*string{
	//		"OUT_IMAGE_QUALITY": jsii.String("50"),
	//		"TARGET_BUCKET":     jsii.String("xpos-image"),
	//	},
	//	Bundling: &awscdklambdanodejs.BundlingOptions{
	//		ExternalModules: &[]*string{
	//			jsii.String("aws-sdk"),
	//			jsii.String("sharp"),
	//		},
	//	},
	//	Layers: &[]awslambda.ILayerVersion{
	//		lambdaLayers,
	//	},
	//	LogRetention: awslogs.RetentionDays_ONE_WEEK,
	//})
	//
	////// Create listener for the image upload notification
	//inputImageBucket.AddEventNotification(awss3.EventType_OBJECT_CREATED_PUT, awss3notifications.NewLambdaDestination(compressImageLambda), &awss3.NotificationKeyFilter{
	//	Prefix: jsii.String("output-zip/"),
	//})

	return imageStack
}
