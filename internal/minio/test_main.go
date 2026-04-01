package minio

import (
	"context"
	"fmt"
	"os"
	"testing"
	"image-processing-pipeline/internal/testutil"
)

var testClient *MinioClient

func TestMain(m *testing.M) {
	ctx := context.Background()
	
	client, _, err := testutil.SetupMinio(ctx)
	if err != nil {
		fmt.Printf("Failed to setup minio: %v\n", err)
		os.Exit(1)
	}
	
	testClient = &MinioClient{
		Client: client,
		Bucket: "test-bucket",
	}
	
	// Run the tests
	os.Exit(m.Run())
}
