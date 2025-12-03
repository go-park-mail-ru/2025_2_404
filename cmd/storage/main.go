package main

import (
	"log/slog"
	"os"
	"2025_2_404/internal/delivery/grpc/interceptor"
	"2025_2_404/internal/service/storage/config"
	storagepb "2025_2_404/protos/gen/go/storage"
	storagehandler "2025_2_404/internal/delivery/grpc/storage"
	usecase "2025_2_404/internal/service/storage/usecase/filestorage"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// === Настройка логгера для разработки ===
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String("time", a.Value.Time().Format("15:04:05.000"))
			}
			return a
		},
	}))
	slog.SetDefault(logger)

	// === Инициализация сервиса ===
	cfg := config.GetConfig()
	useCase := usecase.New(cfg)
	storageHandler := storagehandler.New(useCase)

	// === gRPC сервер с interceptor'ом ===
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.GrpcLoggerInterceptor),
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.AppConfig.PortStorage))
	if err != nil {
		slog.Error("Failed to listen", "port", cfg.AppConfig.PortStorage, "error", err)
		os.Exit(1)
	}

	storagepb.RegisterStorageServer(grpcServer, storageHandler)

	slog.Info("Storage gRPC server started",
		"host", cfg.AppConfig.Host,
		"port", cfg.AppConfig.PortStorage)

	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("gRPC server failed", "error", err)
		os.Exit(1)
	}
}