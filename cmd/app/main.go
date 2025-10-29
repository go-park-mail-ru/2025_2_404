package main

import (
	"2025_2_404/internal/config"
	db "2025_2_404/internal/connections"
	adhandler "2025_2_404/internal/delivery/http/adhandler"
	authhandler "2025_2_404/internal/delivery/http/authhandler"
	middleware "2025_2_404/internal/delivery/http/middleware"
	repo "2025_2_404/internal/repository/postgres"
	usecase "2025_2_404/internal/use_case"
	"time"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"log"
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

	
	mainRouter := mux.NewRouter()
	authSubrouter := mainRouter.PathPrefix("/auth").Subrouter()
	adSubrouter := mainRouter.PathPrefix("/ads").Subrouter()

	authSubrouter.HandleFunc("/signup", handlersAuth.RegisterHandler).Methods(http.MethodPost)
	authSubrouter.HandleFunc("/signin", handlersAuth.LoginHandler).Methods(http.MethodPost)
	authSubrouter.Use(middle.Peflite)

	adSubrouter.HandleFunc("/", handlersAd.Handler).Methods(http.MethodGet)
	adSubrouter.HandleFunc("/create", handlersAd.CreateHandler).Methods(http.MethodPost)
	adSubrouter.HandleFunc("/{ad_id}", handlersAd.UpdateHandler).Methods(http.MethodPut)
	adSubrouter.Use(middle.Peflite, middle.Auth)

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

