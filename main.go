package main

import (
	"2025_2_404/db"
	"2025_2_404/handlers"
	"net/http"
	"2025_2_404/config"
)

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		allowed := map[string]bool{
        "http://localhost:8000": true,
        "http://127.0.0.1:8000": true,
        // добавь нужные домены, если будут
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

	handlers := handlers.New(postgresql)
	http.HandleFunc("/", CORS(handlers.Handle))
	http.HandleFunc("/signup", CORS(handlers.RegisterHandler))
	http.HandleFunc("/signin", CORS(handlers.LoginHandler))

	err = http.ListenAndServe(":"+config.AppConfig.Port, nil)
	if err != nil {
		panic(err)
	}
}
