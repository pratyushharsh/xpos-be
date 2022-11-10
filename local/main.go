package main

import (
	"context"
	f1 "create-new-business/src"
	f5 "create-new-employee/src"
	f9 "create-tax-group/src"
	"fmt"
	f2 "get-business-by-id/src"
	f7 "get-business-for-user/src"
	f6 "get-employee-from-business-by-id/src"
	f4 "get-employee-from-business/src"
	f12 "get-sync-data/src"
	f10 "get-tax-group-for-store/src"
	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	f3 "update-business/src"
	f8 "update-employee/src"
	f11 "update-sync-data/src"
)

type LambdaApiRequest func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

var (
	router *gin.Engine
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "Success")
}

func WrapLambdaFunction(inp LambdaApiRequest) gin.HandlerFunc {
	return func(c *gin.Context) {

		var pathParams = make(map[string]string)

		for _, param := range c.Params {
			pathParams[param.Key] = param.Value
		}

		var headers = make(map[string]string)

		for key, value := range c.Request.Header {
			headers[key] = value[0]
		}

		var queryStrings = make(map[string]string)
		for key, value := range c.Request.URL.Query() {
			queryStrings[key] = value[0]
		}

		body, err := io.ReadAll(c.Request.Body)

		req := events.APIGatewayProxyRequest{
			Resource:                        "",
			Path:                            c.FullPath(),
			HTTPMethod:                      c.Request.Method,
			Headers:                         headers,
			MultiValueHeaders:               nil,
			QueryStringParameters:           queryStrings,
			MultiValueQueryStringParameters: nil,
			PathParameters:                  pathParams,
			StageVariables:                  nil,
			Body:                            string(body),
			IsBase64Encoded:                 false,
		}
		response, err := inp(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.Header("Content-Type", "application/json")
		c.String(response.StatusCode, response.Body)
	}
}

func urlMappings() {
	router.GET("/ping", Ping)
	router.POST("/business", WrapLambdaFunction(f1.CreateBusinessHandler))
	router.GET("/business/:businessId", WrapLambdaFunction(f2.GetBusinessByIdHandler))
	router.PUT("/business/:businessId", WrapLambdaFunction(f3.UpdateBusinessHandler))

	// Configure Sync Service
	router.GET("/business/:businessId/sync", WrapLambdaFunction(f12.GetSyncDataHandler))
	router.POST("/business/:businessId/sync", WrapLambdaFunction(f11.UpdateSyncHandler))

	// Configuration for settings
	router.POST("/business/:businessId/settings/tax", WrapLambdaFunction(f9.CreateNewTaxGroupHandler))
	router.GET("/business/:businessId/settings/tax", WrapLambdaFunction(f10.GetTaxGroupForStore))

	router.GET("/business/:businessId/employee", WrapLambdaFunction(f4.GetEmployeeFromBusinessHandler))
	router.POST("/business/:businessId/employee", WrapLambdaFunction(f5.CreateNewBusinessEmployee))
	router.GET("/business/:businessId/employee/:userid", WrapLambdaFunction(f6.GetEmployeeFromBusinessById))
	router.GET("/user/:userid/business", WrapLambdaFunction(f7.GetBusinessForUserHandler))
	router.PUT("/user/:userid", WrapLambdaFunction(f8.UpdateEmployeeHandler))
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		log.Println("Loaded .env file")
	}

}

func StartApplication() {
	fmt.Println(os.Getenv("GIN_MODE"))
	fmt.Println(os.Getenv("PROJECT_NAME"))
	//gin.SetMode(os.Getenv("GIN_MODE"))
	router = gin.Default()
	urlMappings()

	err := router.Run(":9090")
	if err != nil {
		panic(err)
	}
}

func main() {
	StartApplication()
}
