// Package main initializes and starts the gRPC authentication service.
package main

import (
	handler "2025_2_404/internal/delivery/grpc/auth"
	"2025_2_404/internal/service/auth/config"
	"2025_2_404/internal/service/auth/connections"
	"2025_2_404/internal/service/auth/service"
	jwt "2025_2_404/internal/service/auth/service"
	"2025_2_404/internal/service/auth/storage/postgres"
	"2025_2_404/pkg/logger"
	"2025_2_404/protos/gen/go/auth"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	log, err := logger.New()
	if err != nil {
		panic("failed to initialize logger")
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("starting Auth gRPC server")
	cfg := config.GetConfig()

	db, err := connections.New(cfg)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.CloseAll()

	go func() {
		log.Info("starting metrics server", zap.String("addr", ":9090"))
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Error("metrics server failed", zap.Error(err))
		}
	}()

	userRepo := postgres.New(db.PostgresSQL, log.Named("repo.user"))
	jwtUseCase := jwt.NewJWT(cfg.AppConfig.JwtPrivateKey, cfg.AppConfig.JwtPublicKey, log.Named("jwt.usecase"))
	authUseCase := service.New(userRepo, jwtUseCase, log.Named("auth.usecase"))

	authServer := handler.NewAuthServer(authUseCase, jwtUseCase, log.Named("handler.auth"))

	grpcServer := grpc.NewServer()

	auth.RegisterAuthServer(grpcServer, authServer)
	lis, err := net.Listen("tcp", ":"+cfg.AppConfig.Port)
	if err != nil {
		log.Fatal("failed to listen on port", zap.String("port", cfg.AppConfig.Port), zap.Error(err))
	}

	log.Info("gRPC Auth Service starting", zap.String("port", cfg.AppConfig.Port))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("gRPC server failed", zap.Error(err))
	}
}
