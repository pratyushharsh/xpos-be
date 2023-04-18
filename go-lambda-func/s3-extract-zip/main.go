package main

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"strings"
	"sync"
)

var (
	sess           *session.Session
	s3Client       *s3.S3
	s3UploadClient *s3manager.Uploader
)

func GetSession() *session.Session {
	if sess == nil {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}
	return sess
}

func GetS3Client() *s3.S3 {
	if s3Client == nil {
		s3Client = s3.New(GetSession())
	}
	return s3Client
}

func GetS3Uploader() *s3manager.Uploader {
	if s3UploadClient == nil {
		s3UploadClient = s3manager.NewUploader(GetSession())
	}
	return s3UploadClient
}

func ExtractZipToS3(ctx context.Context, event events.SQSEvent) (events.SQSEventResponse, error) {
	eventJson, _ := json.Marshal(event)
	log.Printf("EVENT: %s", eventJson)

	// Read the SQS message
	for _, record := range event.Records {
		// Parse the record as an S3 event
		var s3Event events.S3Event
		err := json.Unmarshal([]byte(record.Body), &s3Event)
		if err != nil {
			log.Printf("ERROR: %s", err)
			continue
		}
		// Process Each Record From S3 To Extract Zip
		for _, s3Record := range s3Event.Records {
			// Get the bucket name and key
			bucket := s3Record.S3.Bucket.Name
			key := s3Record.S3.Object.Key

			obj := S3Download{
				BucketName:   bucket,
				ObjectKey:    key,
				ObjectLength: s3Record.S3.Object.Size,
			}

			if err != nil {
				fmt.Println(err)
				return events.SQSEventResponse{}, err
			}
			// Create a new zip reader
			zipReader, err := zip.NewReader(&obj, s3Record.S3.Object.Size)
			if err != nil {
				fmt.Println(err)
				return events.SQSEventResponse{}, err
			}
			// Create a new wait group
			wg := sync.WaitGroup{}
			// Loop through the files in the zip
			for _, file := range zipReader.File {

				if !file.FileInfo().IsDir() && strings.HasSuffix(file.Name, ".jpg") {
					bytes, err := file.Open()
					if err != nil {
						fmt.Errorf("failed to upload file, %v", err)
					} else {
						// Increment the wait group
						wg.Add(1)
						go uploadToS3(&wg, "xpos-image-stage", "output-zip/"+file.Name, bytes)
					}
				}

			}
			// Wait for all the files to be processed
			wg.Wait()
		}
	}

	return events.SQSEventResponse{BatchItemFailures: nil}, nil
}

func main() {
	//lambda.Start(ExtractZipToS3)
	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			events.SQSMessage{Body: "{\"Records\":[{\"eventVersion\":\"2.1\",\"eventSource\":\"aws:s3\",\"awsRegion\":\"ap-south-1\",\"eventTime\":\"2022-11-20T18:35:19.255Z\",\"eventName\":\"ObjectCreated:Put\",\"userIdentity\":{\"principalId\":\"A1ROYY0TIYNPCC\"},\"requestParameters\":{\"sourceIPAddress\":\"182.69.243.191\"},\"responseElements\":{\"x-amz-request-id\":\"SG2WDSVHRBNFFPYH\",\"x-amz-id-2\":\"Qq4XovmIytk4/0kF09JJGZ04jlWKshOS4KxArGooMYra6h2i/rartYxmoK/DXkRIj1ivVXReKJUCMVCQttCUyezFE90OklJfO0Kc9IiYN5Y=\"},\"s3\":{\"s3SchemaVersion\":\"1.0\",\"configurationId\":\"arn:aws:cloudformation:ap-south-1:189468856814:stack/XPOS-ImageStack/7c863130-68b2-11ed-a07f-025a43c0aa8e--1606870337993232671\",\"bucket\":{\"name\":\"xpos-image-stage\",\"ownerIdentity\":{\"principalId\":\"A1ROYY0TIYNPCC\"},\"arn\":\"arn:aws:s3:::xpos-image-stage\"},\"object\":{\"key\":\"fileImport/image_compress_test.zip\",\"size\":16662420,\"eTag\":\"4a8c9d583a6653307fe857b89056c288\",\"sequencer\":\"00637A7364EA0034D1\"}}}]}"},
		},
	}
	ExtractZipToS3(context.Background(), sqsEvent)
}
