package minio

import (
	"fmt"
	"io"
	minioSDK "github.com/minio/minio-go/v7"
	"context"
	"strings"
)

func (minioClient *MinioClient) GetImage(ctx context.Context, bucketName string, objectName string) ([]byte, error) {
	objectName = strings.TrimSpace(objectName)
	objectName = strings.TrimPrefix(objectName, "/")
	obj, err := minioClient.Client.GetObject(ctx, bucketName, objectName, minioSDK.GetObjectOptions{})
	if err != nil {
		fmt.Println("ERRRRRRRRRR:", err)
		return nil, err
	}
	fmt.Println("Object retrieved successfully")
	defer obj.Close()
	_, err = obj.Stat()
	if err != nil {
		fmt.Printf("minio stat error for [%s]: %s\n", objectName, err.Error())
		return nil, err
	}
	fmt.Println("Object stats retrieved successfully")
	// this puts the image in a byte slice 
	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}
	fmt.Println("Image data read successfully")

	return data, nil
}