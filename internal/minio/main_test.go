package minio

import (
	"context"
	"os"
	"testing"
	"image-processing-pipeline/internal/testutil"
	"github.com/minio/minio-go/v7"
)

var testClient *MinioClient

func TestMain(m *testing.M) {
	ctx := context.Background()
	client, cleanup, err := testutil.SetupMinio(ctx)
	if err != nil {
		os.Exit(1)
	}
	
	testClient = &MinioClient{
		Client: client,
		Bucket: "integration-test-bucket",
	} 

	_ = client.MakeBucket(ctx, testClient.Bucket, minio.MakeBucketOptions{})
	code := m.Run()
	cleanup()
	os.Exit(code)
}
