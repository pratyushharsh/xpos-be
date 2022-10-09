package main

import (
	"get-business-by-id/src"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(src.GetBusinessByIdHandler)
}
