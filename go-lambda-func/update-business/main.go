package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"update-business/src"
)

func main() {
	lambda.Start(src.UpdateBusinessHandler)
}
