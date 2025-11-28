package main

import (
	"log"
	"net/http"
	"os"
	"time"

	httphandler "2025_2_404/internal/delivery/http"
    "2025_2_404/internal/delivery/http/middleware"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	slotpb "2025_2_404/protos/gen/go/slot"
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
	defer connAuth.Close()
	authClient := pbAuth.NewAuthClient(connAuth)

	connProfile, err := grpc.NewClient(profileAddr, dialOpts)
	if err != nil {
		log.Fatalf("Failed to connect to Profile: %v", err)
	}
	defer connProfile.Close()
	profileClient := pbProfile.NewProfileClient(connProfile)

	connAd, err := grpc.NewClient(adAddr, dialOpts)
	if err != nil {
		log.Fatalf("Failed to connect to Ad: %v", err)
	}
	defer connAd.Close()
	adClient := pbAd.NewAdServClient(connAd)

	connStorage, err := grpc.NewClient(storageAddr, dialOpts)
	if err != nil {
		log.Fatalf("Failed to connect to Storage: %v", err)
	}
	defer connStorage.Close()
	storageClient := pbStorage.NewStorageClient(connStorage)

	connSlot, err := grpc.NewClient(slotAddr, dialOpts)
	if err != nil {
		log.Fatalf("Failed to connect to Slot: %v", err)
	}
	defer connSlot.Close()
	slotClient := slotpb.NewSlotServClient(connSlot)

	r := mux.NewRouter()

	authRouter := r.PathPrefix("/auth").Subrouter()
	profileRouter := r.PathPrefix("/profile").Subrouter()
	adRouter := r.PathPrefix("/ads").Subrouter()
	slotRouter := r.PathPrefix("/slots").Subrouter()
	balanceRouter := r.PathPrefix("/balance").Subrouter()

	// --- HTTP Handlers ---

	authHandler := httphandler.NewAuthHandler(authClient)
	profileHandler := httphandler.NewProfileHandler(profileClient)
	adHandler := httphandler.NewAdHandler(adClient, storageClient)
	slotHandler := httphandler.NewSlotHandler(slotClient, adClient)

	// Auth
	authRouter.HandleFunc("/register", authHandler.Register).Methods("POST")
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Profile
	profileRouter.HandleFunc("", profileHandler.Show).Methods("GET")
	profileRouter.HandleFunc("/update", profileHandler.Update).Methods("POST")
	profileRouter.HandleFunc("", profileHandler.Delete).Methods("DELETE")

	// Balance
	balanceRouter.HandleFunc("", profileHandler.ShowBalance).Methods("GET")
	balanceRouter.HandleFunc("/add", profileHandler.AddBalance).Methods("POST")
	balanceRouter.HandleFunc("/subtract", profileHandler.SubtractBalance).Methods("POST")

	// Ads
	adRouter.HandleFunc("", adHandler.Create).Methods("POST")
	adRouter.HandleFunc("", adHandler.GetAll).Methods("GET")
	adRouter.HandleFunc("/{id}", adHandler.GetOne).Methods("GET")
	adRouter.HandleFunc("/{id}", adHandler.Update).Methods("PUT")
	adRouter.HandleFunc("/{id}", adHandler.Delete).Methods("DELETE")

	// Slots
	slotRouter.HandleFunc("/serving/{id}", slotHandler.ServeSlot).Methods("GET")
	slotRouter.HandleFunc("", slotHandler.Create).Methods("POST")
	slotRouter.HandleFunc("", slotHandler.GetAll).Methods("GET")
	slotRouter.HandleFunc("/{id}", slotHandler.GetOne).Methods("GET")
	slotRouter.HandleFunc("/{id}", slotHandler.Update).Methods("PUT")
	slotRouter.HandleFunc("/{id}", slotHandler.Delete).Methods("DELETE")

	handler := middleware.CorsMiddleware(r)
	log.Printf("API Gateway running on %s", gatewayPort)
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