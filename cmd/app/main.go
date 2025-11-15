package main

import (
	"2025_2_404/internal/config"
	db "2025_2_404/internal/connections"
	adhandler "2025_2_404/internal/delivery/http/adhandler"
	authhandler "2025_2_404/internal/delivery/http/authhandler"
	balancehandler "2025_2_404/internal/delivery/http/balancehandler"
	feedhandler "2025_2_404/internal/delivery/http/feedhandler"
	middleware "2025_2_404/internal/delivery/http/middleware"
	"2025_2_404/internal/delivery/http/profilehandler"
	repo "2025_2_404/internal/repository/postgres"
	usecase "2025_2_404/internal/use_case"
	"fmt"
	supporthandler "2025_2_404/internal/delivery/http/supporthandler"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const(
	Timeout = time.Second * 5
)

func main() {
	config := config.GetConfig()
	connCfg, err := db.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer connCfg.CloseAll()
	repoCfg := repo.New(connCfg)
	useCaseCfg := usecase.New(config, repoCfg)
	
	middle := middleware.New(useCaseCfg.TokenUsecase)
	handlersAd := adhandler.New(useCaseCfg.AdUsecase)
	handlersAuth := authhandler.New(useCaseCfg.AuthUsecase)
	handlersProfile := profilehandler.New(useCaseCfg.ProfileUsecase)
	handlersBalance := balancehandler.New(useCaseCfg.BalanceUsecase)
	handlersFeed := feedhandler.New(useCaseCfg.FeedUsecase)
	handlersSupport := supporthandler.New(useCaseCfg.SupportUsecase)

	
	mainRouter := mux.NewRouter()
	authSubrouter := mainRouter.PathPrefix("/auth").Subrouter()
	adSubrouter := mainRouter.PathPrefix("/ads").Subrouter()
	clientSubroute := mainRouter.PathPrefix("/profile").Subrouter()
	balanceSubrouter := mainRouter.PathPrefix("/wallet").Subrouter()
	feedSurouter := mainRouter.PathPrefix("/feed").Subrouter()
	supportSubrouter := mainRouter.PathPrefix("/support").Subrouter()
	supportsSubrouter := mainRouter.PathPrefix("/supports").Subrouter()

	authSubrouter.HandleFunc("/signup", handlersAuth.RegisterHandler).Methods(http.MethodPost, http.MethodOptions)
	authSubrouter.HandleFunc("/signin", handlersAuth.LoginHandler).Methods(http.MethodPost, http.MethodOptions)
	authSubrouter.Use(middle.Peflite)

	adSubrouter.HandleFunc("/", handlersAd.Handler).Methods(http.MethodGet, http.MethodOptions)
	adSubrouter.HandleFunc("/", handlersAd.CreateHandler).Methods(http.MethodPost, http.MethodOptions)
	adSubrouter.HandleFunc("/{ad_id}", handlersAd.UpdateHandler).Methods(http.MethodPut, http.MethodOptions)
	adSubrouter.HandleFunc("/{ad_id}", handlersAd.DeleteHandler).Methods(http.MethodDelete, http.MethodOptions)
	adSubrouter.HandleFunc("/{ad_id}", handlersAd.GetOneAd).Methods(http.MethodGet, http.MethodOptions)
	adSubrouter.Use(middle.Peflite, middle.Auth)

	clientSubroute.HandleFunc("/", handlersProfile.ShowHandler).Methods(http.MethodGet, http.MethodOptions)
	clientSubroute.HandleFunc("/", handlersProfile.UpdateHandler).Methods(http.MethodPut, http.MethodOptions)
	clientSubroute.HandleFunc("/", handlersProfile.DeleteHandler).Methods(http.MethodDelete, http.MethodOptions)
	clientSubroute.Use(middle.Peflite, middle.Auth)

	balanceSubrouter.HandleFunc("/", handlersBalance.Show).Methods(http.MethodGet, http.MethodOptions)
	balanceSubrouter.Use(middle.Peflite, middle.Auth)

	feedSurouter.HandleFunc("/{platform_name}", handlersFeed.GetAdFeedHandler).Methods(http.MethodGet, http.MethodOptions)
	feedSurouter.Use(middle.Peflite)

	supportSubrouter.HandleFunc("/", handlersSupport.Handler).Methods(http.MethodGet, http.MethodOptions)
	
	supportSubrouter.HandleFunc("/{uved_id}", handlersSupport.GetOneUved).Methods(http.MethodGet, http.MethodOptions)
	supportSubrouter.HandleFunc("/", handlersSupport.CreateHandler).Methods(http.MethodPost, http.MethodOptions)
	supportSubrouter.HandleFunc("/{uved_id}", handlersSupport.UpdateHandler).Methods(http.MethodPut, http.MethodOptions)
	supportSubrouter.HandleFunc("/{uved_id}", handlersSupport.DeleteHandler).Methods(http.MethodDelete, http.MethodOptions)
	supportSubrouter.Use(middle.Peflite, middle.Auth)

	supportsSubrouter.HandleFunc("/", handlersSupport.HandlerGetAll).Methods(http.MethodGet, http.MethodOptions)
	supportsSubrouter.Use(middle.Peflite, middle.Auth)


	srv := &http.Server{
        Addr:         fmt.Sprintf("%s:%s", config.AppConfig.Host, config.AppConfig.Port),
        WriteTimeout: Timeout,
        ReadTimeout:  Timeout,
        IdleTimeout:  Timeout,
        Handler: mainRouter,
    }
	log.Println("Starting server on", fmt.Sprintf("%s:%s", config.AppConfig.Host, config.AppConfig.Port))
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

