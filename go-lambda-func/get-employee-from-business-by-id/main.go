package main

import (
	"get-employee-from-business-by-id/src"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(src.GetEmployeeFromBusinessById)
}
