package convertimage

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestConvertImageToBase64_JPG(t *testing.T) {
	imageData := []byte("test image data")
	result := ConvertImageToBase64(imageData, "jpg")

	expectedEncoding := base64.StdEncoding.EncodeToString(imageData)
	expectedResult := "data:image/jpeg;base64," + expectedEncoding

	if result != expectedResult {
		t.Errorf("expected %s, got %s", expectedResult, result)
	}
}

func TestConvertImageToBase64_PNG(t *testing.T) {
	imageData := []byte("test png data")
	result := ConvertImageToBase64(imageData, "image/png")

	expectedEncoding := base64.StdEncoding.EncodeToString(imageData)
	expectedResult := "data:image/png;base64," + expectedEncoding

	if result != expectedResult {
		t.Errorf("expected %s, got %s", expectedResult, result)
	}
}

func TestConvertImageToBase64_WebP(t *testing.T) {
	imageData := []byte("test webp data")
	result := ConvertImageToBase64(imageData, "image/webp")

	expectedEncoding := base64.StdEncoding.EncodeToString(imageData)
	expectedResult := "data:image/webp;base64," + expectedEncoding

	if result != expectedResult {
		t.Errorf("expected %s, got %s", expectedResult, result)
	}
}

func TestConvertImageToBase64_EmptyData(t *testing.T) {
	imageData := []byte{}
	result := ConvertImageToBase64(imageData, "jpg")

	if !strings.HasPrefix(result, "data:image/jpeg;base64,") {
		t.Errorf("expected result to start with 'data:image/jpeg;base64,', got %s", result)
	}
}

func TestConvertImageToBase64_LargeData(t *testing.T) {
	imageData := make([]byte, 10000)
	for i := range imageData {
		imageData[i] = byte(i % 256)
	}

	result := ConvertImageToBase64(imageData, "image/png")

	expectedEncoding := base64.StdEncoding.EncodeToString(imageData)
	expectedResult := "data:image/png;base64," + expectedEncoding

	if result != expectedResult {
		t.Errorf("base64 encoding mismatch for large data")
	}
}

func TestConvertImageToBase64_ValidBase64Output(t *testing.T) {
	imageData := []byte("test data for validation")
	result := ConvertImageToBase64(imageData, "jpg")

	parts := strings.Split(result, ",")
	if len(parts) != 2 {
		t.Fatalf("expected result to have 2 parts separated by comma, got %d", len(parts))
	}

	encodedPart := parts[1]
	_, err := base64.StdEncoding.DecodeString(encodedPart)
	if err != nil {
		t.Errorf("result contains invalid base64: %v", err)
	}
}
