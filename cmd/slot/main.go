package main

import (
	"2025_2_404/internal/delivery/grpc/interceptor"
	"2025_2_404/internal/delivery/grpc/slot"
	"2025_2_404/internal/service/slot/config"
	db "2025_2_404/internal/service/slot/connections"
	metricRepo "2025_2_404/internal/service/slot/repository/postgres/metric"
	repo "2025_2_404/internal/service/slot/repository/postgres/slot"
	metricusecase "2025_2_404/internal/service/slot/usecase/metric"
	usecase "2025_2_404/internal/service/slot/usecase/slot"
	slotpb "2025_2_404/protos/gen/go/slot"
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

	adServiceAddr := fmt.Sprintf("ad-service:%s", config.AppConfig.PortAD)
	adConn, err := grpc.NewClient(adServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to Ad Service: %v", err)
	}
	defer adConn.Close()

	adClient := adpb.NewAdServClient(adConn)

	repoCfg := repo.New(connCfg.PostgresSQL)
	metricRepo := metricRepo.New(connCfg.PostgresSQL)
	useCaseCfg := usecase.New(repoCfg)
	metricUsecase := metricusecase.New(metricRepo)
	slotHandler := slot.New(useCaseCfg, metricUsecase, adClient)
	
	authInterceptor, authConn := interceptor.InitAuthInterceptor()
    defer authConn.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.AppConfig.PortSlot))
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(authInterceptor),
    )

	slotpb.RegisterSlotServServer(grpcServer, slotHandler)

	log.Println("Starting Slot Server on", fmt.Sprintf("%s:%s", config.AppConfig.Host, config.AppConfig.PortSlot))
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
	log.Println("Slot Server stopped")
}