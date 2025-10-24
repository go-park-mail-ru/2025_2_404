package main

import (
	"2025_2_404/internal/config"
	db "2025_2_404/internal/connections"
	adhandler "2025_2_404/internal/delivery/http/adhandler"
	authhandler "2025_2_404/internal/delivery/http/authhandler"
	middleware "2025_2_404/internal/delivery/http/middleware"
	adrepo "2025_2_404/internal/repository/postgres/ad"
	authrepo "2025_2_404/internal/repository/postgres/auth"
	usecasead "2025_2_404/internal/use_case/ad"
	usecaseauth "2025_2_404/internal/use_case/auth"
	tokenusecase "2025_2_404/internal/use_case/token"
	"net/http"
)

func main() {
	config := config.GetConfig()
	postgresql, err := db.ConnectDB(config.DBConfig)
	if err != nil {
		panic(err)
	}
	defer postgresql.Close()

	repoAuth := authrepo.New(postgresql)
	repoAd := adrepo.New(postgresql)	

	tokenUsecae := tokenusecase.New(&config)
	authUsecase := usecaseauth.New(repoAuth, tokenUsecae)
	adUsecase := usecasead.New(repoAd)

	middle := middleware.New(tokenUsecae)
	handlersAd := adhandler.New(adUsecase)
	handlersAuth := authhandler.New(authUsecase)
	http.HandleFunc("/", middle.Auth(middle.Peflite(handlersAd.Handler)))
	http.HandleFunc("/signup", middle.Peflite(handlersAuth.RegisterHandler))
	http.HandleFunc("/signin", middle.Peflite(handlersAuth.LoginHandler))

	err = http.ListenAndServe(":"+config.AppConfig.Port, nil)
	if err != nil {
		panic(err)
	}
}

