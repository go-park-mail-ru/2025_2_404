package main

import (
	"2025_2_404/internal/delivery/grpc/interceptor"
	handler "2025_2_404/internal/delivery/grpc/profile"
	"2025_2_404/internal/service/profile/config"
	"2025_2_404/internal/service/profile/connections"
	service "2025_2_404/internal/service/profile/service"
	repository "2025_2_404/internal/service/profile/storage/postgres"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	authProto "2025_2_404/protos/auth"
	pb "2025_2_404/protos/profile"
)

func main() {
	cfg := config.GetConfig()

	db, err := connections.New(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.CloseAll()

	authServiceAddr := "auth:50001" 

	authConn, err := grpc.NewClient(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to Auth Service: %v", err)
	}
	defer authConn.Close()

	authClient := authProto.NewAuthClient(authConn)

	userRepo := repository.New(db.PostgresSQL)
	authInterceptor := interceptor.AuthInterceptor(authClient)
	profileUseCase := service.New(userRepo)
	profileServer := handler.NewProfileServer(profileUseCase)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
	)

	pb.RegisterProfileServer(grpcServer, profileServer)
	lis, err := net.Listen("tcp", ":"+cfg.AppConfig.Port)
	if err != nil {
		log.Fatalf("failed to listen on :%s: %v", cfg.AppConfig.Port, err)
	}

	go func() {
		log.Printf("gRPC Profile Service running on :%s", cfg.AppConfig.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	gracefulShutdown(grpcServer, lis)
}

func gracefulShutdown(grpcServer *grpc.Server, lis net.Listener) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	sig := <-c
	log.Printf("Received signal %v. Shutting down gracefully...", sig)

	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
