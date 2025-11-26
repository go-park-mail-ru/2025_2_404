package main

import (
	"2025_2_404/internal/delivery/grpc/interceptor"
	"2025_2_404/internal/delivery/grpc/slot"
	"2025_2_404/internal/service/slot/config"
	db "2025_2_404/internal/service/slot/connections"
	repo "2025_2_404/internal/service/slot/repository/postgres"
	usecase "2025_2_404/internal/service/slot/usecase/slot"
	slotpb "2025_2_404/protos/gen/go/slot"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	config := config.GetConfig()
	connCfg, err := db.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer connCfg.CloseAll()

	repoCfg := repo.New(connCfg.PostgresSQL)
	useCaseCfg := usecase.New(repoCfg)
	slotHandler := slot.New(useCaseCfg)
	
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