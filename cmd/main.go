// Package main implements the API Gateway for the advertisement service platform.
package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	httphandler "2025_2_404/internal/delivery/http"
	"2025_2_404/internal/delivery/http/middleware"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbAd "2025_2_404/protos/gen/go/ad"
	pbAuth "2025_2_404/protos/gen/go/auth"
	pbProfile "2025_2_404/protos/gen/go/profile"
	slotpb "2025_2_404/protos/gen/go/slot"
	pbStorage "2025_2_404/protos/gen/go/storage"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String("time", a.Value.Time().Format("15:04:05.000"))
			}
			return a
		},
	}))

	slog.SetDefault(logger)

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
	storageAddr := os.Getenv("STORAGE_ADDR")
	if storageAddr == "" {
		storageAddr = "localhost:8078"
	}
	slotAddr := os.Getenv("SLOT_ADDR")
	if slotAddr == "" {
		slotAddr = "localhost:8081"
	}
	gatewayPort := os.Getenv("APP_PORT")
	if gatewayPort == "" {
		gatewayPort = "8080"
	}

	// --- gRPC соединения ---
	dialOpts := grpc.WithTransportCredentials(insecure.NewCredentials())

	connAuth, err := grpc.NewClient(authAddr, dialOpts)
	if err != nil {
		log.Fatalf("Failed to connect to Auth: %v", err)
	}
	defer func() {
		_ = connAuth.Close()
	}()
	authClient := pbAuth.NewAuthClient(connAuth)

	connProfile, err := grpc.NewClient(profileAddr, dialOpts)
	if err != nil {
		log.Fatalf("Failed to connect to Profile: %v", err)
	}
	defer func() {
		_ = connProfile.Close()
	}()
	profileClient := pbProfile.NewProfileClient(connProfile)

	connAd, err := grpc.NewClient(adAddr, dialOpts)
	if err != nil {
		log.Fatalf("Failed to connect to Ad: %v", err)
	}
	defer func() {
		_ = connAd.Close()
	}()
	adClient := pbAd.NewAdServClient(connAd)

	connStorage, err := grpc.NewClient(storageAddr, dialOpts)
	if err != nil {
		log.Fatalf("Failed to connect to Storage: %v", err)
	}
	defer func() {
		_ = connStorage.Close()
	}()
	storageClient := pbStorage.NewStorageClient(connStorage)

	connSlot, err := grpc.NewClient(slotAddr, dialOpts)
	if err != nil {
		log.Fatalf("Failed to connect to Slot: %v", err)
	}
	defer func() {
		_ = connSlot.Close()
	}()
	slotClient := slotpb.NewSlotServClient(connSlot)

	r := mux.NewRouter()

	authRouter := r.PathPrefix("/api/auth").Subrouter()
	profileRouter := r.PathPrefix("/api/profile").Subrouter()
	adRouter := r.PathPrefix("/api/ads").Subrouter()
	slotRouter := r.PathPrefix("/api/slots").Subrouter()
	balanceRouter := r.PathPrefix("/api/balance").Subrouter()
	metricRouter := r.PathPrefix("/api/metric").Subrouter()

	// --- HTTP Handlers ---

	authHandler := httphandler.NewAuthHandler(authClient)
	profileHandler := httphandler.NewProfileHandler(profileClient, storageClient, adClient)
	adHandler := httphandler.NewAdHandler(adClient, storageClient, profileClient)
	slotHandler := httphandler.NewSlotHandler(slotClient, adClient, storageClient)

	// Auth
	authRouter.HandleFunc("/register", authHandler.Register).Methods("POST")
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Profile
	profileRouter.HandleFunc("", profileHandler.Show).Methods("GET")
	profileRouter.HandleFunc("/update", profileHandler.Update).Methods("POST")
	profileRouter.HandleFunc("", profileHandler.Delete).Methods("DELETE")

	// Balance
	balanceRouter.HandleFunc("", profileHandler.ShowBalance).Methods("GET")
	balanceRouter.HandleFunc("/subtract", profileHandler.SubtractBalance).Methods("POST")
	balanceRouter.HandleFunc("/payment", profileHandler.CreatePayment).Methods("POST")
	balanceRouter.HandleFunc("/payment/status", profileHandler.HandleYooKassaWebhook).Methods("POST")

	// Ads
	adRouter.HandleFunc("", adHandler.Create).Methods("POST")
	adRouter.HandleFunc("", adHandler.GetAll).Methods("GET")
	adRouter.HandleFunc("/{id}", adHandler.GetOne).Methods("GET")
	adRouter.HandleFunc("/{id}", adHandler.Update).Methods("PUT")
	adRouter.HandleFunc("/{id}", adHandler.Delete).Methods("DELETE")
	adRouter.HandleFunc("/{id}/addBudget", adHandler.UpdateBudget).Methods("PUT")

	// Slots
	slotRouter.HandleFunc("/serving/{id}", slotHandler.ServeSlot).Methods("GET")
	slotRouter.HandleFunc("", slotHandler.Create).Methods("POST")
	slotRouter.HandleFunc("", slotHandler.GetAll).Methods("GET")
	slotRouter.HandleFunc("/{id}", slotHandler.GetOne).Methods("GET")
	slotRouter.HandleFunc("/{id}", slotHandler.Update).Methods("PUT")
	slotRouter.HandleFunc("/{id}", slotHandler.Delete).Methods("DELETE")
	slotRouter.HandleFunc("/{id}/statistics", slotHandler.GetMetrics).Methods("GET")

	//Metric
	metricRouter.HandleFunc("", slotHandler.CreateMetric).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())

	handler := middleware.CorsMiddleware(r)
	handler = middleware.AccessLogMiddleware(handler)
	slog.Info("API Gateway running on " + gatewayPort)
	srv := &http.Server{
		Addr:         ":" + gatewayPort,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to run gateway: %v", err)
	}
}
