package logger

import (
	"os"
	"testing"
)

func TestNew_Production(t *testing.T) {
	// Clear ENV or set to production
	originalEnv := os.Getenv("ENV")
	os.Setenv("ENV", "production")
	defer os.Setenv("ENV", originalEnv)

	logger, err := New()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if logger == nil {
		t.Fatal("expected logger to be created, got nil")
	}

	// Cleanup
	_ = logger.Sync()
}

func TestNew_Development(t *testing.T) {
	originalEnv := os.Getenv("ENV")
	os.Setenv("ENV", "development")
	defer os.Setenv("ENV", originalEnv)

	logger, err := New()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if logger == nil {
		t.Fatal("expected logger to be created, got nil")
	}

	// Cleanup
	_ = logger.Sync()
}

func TestNew_Dev(t *testing.T) {
	originalEnv := os.Getenv("ENV")
	os.Setenv("ENV", "dev")
	defer os.Setenv("ENV", originalEnv)

	logger, err := New()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if logger == nil {
		t.Fatal("expected logger to be created, got nil")
	}

	// Cleanup
	_ = logger.Sync()
}

func TestNew_DefaultProduction(t *testing.T) {
	originalEnv := os.Getenv("ENV")
	os.Unsetenv("ENV")
	defer os.Setenv("ENV", originalEnv)

	logger, err := New()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if logger == nil {
		t.Fatal("expected logger to be created, got nil")
	}

	// Test that logger can be used
	logger.Info("test message")

	// Cleanup
	_ = logger.Sync()
}

func TestNew_LoggerWorks(t *testing.T) {
	logger, err := New()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Test that logging doesn't panic
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")

	// If we reach here, logger is working
}
