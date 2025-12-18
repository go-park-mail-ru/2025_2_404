package config

import (
	"os"
	"testing"
)

func TestGetPostgresConfig(t *testing.T) {
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
	os.Setenv("GRPC_PORT_PROFILE", "50053")
	os.Setenv("IMG_PATH", "/var/images")
	defer func() {
		os.Unsetenv("GRPC_PORT_PROFILE")
		os.Unsetenv("IMG_PATH")
	}()

	cfg := GetAppConfig()

	if cfg.Port != "50053" {
		t.Errorf("expected Port '50053', got %s", cfg.Port)
	}
	if cfg.ImgPath != "/var/images" {
		t.Errorf("expected ImgPath '/var/images', got %s", cfg.ImgPath)
	}
}

func TestGetPaymentConfig(t *testing.T) {
	os.Setenv("YOOKASSA_SHOP_ID", "123456")
	os.Setenv("YOOKASSA_SECRET_KEY", "secret123")
	defer func() {
		os.Unsetenv("YOOKASSA_SHOP_ID")
		os.Unsetenv("YOOKASSA_SECRET_KEY")
	}()

	cfg := GetPaymentConfig()

	if cfg.ShopID != "123456" {
		t.Errorf("expected ShopID '123456', got %s", cfg.ShopID)
	}
	if cfg.SecretKey != "secret123" {
		t.Errorf("expected SecretKey 'secret123', got %s", cfg.SecretKey)
	}
}

func TestGetPostgresConfig_EmptyEnv(t *testing.T) {
	os.Clearenv()

	cfg := GetPostgresConfig()

	if cfg.User != "" {
		t.Errorf("expected empty User, got %s", cfg.User)
	}
	if cfg.Password != "" {
		t.Errorf("expected empty Password, got %s", cfg.Password)
	}
}
