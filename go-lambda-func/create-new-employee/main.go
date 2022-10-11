package main

import (
	"create-new-employee/src"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(src.CreateNewBusinessEmployee)
}
