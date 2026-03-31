package minio

// import (
// 	"bytes"
// 	"context"
// 	"github.com/minio/minio-go/v7"
// )

// func WriteImage(client *MinioClient, bucketName, fileName string, data []byte) (string, error) {
// 	filePath := bucketName + "/" + fileName
// 	_, err := client.Client.PutObject(context.Background(), bucketName, fileName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
// 	if err != nil {
// 		return "", err
// 	}

// 	return filePath, nil
// }