package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

// GrpcLoggerInterceptor — логирует все unary gRPC вызовы
func GrpcLoggerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	slog.Debug("gRPC request received",
		"method", info.FullMethod,
		"request", req,
	)

	// Выполняем хендлер
	resp, err := handler(ctx, req)

	durationMs := float64(time.Since(start).Microseconds()) / 1000.0

	if err != nil {
		slog.Error("gRPC request failed",
			"method", info.FullMethod,
			"error", err,
			"took_ms", durationMs,
		)
	} else {
		slog.Info("gRPC request succeeded",
			"method", info.FullMethod,
			"took_ms", durationMs,
		)
	}

	return resp, err
}
