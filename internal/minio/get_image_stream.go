package minio

import (
	"fmt"
	"io"
	minioSDK "github.com/minio/minio-go/v7"
	"context"
	"strings"
)

func (minioClient *MinioClient) GetImageStream(ctx context.Context, bucketName string, objectName string) (io.ReadCloser, error) {
	objectName = strings.TrimSpace(objectName)
	objectName = strings.TrimPrefix(objectName, "/")

	// This doesn't download the object yet 
	obj, err := minioClient.Client.GetObject(ctx, bucketName, objectName, minioSDK.GetObjectOptions{})
	if err != nil {
		fmt.Println("ERRRRRRRRRR:", err)
		return nil, err
	}
	fmt.Println("Object retrieved successfully")

	// Get the object stats to make sure the entire thing is there
	_, err = obj.Stat()
	if err != nil {
		obj.Close()
		fmt.Printf("minio stat error for [%s]: %s\n", objectName, err.Error())
		return nil, err
	}
	fmt.Println("Object stats retrieved successfully")

	return obj, nil
}