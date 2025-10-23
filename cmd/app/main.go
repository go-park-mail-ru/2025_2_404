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

func pefliteMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		allowed := map[string]bool{
        "http://localhost:8000": true,
        "http://127.0.0.1:8000": true,
		"http://89.208.230.119:8000": true,
    }
		origin := r.Header.Get("Origin")
		if allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-CSRF-Token")

		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

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
	http.HandleFunc("/", middle.AuthMiddleware(pefliteMiddleware(handlersAd.Handler)))
	http.HandleFunc("/signup", pefliteMiddleware(handlersAuth.RegisterHandler))
	http.HandleFunc("/signin", pefliteMiddleware(handlersAuth.LoginHandler))

	err = http.ListenAndServe(":"+config.AppConfig.Port, nil)
	if err != nil {
		panic(err)
	}
}

