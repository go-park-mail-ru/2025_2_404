// Package main initializes and starts the gRPC storage service.
package main

import (
	"2025_2_404/internal/delivery/grpc/interceptor"
	storagehandler "2025_2_404/internal/delivery/grpc/storage"
	"2025_2_404/internal/service/storage/config"
	usecase "2025_2_404/internal/service/storage/usecase/filestorage"
	storagepb "2025_2_404/protos/gen/go/storage"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
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

	go func() {
		log.Println("Starting metrics server on :9090")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Printf("Metrics server failed: %v", err)
		}
	}()

	cfg := config.GetConfig()
	useCase := usecase.New(cfg)
	storageHandler := storagehandler.New(useCase)

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
