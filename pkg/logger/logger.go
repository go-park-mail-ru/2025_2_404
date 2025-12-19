// Package logger provides a zap-based structured logger.
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.Logger, error) {
	env := os.Getenv("ENV")
	if env == "development" || env == "dev" {
		return zap.NewDevelopment()
	}

	return zap.NewProduction(
		zap.AddStacktrace(zapcore.FatalLevel),
	)
}
