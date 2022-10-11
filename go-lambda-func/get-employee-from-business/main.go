package main

import (
	"get-employee-from-business/src"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(src.GetEmployeeFromBusinessHandler)
}
