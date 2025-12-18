package config

import (
	"os"
	"testing"
)

func TestGetPostgresConfig(t *testing.T) {
	os.Setenv("POSTGRES_USER", "slotuser")
	os.Setenv("POSTGRES_PASSWORD", "slotpass")
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_DB", "slotdb")
	defer func() {
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("POSTGRES_DB")
	}()

	cfg := GetPostgresConfig()

	if cfg.User != "slotuser" {
		t.Errorf("expected User 'slotuser', got %s", cfg.User)
	}
	if cfg.Password != "slotpass" {
		t.Errorf("expected Password 'slotpass', got %s", cfg.Password)
	}
	if cfg.Host != "localhost" {
		t.Errorf("expected Host 'localhost', got %s", cfg.Host)
	}
	if cfg.Port != "5432" {
		t.Errorf("expected Port '5432', got %s", cfg.Port)
	}
	if cfg.DB != "slotdb" {
		t.Errorf("expected DB 'slotdb', got %s", cfg.DB)
	}
}

func TestGetAppConfig(t *testing.T) {
	os.Setenv("APP_HOST", "127.0.0.1")
	os.Setenv("APP_PORT", "8080")
	os.Setenv("GRPC_AD_PORT", "50051")
	os.Setenv("GRPC_PORT_PROFILE", "50053")
	os.Setenv("GRPC_STORAGE_PORT", "50052")
	os.Setenv("GRPC_PORT_SLOT", "50054")
	os.Setenv("IMG_PATH", "/data/images")
	defer func() {
		os.Unsetenv("APP_HOST")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("GRPC_AD_PORT")
		os.Unsetenv("GRPC_PORT_PROFILE")
		os.Unsetenv("GRPC_STORAGE_PORT")
		os.Unsetenv("GRPC_PORT_SLOT")
		os.Unsetenv("IMG_PATH")
	}()

	cfg := GetAppConfig()

	if cfg.Host != "127.0.0.1" {
		t.Errorf("expected Host '127.0.0.1', got %s", cfg.Host)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected Port '8080', got %s", cfg.Port)
	}
	if cfg.PortAD != "50051" {
		t.Errorf("expected PortAD '50051', got %s", cfg.PortAD)
	}
	if cfg.PortProfile != "50053" {
		t.Errorf("expected PortProfile '50053', got %s", cfg.PortProfile)
	}
	if cfg.PortStorage != "50052" {
		t.Errorf("expected PortStorage '50052', got %s", cfg.PortStorage)
	}
	if cfg.PortSlot != "50054" {
		t.Errorf("expected PortSlot '50054', got %s", cfg.PortSlot)
	}
	if cfg.ImgPath != "/data/images" {
		t.Errorf("expected ImgPath '/data/images', got %s", cfg.ImgPath)
	}
}

func TestGetPostgresConfig_EmptyEnv(t *testing.T) {
	os.Clearenv()

	cfg := GetPostgresConfig()

	if cfg.User != "" {
		t.Errorf("expected empty User, got %s", cfg.User)
	}
}
