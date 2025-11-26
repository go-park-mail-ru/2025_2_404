package main

import (
	adhandler "2025_2_404/internal/delivery/grpc/ad"
	"2025_2_404/internal/delivery/grpc/interceptor"
	"2025_2_404/internal/service/ad/config"
	db "2025_2_404/internal/service/ad/connections"
	repo "2025_2_404/internal/service/ad/repository/postgres"
	usecase "2025_2_404/internal/service/ad/usecase/ad"
	// authProto "2025_2_404/protos/auth"
	adpb "2025_2_404/protos/gen/go/ad"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config := config.GetConfig()
	connCfg, err := db.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer connCfg.CloseAll()

	authServiceAddr := "auth:8077" 

	authConn, err := grpc.NewClient(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to Auth Service: %v", err)
	}
	defer authConn.Close()

	// authClient := authProto.NewAuthClient(authConn)

	repoCfg := repo.New(connCfg.PostgresSQL)
	useCaseCfg := usecase.New(repoCfg)
	// authInterceptor := interceptor.AuthInterceptor(authClient)
	authInterceptor, authConn := interceptor.InitAuthInterceptor()
    defer authConn.Close()
	adHandler := adhandler.New(useCaseCfg)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.AppConfig.PortAD))
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor))

	adpb.RegisterAdServServer(grpcServer, adHandler)

	log.Println("Starting server on", fmt.Sprintf("%s:%s", config.AppConfig.Host, config.AppConfig.PortAD))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

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
