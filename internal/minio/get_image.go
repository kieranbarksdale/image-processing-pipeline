package minio

import (
	"io"
	minioSDK "github.com/minio/minio-go/v7"
	"context"
)

func (minioClient *MinioClient) GetImage(ctx context.Context, bucketName string, objectName string) ([]byte, error) {
	obj, err := minioClient.Client.GetObject(ctx, bucketName, objectName, minioSDK.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()
	_, err = obj.Stat()
	if err != nil {
		return nil, err
	}
	// this puts the image in a byte slice 
	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}

	return data, nil
}