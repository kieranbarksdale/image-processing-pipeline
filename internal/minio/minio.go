package minio

import (
	"context"
	"fmt"
	"os"
	"time"
	minioSDK "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)


type MinioClient struct {
	Client *minioSDK.Client
	Bucket string
}

func waitUntilMinioIsReady(client *minioSDK.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i := 0; i < 10; i++ {
		_, err := client.ListBuckets(ctx)
		if err == nil {
			fmt.Println("Minio is ready")
			return
		}
		fmt.Printf("Minio is not ready yet, retrying... (%d/10) error: %v\n", i+1, err)
		time.Sleep(2 * time.Second)
	}
	panic("MinIO never became ready")
}


func NewClient(endpoint, accessKey, secretKey string, ssl bool) (*MinioClient, error) {
	client, err := minioSDK.New(endpoint, &minioSDK.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: ssl,
	})
	if err != nil {
		return nil, err
	}

	waitUntilMinioIsReady(client) 

	ctx := context.Background()
	bucketName := os.Getenv("MINIO_BUCKET")
	err = client.MakeBucket(ctx, bucketName, minioSDK.MakeBucketOptions{})
	if err != nil {
		exists, err := client.BucketExists(ctx, bucketName)
		if err == nil && exists {
			fmt.Println("Bucket '" + bucketName + "' already exists")
		} else {
			return nil, err
		}
	}
	
	return &MinioClient{Client: client, Bucket: bucketName}, nil
}