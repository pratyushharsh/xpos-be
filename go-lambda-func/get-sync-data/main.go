package main

import (
	"get-sync-data/src"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(src.GetSyncDataHandler)
}
