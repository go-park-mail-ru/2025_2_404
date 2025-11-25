package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	
	gatewayHttp "2025_2_404/internal/delivery/http"
	pbAuth "2025_2_404/protos/auth"
	pbAd "2025_2_404/protos/gen/go/ad"
	pbProfile "2025_2_404/protos/profile"
	pbStorage "2025_2_404/protos/gen/go/storage" 
)

func main() {
	authAddr := os.Getenv("AUTH_ADDR")
	if authAddr == "" {
		authAddr = "localhost:8077" 
	}

	profileAddr := os.Getenv("PROFILE_ADDR")
	if profileAddr == "" {
		profileAddr = "localhost:8076"
	}

	adAddr := os.Getenv("AD_ADDR")
	if adAddr == "" {
		adAddr = "localhost:8079"
	}
	
	StoragePort := os.Getenv("STORAGE_ADDR")
	if adAddr == "" {
		adAddr = "localhost:8078"
	}

	gatewayPort := os.Getenv("APP_PORT")
	if gatewayPort == "" {
		gatewayPort = "8080"
	}																																																							

	connAuth, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Auth: %v", err)
	}
	defer connAuth.Close()
	authClient := pbAuth.NewAuthClient(connAuth)

	connProfile, err := grpc.NewClient(profileAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Profile: %v", err)
	}
	defer connProfile.Close()
	profileClient := pbProfile.NewProfileClient(connProfile)

	connAd, err := grpc.NewClient(adAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Ad: %v", err)
	}
	defer connAd.Close()
	adClient := pbAd.NewAdServClient(connAd)

	connStorage, err := grpc.NewClient(StoragePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Storage: %v", err)
	}
	defer connStorage.Close()
	storageClient := pbStorage.NewStorageClient(connStorage)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	authHandler := gatewayHttp.NewAuthHandler(authClient)
	authHandler.RegisterRoutes(r)

	profileHandler := gatewayHttp.NewProfileHandler(profileClient)
	profileHandler.RegisterRoutes(r)

	adHandler := gatewayHttp.NewAdHandler(adClient)
	adHandler.RegisterRoutes(r)

	storageHandler := gatewayHttp.NewStorageHandler(storageClient)
	storageHandler.RegisterRoutes(r)

	//  r.Run(":" + gatewayPort)
	log.Printf("API Gateway running on %s", gatewayPort)
	if err := r.Run(":" + gatewayPort); err != nil {
		log.Fatalf("Failed to run gateway: %v", err)
	}
}