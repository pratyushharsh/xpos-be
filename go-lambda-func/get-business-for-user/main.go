package main

import (
	"get-business-for-user/src"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(src.GetBusinessForUserHandler)
}
