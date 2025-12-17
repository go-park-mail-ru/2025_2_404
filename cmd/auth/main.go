// Package main initializes and starts the gRPC authentication service.
package main

import (
	handler "2025_2_404/internal/delivery/grpc/auth"
	"2025_2_404/internal/service/auth/config"
	"2025_2_404/internal/service/auth/connections"
	"2025_2_404/internal/service/auth/service"
	jwt "2025_2_404/internal/service/auth/service"
	"2025_2_404/internal/service/auth/storage/postgres"
	"2025_2_404/protos/gen/go/auth"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.GetConfig()

	db, err := connections.New(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.CloseAll()

	go func() {
		log.Println("Starting metrics server on :9090")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Printf("Metrics server failed: %v", err)
		}
	}()

	userRepo := postgres.New(db.PostgresSQL)
	jwtUseCase := jwt.NewJWT(cfg.AppConfig.JwtPrivateKey, cfg.AppConfig.JwtPublicKey)
	authUseCase := service.New(userRepo, jwtUseCase)

	authServer := handler.NewAuthServer(authUseCase, jwtUseCase)

	grpcServer := grpc.NewServer()

	auth.RegisterAuthServer(grpcServer, authServer)
	lis, err := net.Listen("tcp", ":"+cfg.AppConfig.Port)
	if err != nil {
		log.Fatalf("failed to listen on :%s: %v", cfg.AppConfig.Port, err)
	}

	go func() {
		log.Printf("gRPC Auth Service running on :%s", cfg.AppConfig.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()
}
