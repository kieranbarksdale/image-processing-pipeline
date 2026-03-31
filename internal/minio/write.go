package minio

import (
	"context"
	"time"
	"fmt"
	"net/http"
	"io"
	"strings"
	"github.com/minio/minio-go/v7"
	"github.com/google/uuid"
)


func WriteImage(client *MinioClient, file io.ReadSeeker, size int64) (string, error) {
	if client == nil {
		return "", fmt.Errorf("minio client is nil")
	}
	if size == 0 {
		return "", fmt.Errorf("data is empty")
	}

	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	contentType := http.DetectContentType(buffer) 

	if contentType != "image/jpeg" && 
		contentType != "image/png" && 
		contentType != "image/gif" && 
		contentType != "image/webp" {
		return "", fmt.Errorf("invalid content type")
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	uuid := uuid.New().String()
	objectKey := uuid + "_" + time.Now().Format("2006-01-02_15-04-05.000") + "." + strings.Split(contentType, "/")[1]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uploadInfo, err := client.Client.PutObject(
		ctx, 
		client.Bucket, 
		objectKey, 
		file, 
		size, 
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", err
	}
	fmt.Println("Uploaded", objectKey, "of size: ", uploadInfo.Size, "to bucket", client.Bucket + " key:" + uploadInfo.Key)

	fullURL := client.Client.EndpointURL().String() + "/" + client.Bucket + "/" + objectKey
	return fullURL, nil
}