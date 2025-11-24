package main

import (
	"2025_2_404/internal/service/ad/config"
	adpb "2025_2_404/protos/gen/go/ad"
	db "2025_2_404/internal/service/ad/connections"
	adhandler "2025_2_404/internal/delivery/grpc/ad"
	repo "2025_2_404/internal/service/ad/repository/postgres"
	usecase "2025_2_404/internal/service/ad/usecase/ad"
	"fmt"
	"log"
	"net"

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

	adHandler := adhandler.New(useCaseCfg)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.AppConfig.PortAD))
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	grpcServer := grpc.NewServer()

	adpb.RegisterAdServServer(grpcServer, adHandler)

	log.Println("Starting server on", fmt.Sprintf("%s:%s", config.AppConfig.Host, config.AppConfig.PortAD))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
