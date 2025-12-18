package config

import (
	"os"
	"testing"
)

func TestGetPostgresConfig(t *testing.T) {
	// Set test environment variables
	os.Setenv("POSTGRES_USER", "testuser")
	os.Setenv("POSTGRES_PASSWORD", "testpass")
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_DB", "testdb")
	defer func() {
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("POSTGRES_DB")
	}()

	cfg := GetPostgresConfig()

	if cfg.User != "testuser" {
		t.Errorf("expected User 'testuser', got %s", cfg.User)
	}
	if cfg.Password != "testpass" {
		t.Errorf("expected Password 'testpass', got %s", cfg.Password)
	}
	if cfg.Host != "localhost" {
		t.Errorf("expected Host 'localhost', got %s", cfg.Host)
	}
	if cfg.Port != "5432" {
		t.Errorf("expected Port '5432', got %s", cfg.Port)
	}
	if cfg.DB != "testdb" {
		t.Errorf("expected DB 'testdb', got %s", cfg.DB)
	}
}

func TestGetAppConfig(t *testing.T) {
	os.Setenv("APP_HOST", "0.0.0.0")
	os.Setenv("APP_PORT", "8080")
	os.Setenv("GRPC_AD_PORT", "50051")
	os.Setenv("GRPC_STORAGE_PORT", "50052")
	os.Setenv("IMG_PATH", "/tmp/images")
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
	if cfg.Port != "8080" {
		t.Errorf("expected Port '8080', got %s", cfg.Port)
	}
	if cfg.PortAD != "50051" {
		t.Errorf("expected PortAD '50051', got %s", cfg.PortAD)
	}
	if cfg.PortStorage != "50052" {
		t.Errorf("expected PortStorage '50052', got %s", cfg.PortStorage)
	}
	if cfg.ImgPath != "/tmp/images" {
		t.Errorf("expected ImgPath '/tmp/images', got %s", cfg.ImgPath)
	}
}

func TestGetPostgresConfig_EmptyEnv(t *testing.T) {
	// Clear all environment variables
	os.Clearenv()

	cfg := GetPostgresConfig()

	if cfg.User != "" {
		t.Errorf("expected empty User, got %s", cfg.User)
	}
	if cfg.Password != "" {
		t.Errorf("expected empty Password, got %s", cfg.Password)
	}
}
