package filestorage

import (
	"2025_2_404/internal/service/storage/config"
	"2025_2_404/pkg/globalerrors"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		AppConfig: &config.AppConfig{
			ImgPath: "/tmp/test",
		},
	}

	useCase := New(cfg)

	if useCase == nil {
		t.Fatal("expected non-nil useCase")
	}

	if useCase.BaseDir() != "/tmp/test" {
		t.Errorf("expected baseDir '/tmp/test', got %s", useCase.BaseDir())
	}
}

func TestCreate_Success(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		AppConfig: &config.AppConfig{
			ImgPath: tempDir,
		},
	}

	useCase := New(cfg)
	imageData := []byte("test image data")
	imagePath := "test/image.png"

	err := useCase.Create(context.Background(), imageData, imagePath)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	fullPath := filepath.Join(tempDir, imagePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("expected file to exist")
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(data) != string(imageData) {
		t.Errorf("expected data %s, got %s", imageData, data)
	}
}

func TestCreate_EmptyPath(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		AppConfig: &config.AppConfig{
			ImgPath: tempDir,
		},
	}

	useCase := New(cfg)
	imageData := []byte("test image data")

	err := useCase.Create(context.Background(), imageData, "")
	if err != globalerrors.ErrInvalidPath {
		t.Errorf("expected ErrInvalidPath, got %v", err)
	}
}

func TestCreate_PathTraversal(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		AppConfig: &config.AppConfig{
			ImgPath: tempDir,
		},
	}

	useCase := New(cfg)
	imageData := []byte("test image data")

	testCases := []string{
		"../../../etc/passwd",
		"/etc/passwd",
		"test/../../../etc/passwd",
		"./test/../../../etc/passwd",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			err := useCase.Create(context.Background(), imageData, tc)
			if err != globalerrors.ErrInvalidPath {
				t.Errorf("expected ErrInvalidPath for path %s, got %v", tc, err)
			}
		})
	}
}

func TestGet_Success(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		AppConfig: &config.AppConfig{
			ImgPath: tempDir,
		},
	}

	useCase := New(cfg)
	imageData := []byte("test image data")
	imagePath := "test/image.png"

	err := useCase.Create(context.Background(), imageData, imagePath)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	data, contentType, err := useCase.Get(context.Background(), imagePath)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if string(data) != string(imageData) {
		t.Errorf("expected data %s, got %s", imageData, data)
	}

	if contentType != "image/png" {
		t.Errorf("expected content type 'image/png', got %s", contentType)
	}
}

func TestGet_FileNotFound(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		AppConfig: &config.AppConfig{
			ImgPath: tempDir,
		},
	}

	useCase := New(cfg)

	_, _, err := useCase.Get(context.Background(), "nonexistent/file.png")
	if err != globalerrors.ErrFileNotFound {
		t.Errorf("expected ErrFileNotFound, got %v", err)
	}
}

func TestGet_EmptyPath(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		AppConfig: &config.AppConfig{
			ImgPath: tempDir,
		},
	}

	useCase := New(cfg)

	_, _, err := useCase.Get(context.Background(), "")
	if err != globalerrors.ErrInvalidPath {
		t.Errorf("expected ErrInvalidPath, got %v", err)
	}
}

func TestGet_PathTraversal(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		AppConfig: &config.AppConfig{
			ImgPath: tempDir,
		},
	}

	useCase := New(cfg)

	testCases := []string{
		"../../../etc/passwd",
		"/etc/passwd",
		"test/../../../etc/passwd",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			_, _, err := useCase.Get(context.Background(), tc)
			if err != globalerrors.ErrInvalidPath {
				t.Errorf("expected ErrInvalidPath for path %s, got %v", tc, err)
			}
		})
	}
}

func TestDelete_Success(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		AppConfig: &config.AppConfig{
			ImgPath: tempDir,
		},
	}

	useCase := New(cfg)
	imageData := []byte("test image data")
	imagePath := "test/image.png"

	err := useCase.Create(context.Background(), imageData, imagePath)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	err = useCase.Delete(context.Background(), imagePath)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	fullPath := filepath.Join(tempDir, imagePath)
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		t.Error("expected file to be deleted")
	}
}

func TestDelete_FileNotExist(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		AppConfig: &config.AppConfig{
			ImgPath: tempDir,
		},
	}

	useCase := New(cfg)

	err := useCase.Delete(context.Background(), "nonexistent/file.png")
	if err != nil {
		t.Errorf("expected no error for deleting non-existent file, got %v", err)
	}
}

func TestDelete_EmptyPath(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		AppConfig: &config.AppConfig{
			ImgPath: tempDir,
		},
	}

	useCase := New(cfg)

	err := useCase.Delete(context.Background(), "")
	if err != globalerrors.ErrInvalidPath {
		t.Errorf("expected ErrInvalidPath, got %v", err)
	}
}

func TestDelete_PathTraversal(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		AppConfig: &config.AppConfig{
			ImgPath: tempDir,
		},
	}

	useCase := New(cfg)

	testCases := []string{
		"../../../etc/passwd",
		"/etc/passwd",
		"test/../../../etc/passwd",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			err := useCase.Delete(context.Background(), tc)
			if err != globalerrors.ErrInvalidPath {
				t.Errorf("expected ErrInvalidPath for path %s, got %v", tc, err)
			}
		})
	}
}

func TestGet_ContentTypeNoExtension(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		AppConfig: &config.AppConfig{
			ImgPath: tempDir,
		},
	}

	useCase := New(cfg)
	imageData := []byte("test data")
	imagePath := "test/file"

	err := useCase.Create(context.Background(), imageData, imagePath)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	_, contentType, err := useCase.Get(context.Background(), imagePath)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if contentType != "application/octet-stream" {
		t.Errorf("expected content type 'application/octet-stream', got %s", contentType)
	}
}
