package main

import (
	"context"
	f1 "create-new-business/src"
	f2 "get-business-by-id/src"
	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	f3 "update-business/src"
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

		body, err := ioutil.ReadAll(c.Request.Body)

		req := events.APIGatewayProxyRequest{
			Resource:                        "",
			Path:                            c.FullPath(),
			HTTPMethod:                      c.Request.Method,
			Headers:                         headers,
			MultiValueHeaders:               nil,
			QueryStringParameters:           nil,
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
	//router.GET("/business/:businessId/user", WrapLambdaFunction(user.GetUserForStoreHandler))
	//router.POST("/business/:businessId/user", WrapLambdaFunction(user.CreateNewStoreEmployee))
	//router.GET("/business/:businessId/user/:userid", WrapLambdaFunction(user.GetEmployeeFromStoreById))
	//router.PUT("/user/:userid", WrapLambdaFunction(user.UpdateEmployeeHandler))
	//router.GET("/user/:userid/business", WrapLambdaFunction(user.GetBusinessForUserHandler))
}

func StartApplication() {
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
