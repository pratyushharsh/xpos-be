package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"update-sync-data/src"
)

func main() {
	lambda.Start(src.UpdateSyncHandler)
}
