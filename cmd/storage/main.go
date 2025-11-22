package storage

import (
	"2025_2_404/internal/config"
	db "2025_2_404/internal/connections"
	storagepb "2025_2_404/protos/gen/go/storage"
	storagehandler "2025_2_404/internal/delivery/grpc/storage"
	repo "2025_2_404/internal/repository/postgres"
	usecase "2025_2_404/internal/use_case"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main(){
	config := config.GetConfig()
	connCfg, err := db.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer connCfg.CloseAll()

	repoCfg := repo.New(connCfg)
	useCaseCfg := usecase.New(config, repoCfg)
	storageHandler := storagehandler.New(useCaseCfg.StorageUsecase)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.AppConfig.PortStorage))
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	grpcServer := grpc.NewServer()

	storagepb.RegisterStorageServer(grpcServer, storageHandler)

	log.Println("Starting server on", fmt.Sprintf("%s:%s", config.AppConfig.Host, config.AppConfig.PortStorage))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}