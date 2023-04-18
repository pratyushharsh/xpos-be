package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/imagekit-developer/imagekit-go"
)

type ImageKitSignatureResponse struct {
	Token     string `json:"token"`
	Expires   int64  `json:"expires"`
	Signature string `json:"signature"`
}

func GenerateImageKitToken(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	privateKey := "private_ZzWCv71lj1LeNyRCSbSF+7anzAE="
	publicKey := "xxxxxxxx"
	urlEndpoint := "https://ik.imagekit.io/6z4jinkib"
	ik := imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		UrlEndpoint: urlEndpoint,
	})

	// Using auto-generated token and expiration
	resp := ik.SignToken(imagekit.SignTokenParam{})
	res, _ := json.Marshal(ImageKitSignatureResponse{
		Token:     resp.Token,
		Expires:   resp.Expires,
		Signature: resp.Signature,
	})

	return events.APIGatewayProxyResponse{Body: string(res), StatusCode: 200}, nil
}

func main() {
	lambda.Start(GenerateImageKitToken)

}
