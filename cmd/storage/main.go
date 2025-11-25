package storage

import (
	"2025_2_404/internal/service/storage/config"
	storagepb "2025_2_404/protos/gen/go/storage"
	storagehandler "2025_2_404/internal/delivery/grpc/storage"
	usecase "2025_2_404/internal/service/storage/usecase/filestorage"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main(){
	config := config.GetConfig()
	useCaseCfg := usecase.New(config)
	storageHandler := storagehandler.New(useCaseCfg)

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