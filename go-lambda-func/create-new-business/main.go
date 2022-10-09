package main

import (
	"create-new-business/src"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(src.CreateBusinessHandler)
}
