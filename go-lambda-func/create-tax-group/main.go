package main

import (
	"create-tax-group/src"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(src.CreateNewTaxGroupHandler)
}
