package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"2025_2_404/internal/delivery/grpc/ad"
	"2025_2_404/internal/delivery/grpc/auth/handler"
	"2025_2_404/internal/delivery/grpc/profile/handler"
	"2025_2_404/internal/delivery/grpc/storage"

	adv1 "2025_2_404/protos/gen/go/ad"
	authv1 "2025_2_404/protos/auth"
	profilev1 "2025_2_404/protos/profile"
	storagev1 "2025_2_404/protos/gen/go/storage"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()


	grpcConn, err := grpc.DialContext(ctx, "localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial gRPC server: %v", err)
	}
	defer grpcConn.Close()

	gwmux := runtime.NewServeMux()

	err = adv1.RegisterAdServServer(ctx, gwmux, grpcConn)
	if err != nil {
		log.Fatalf("Failed to register AdServ gateway: %v", err)
	}

	err = authv1.RegisterAuthHandler(ctx, gwmux, grpcConn)
	if err != nil {
		log.Fatalf("Failed to register Auth gateway: %v", err)
	}

	err = profilev1.RegisterProfileHandler(ctx, gwmux, grpcConn)
	if err != nil {
		log.Fatalf("Failed to register Profile gateway: %v", err)
	}

	err = storagev1.RegisterStorageHandler(ctx, gwmux, grpcConn)
	if err != nil {
		log.Fatalf("Failed to register Storage gateway: %v", err)
	}

	// HTTP-сервер
	httpSrv := &http.Server{
		Addr:         ":8080",
		Handler:      gwmux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting HTTP gateway on :8080")
	log.Fatal(httpSrv.ListenAndServe())
}

// package main

// import (
// 	"log"

// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"

// 	profileHttp "2025_2_404/internal/delivery/http/profile"
// 	authHttp "2025_2_404/internal/delivery/http/auth"
// 	pbAuth "2025_2_404/protos/auth"
// 	pbProfile "2025_2_404/protos/profile"
// )

// func main() {
// 	authAddr := "localhost:50001"    // Auth Service
// 	profileAddr := "localhost:50002" // Profile Service
// 	gatewayPort := ":8080"          // API Gateway

// 	connAuth, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("Failed to connect to Auth Service: %v", err)
// 	}
// 	defer connAuth.Close()
// 	authClient := pbAuth.NewAuthClient(connAuth) 

// 	connProfile, err := grpc.NewClient(profileAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("Failed to connect to Profile Service: %v", err)
// 	}
// 	defer connProfile.Close()
// 	profileClient := pbProfile.NewProfileClient(connProfile)

// 	r := gin.Default()
// 	r.Use(corsMiddleware())

// 	authHandler := authHttp.NewAuthHandler(authClient)
// 	authHandler.RegisterRoutes(r)

// 	profileHandler := profileHttp.NewProfileHandler(profileClient)
// 	profileHandler.RegisterRoutes(r)

// 	log.Printf("API Gateway running on %s", gatewayPort)
// 	if err := r.Run(gatewayPort); err != nil {
// 		log.Fatalf("Failed to run gateway: %v", err)
// 	}
// }

// func corsMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") 
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	}
// }

// package main

// import (
// 	"2025_2_404/internal/config"
// 	db "2025_2_404/internal/connections"
// 	adhandler "2025_2_404/internal/delivery/http/adhandler"
// 	authhandler "2025_2_404/internal/delivery/http/authhandler"
// 	balancehandler "2025_2_404/internal/delivery/http/balancehandler"
// 	feedhandler "2025_2_404/internal/delivery/http/feedhandler"
// 	middleware "2025_2_404/internal/delivery/http/middleware"
// 	"2025_2_404/internal/delivery/http/profilehandler"
// 	repo "2025_2_404/internal/repository/postgres"
// 	usecase "2025_2_404/internal/use_case"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gorilla/mux"
// )

// const(
// 	Timeout = time.Second * 5
// )

// func main() {
// 	config := config.GetConfig()
// 	connCfg, err := db.New(config)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer connCfg.CloseAll()
// 	repoCfg := repo.New(connCfg)
// 	useCaseCfg := usecase.New(config, repoCfg)
	
// 	middle := middleware.New(useCaseCfg.TokenUsecase)
// 	handlersAd := adhandler.New(useCaseCfg.AdUsecase)
// 	handlersAuth := authhandler.New(useCaseCfg.AuthUsecase)
// 	handlersProfile := profilehandler.New(useCaseCfg.ProfileUsecase)
// 	handlersBalance := balancehandler.New(useCaseCfg.BalanceUsecase)
// 	handlersFeed := feedhandler.New(useCaseCfg.FeedUsecase)

	
// 	mainRouter := mux.NewRouter()
// 	authSubrouter := mainRouter.PathPrefix("/auth").Subrouter()
// 	adSubrouter := mainRouter.PathPrefix("/ads").Subrouter()
// 	clientSubroute := mainRouter.PathPrefix("/profile").Subrouter()
// 	balanceSubrouter := mainRouter.PathPrefix("/wallet").Subrouter()
// 	feedSurouter := mainRouter.PathPrefix("/feed").Subrouter()

// 	authSubrouter.HandleFunc("/signup", handlersAuth.RegisterHandler).Methods(http.MethodPost, http.MethodOptions)
// 	authSubrouter.HandleFunc("/signin", handlersAuth.LoginHandler).Methods(http.MethodPost, http.MethodOptions)
// 	authSubrouter.Use(middle.Peflite)

// 	adSubrouter.HandleFunc("/", handlersAd.Handler).Methods(http.MethodGet, http.MethodOptions)
// 	adSubrouter.HandleFunc("/", handlersAd.CreateHandler).Methods(http.MethodPost, http.MethodOptions)
// 	adSubrouter.HandleFunc("/{ad_id}", handlersAd.UpdateHandler).Methods(http.MethodPut, http.MethodOptions)
// 	adSubrouter.HandleFunc("/{ad_id}", handlersAd.DeleteHandler).Methods(http.MethodDelete, http.MethodOptions)
// 	adSubrouter.HandleFunc("/{ad_id}", handlersAd.GetOneAd).Methods(http.MethodGet, http.MethodOptions)
// 	adSubrouter.Use(middle.Peflite, middle.Auth)

// 	clientSubroute.HandleFunc("/", handlersProfile.ShowHandler).Methods(http.MethodGet, http.MethodOptions)
// 	clientSubroute.HandleFunc("/", handlersProfile.UpdateHandler).Methods(http.MethodPut, http.MethodOptions)
// 	clientSubroute.HandleFunc("/", handlersProfile.DeleteHandler).Methods(http.MethodDelete, http.MethodOptions)
// 	clientSubroute.Use(middle.Peflite, middle.Auth)

// 	balanceSubrouter.HandleFunc("/", handlersBalance.Show).Methods(http.MethodGet, http.MethodOptions)
// 	balanceSubrouter.Use(middle.Peflite, middle.Auth)

// 	feedSurouter.HandleFunc("/{platform_name}", handlersFeed.GetAdFeedHandler).Methods(http.MethodGet, http.MethodOptions)
// 	feedSurouter.Use(middle.Peflite)

// 	srv := &http.Server{
//         Addr:         fmt.Sprintf("%s:%s", config.AppConfig.Host, config.AppConfig.Port),
//         WriteTimeout: Timeout,
//         ReadTimeout:  Timeout,
//         IdleTimeout:  Timeout,
//         Handler: mainRouter,
//     }
// 	log.Println("Starting server on", fmt.Sprintf("%s:%s", config.AppConfig.Host, config.AppConfig.Port))
// 	err = srv.ListenAndServe()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

