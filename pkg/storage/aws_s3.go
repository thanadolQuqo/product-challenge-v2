package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	product_config "product-challenge/pkg/config"
)

func NewS3Client(ctx context.Context, productCfg *product_config.Config) (*s3.Client, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(productCfg.Aws.Region))
	if err != nil {
		fmt.Println("Unable to load AWS config: ", err)
		return nil, err
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)
	return client, nil
}

// may not need these functions

// ListBucket is for checking if correct or not
func ListBucket(ctx context.Context, client *s3.Client) []string {
	// Call ListBuckets API
	resp, err := client.ListBuckets(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list S3 buckets: %v", err)
	}
	var bucketList []string

	for _, bucket := range resp.Buckets {
		bucketList = append(bucketList, *bucket.Name)
	}
	return bucketList
}

func UploadFileToS3(client *s3.Client, bucketName, fileName string, file multipart.File) (string, error) {
	region := os.Getenv("AWS_REGION")

	// Determine the content type
	contentType := http.DetectContentType(make([]byte, 512))
	_, err := file.Read(make([]byte, 512)) // Pre-read a portion of the file
	if err != nil {
		contentType = "application/octet-stream"
	}

	// Reset the file pointer to the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("failed to reset file pointer: %v", err)
	}

	// Upload the file
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileName),
		Body:        file, // Convert to a reader
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// Generate the S3 public URL
	s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, fileName)

	return s3URL, nil
}
