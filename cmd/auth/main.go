package main

import (
	handler "2025_2_404/internal/delivery/grpc/auth"
	jwt "2025_2_404/internal/service/auth/service"
	"2025_2_404/internal/service/auth/config"
	"2025_2_404/internal/service/auth/connections"
	"2025_2_404/internal/service/auth/storage/postgres"
	"2025_2_404/internal/service/auth/service"
	"2025_2_404/protos/auth"
	"net"
	"os"
	"os/signal"
	"syscall"

	"log"
	"google.golang.org/grpc"
)

func main(){
	cfg := config.GetConfig()

	db, err := connections.New(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.CloseAll()

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

	gracefulShutdown(grpcServer, lis)
}

func gracefulShutdown(grpcServer *grpc.Server, lis net.Listener) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down gracefully...")
	grpcServer.GracefulStop()
	lis.Close()
	log.Println("Server stopped")
}