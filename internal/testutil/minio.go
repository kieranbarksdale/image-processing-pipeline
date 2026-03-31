package testutil

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	minioContainer "github.com/testcontainers/testcontainers-go/modules/minio"
)

func SetupMinio(ctx context.Context) (*minio.Client, func(), error) {

	container, err := minioContainer.Run(ctx, "minio/minio:latest")
	if err != nil {
		return nil, nil, err
	}

	accessKey := "minioadmin"
	secretKey := "minioadmin"
	
	endpoint, _ := container.ConnectionString(ctx)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, nil, err
	}
	
	cleanup := func() {
		container.Terminate(ctx)
	}
	
	return client, cleanup, nil
}


