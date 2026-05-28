package controller

import (
	"log"
	"os"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

type Presigner struct {
	PresignClient *s3.PresignClient
}

func (presigner Presigner) GetObject(ctx *gin.Context, objectKey string) (*v4.PresignedHTTPRequest, error) {
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	request, err := presigner.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(120 * int64(time.Second)) // keep url alive for 2 minutes
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

func (presigner Presigner) PutObject(ctx *gin.Context, objectKey string) (*v4.PresignedHTTPRequest, error) {
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	request, err := presigner.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(120 * int64(time.Second)) // keep url alive for 2 minutes
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}
