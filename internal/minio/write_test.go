package minio

import (
	"strings"
	"testing"
)

func extractKeyFromURL(fullURL string) string {
	parts := strings.Split(fullURL, "/")
	return parts[len(parts)-1]
}

func TestWriteImage(t *testing.T) {

	data := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0, 0, 0} 
	reader := strings.NewReader(string(data))

	_, err := WriteImage(testClient, reader, int64(len(data)))
	if err != nil {
		t.Fatalf("failed to write image: %v", err)
	} 

	t.Log("Image written successfully and running through the main_test and booted from the testutil")
}

func TestWriteImage_InvalidContentType(t *testing.T) {

	invalidData := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05} 
	reader := strings.NewReader(string(invalidData))

	_, err := WriteImage(testClient, reader, int64(len(invalidData)))
	if err == nil {
		t.Fatalf("expected error for invalid content type, got nil")
	}

	t.Log("Invalid content type handled correctly")
}

func TestWriteImage_EmptyData(t *testing.T) {

	emptyData := []byte{}
	reader := strings.NewReader(string(emptyData))

	_, err := WriteImage(testClient, reader, int64(len(emptyData)))
	if err == nil {
		t.Fatalf("expected error for empty data, got nil")
	}

	t.Log("Empty data handled correctly")
}



