package config

import (
	"os"
	"testing"
)

func TestGetPostgresConfig(t *testing.T) {
	os.Setenv("POSTGRES_USER", "storageuser")
	os.Setenv("POSTGRES_PASSWORD", "storagepass")
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_DB", "storagedb")
	defer func() {
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("POSTGRES_DB")
	}()

	cfg := GetPostgresConfig()

	if cfg.User != "storageuser" {
		t.Errorf("expected User 'storageuser', got %s", cfg.User)
	}
	if cfg.Password != "storagepass" {
		t.Errorf("expected Password 'storagepass', got %s", cfg.Password)
	}
	if cfg.Host != "localhost" {
		t.Errorf("expected Host 'localhost', got %s", cfg.Host)
	}
	if cfg.Port != "5432" {
		t.Errorf("expected Port '5432', got %s", cfg.Port)
	}
	if cfg.DB != "storagedb" {
		t.Errorf("expected DB 'storagedb', got %s", cfg.DB)
	}
}

func TestGetAppConfig(t *testing.T) {
	os.Setenv("APP_HOST", "0.0.0.0")
	os.Setenv("APP_PORT", "9000")
	os.Setenv("GRPC_AD_PORT", "50051")
	os.Setenv("GRPC_STORAGE_PORT", "50052")
	os.Setenv("IMG_PATH", "/storage/images")
	defer func() {
		os.Unsetenv("APP_HOST")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("GRPC_AD_PORT")
		os.Unsetenv("GRPC_STORAGE_PORT")
		os.Unsetenv("IMG_PATH")
	}()

	cfg := GetAppConfig()

	if cfg.Host != "0.0.0.0" {
		t.Errorf("expected Host '0.0.0.0', got %s", cfg.Host)
	}
	if cfg.Port != "9000" {
		t.Errorf("expected Port '9000', got %s", cfg.Port)
	}
	if cfg.PortAD != "50051" {
		t.Errorf("expected PortAD '50051', got %s", cfg.PortAD)
	}
	if cfg.PortStorage != "50052" {
		t.Errorf("expected PortStorage '50052', got %s", cfg.PortStorage)
	}
	if cfg.ImgPath != "/storage/images" {
		t.Errorf("expected ImgPath '/storage/images', got %s", cfg.ImgPath)
	}
}

func TestGetPostgresConfig_EmptyEnv(t *testing.T) {
	os.Clearenv()

	cfg := GetPostgresConfig()

	if cfg.User != "" {
		t.Errorf("expected empty User, got %s", cfg.User)
	}
}
