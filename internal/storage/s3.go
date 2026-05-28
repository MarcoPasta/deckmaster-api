package storage

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

type Presigner struct {
	Client        *s3.Client
	PresignClient *s3.PresignClient
}

func CreateS3Client() Presigner {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	S3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = &endpoint
			o.UsePathStyle = true
		}
	})
	presignClient := s3.NewPresignClient(S3Client)
	presigner := Presigner{Client: S3Client, PresignClient: presignClient}

	return presigner
}

func (presigner Presigner) ListObjects(ctx *gin.Context, prefix string) ([]string, error) {
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	log.Println("AWS_S3_BUCKET_NAME: ", bucketName)
	var keys []string
	paginator := s3.NewListObjectsV2Paginator(presigner.Client, &s3.ListObjectsV2Input{
		Bucket: &bucketName,
		Prefix: aws.String(prefix),
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, obj := range page.Contents {
			keys = append(keys, *obj.Key)
		}
	}
	return keys, nil
}

func (presigner Presigner) GetObject(ctx *gin.Context, objectKey string) (*v4.PresignedHTTPRequest, error) {
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	request, err := presigner.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(3600 * int64(time.Second))
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
