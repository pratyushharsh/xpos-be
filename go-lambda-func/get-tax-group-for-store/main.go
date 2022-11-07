package main

import (
	"get-tax-group-for-store/src"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(src.GetTaxGroupForStore)
}
