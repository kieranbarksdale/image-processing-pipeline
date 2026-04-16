package minio

import (
	"io"
	minioSDK "github.com/minio/minio-go/v7"
	"context"
	"strings"
	"log"
)

func (minioClient *MinioClient) GetImageStream(ctx context.Context, bucketName string, objectName string) (io.ReadCloser, error) {
	objectName = strings.TrimSpace(objectName)
	objectName = strings.TrimPrefix(objectName, "/")

	// This doesn't download the object yet 
	obj, err := minioClient.Client.GetObject(ctx, bucketName, objectName, minioSDK.GetObjectOptions{})
	if err != nil {
		log.Printf("Error getting object: %v", err)
		return nil, err
	}

	// Get the object stats to make sure the entire thing is there
	_, err = obj.Stat()
	if err != nil {
		obj.Close()
		log.Printf("Minio stat error for [%s]: %s", objectName, err.Error())
		return nil, err
	}

	return obj, nil
}