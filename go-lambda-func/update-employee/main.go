package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"update-employee/src"
)

func main() {
	lambda.Start(src.UpdateEmployeeHandler)
}
