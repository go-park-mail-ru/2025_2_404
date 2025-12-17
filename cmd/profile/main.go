// Package main initializes and starts the gRPC profile service.
package main

import (
	"2025_2_404/internal/delivery/grpc/interceptor"
	handler "2025_2_404/internal/delivery/grpc/profile"
	"2025_2_404/internal/service/profile/config"
	"2025_2_404/internal/service/profile/connections"
	external "2025_2_404/internal/service/profile/external/http"
	service "2025_2_404/internal/service/profile/service"
	repository "2025_2_404/internal/service/profile/storage/postgres"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// authProto "2025_2_404/protos/auth"
	pb "2025_2_404/protos/gen/go/profile"
)

func main() {
	cfg := config.GetConfig()

	db, err := connections.New(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.CloseAll()

	authServiceAddr := "auth:8077"

	authConn, err := grpc.NewClient(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to Auth Service: %v", err)
	}
	defer func() {
		_ = authConn.Close()
	}()

	go func() {
		log.Println("Starting metrics server on :9090")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Printf("Metrics server failed: %v", err)
		}
	}()

	// authClient := authProto.NewAuthClient(authConn)

	userRepo := repository.New(db.PostgresSQL)
	yooExternal := external.New(cfg)
	authInterceptor, authConn := interceptor.InitAuthInterceptor()
	profileUseCase := service.New(userRepo, yooExternal)
	profileServer := handler.NewProfileServer(profileUseCase)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
	)

	pb.RegisterProfileServer(grpcServer, profileServer)
	lis, err := net.Listen("tcp", ":"+cfg.AppConfig.Port)
	if err != nil {
		log.Fatalf("failed to listen on :%s: %v", cfg.AppConfig.Port, err)
	}

	log.Printf("gRPC Profile Service starting on :%s", cfg.AppConfig.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}
