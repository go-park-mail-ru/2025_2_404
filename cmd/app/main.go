package main

import (
	"2025_2_404/internal/config"
	httphandler "2025_2_404/internal/delivery/http"
	db "2025_2_404/internal/connections"
	authrepo "2025_2_404/internal/repository/postgres/auth"
	adrepo "2025_2_404/internal/repository/postgres/ad"
	usecaseauth "2025_2_404/internal/use_case/auth"
	usecasead "2025_2_404/internal/use_case/ad"
	"net/http"
)

func pefliteMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		allowed := map[string]bool{
        "http://localhost:8000": true,
        "http://127.0.0.1:8000": true,
		"http://89.208.230.119:8000": true,
        // добавь нужные домены, если будут
    }
		origin := r.Header.Get("Origin")
		if allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		if r.Method != http.MethodOptions {
			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")

			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization, X-Requested-With, X-CSRF-Token")

			w.Header().Set("Access-Control-Allow-Credentials", "true")

			w.Header().Set("Access-Control-Max-Age", "86400")

			w.WriteHeader(http.StatusOK)
			return
		}
	}
}

func main() {
	config := config.GetConfig()
	postgresql, err := db.ConnectDB(config.DBConfig)
	if err != nil {
		panic(err)
	}
	defer postgresql.Close()

	repoauth := authrepo.New(postgresql)
	repoad := adrepo.New(postgresql)	

	authUsecase := usecaseauth.New(repoauth)
	adUsecase := usecasead.New(repoad)

	handlers := httphandler.New(authUsecase, adUsecase, config.AppConfig.JwtPrivateKey, config.AppConfig.JwtPublicKey)
	http.HandleFunc("/", httphandler.AuthMiddleware(handlers.JwtPublicKey, pefliteMiddleware(handlers.AdHandler)))
	http.HandleFunc("/signup", pefliteMiddleware(handlers.RegisterHandler))
	http.HandleFunc("/signin", pefliteMiddleware(handlers.LoginHandler))

	err = http.ListenAndServe(":"+config.AppConfig.Port, nil)
	if err != nil {
		panic(err)
	}
}

