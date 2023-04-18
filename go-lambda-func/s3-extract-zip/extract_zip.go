package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"image"
	"image/jpeg"
	"io"
	"sync"
)

type S3Download struct {
	BucketName   string
	ObjectKey    string
	ObjectLength int64
}

func (sd *S3Download) ReadAt(p []byte, offset int64) (int, error) {
	// #1
	if offset < 0 || offset >= sd.ObjectLength {
		return 0, fmt.Errorf("invalid offset")
	}
	svc := GetS3Client()
	input := &s3.GetObjectInput{
		Bucket: aws.String(sd.BucketName),
		Key:    aws.String(sd.ObjectKey),
		// #2
		Range: aws.String(fmt.Sprintf("bytes=%d-%d", offset, offset+int64(len(p)))),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return 0, fmt.Errorf(s3.ErrCodeNoSuchKey, aerr.Error())
			case s3.ErrCodeInvalidObjectState:
				return 0, fmt.Errorf(s3.ErrCodeInvalidObjectState, aerr.Error())
			default:
				return 0, fmt.Errorf(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return 0, fmt.Errorf(err.Error())
		}
	}
	// #3
	n, err := result.Body.Read(p)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return n, nil
		}
		return 0, err
	}
	return n, nil
}

//func (sd *S3Download) S3ObjectLength() (int64, error) {
//	svc := GetS3Client()
//	input := &s3.GetObjectInput{
//		Bucket: aws.String(sd.BucketName),
//		Key:    aws.String(sd.ObjectKey),
//	}
//
//	result, err := svc.GetObject(input)
//	if err != nil {
//		if aerr, ok := err.(awserr.Error); ok {
//			switch aerr.Code() {
//			case s3.ErrCodeNoSuchKey:
//				return 0, fmt.Errorf(s3.ErrCodeNoSuchKey, aerr.Error())
//			case s3.ErrCodeInvalidObjectState:
//				return 0, fmt.Errorf(s3.ErrCodeInvalidObjectState, aerr.Error())
//			default:
//				return 0, aerr.OrigErr()
//			}
//		} else {
//			// Print the error, cast err to awserr.Error to get the Code and
//			// Message from an error.
//			fmt.Errorf(err.Error())
//		}
//	}
//	sd.ObjectLength = *result.ContentLength
//	return *result.ContentLength, nil
//}

func compressImageResource(data []byte) []byte {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return data
	}
	//newSrc := resize.Resize(1000, 0, img, resize.Lanczos3)
	buf := bytes.Buffer{}
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 20})
	if err != nil {
		return data
	}
	if buf.Len() > len(data) {
		return data
	}
	return buf.Bytes()
}

func uploadToS3(WaitGroup *sync.WaitGroup, bucketName string, key string, body io.Reader) {
	defer WaitGroup.Done()

	// Get Bytes from reader
	data, err := io.ReadAll(body)
	if err != nil {
		fmt.Println(err)
		return
	}

	compressedImage := compressImageResource(data)

	// Convert Byte[] to io.Reader
	s3Body := bytes.NewReader(compressedImage)

	_, err = GetS3Uploader().Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   s3Body,
	})
	if err != nil {
		fmt.Println(err)
	}
}
