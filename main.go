package main

import (
	"2025_2_404/db"
	"2025_2_404/handlers"
	"net/http"
	"2025_2_404/config"
)

func main() {
	config := config.GetConfig()
	postgresql, err := db.ConnectDB(config.DBConfig)
	if err != nil {
		panic(err)
	}
	defer postgresql.Close()

	handlers := handlers.New(postgresql)
	http.HandleFunc("/", handlers.Handle)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	err = http.ListenAndServe(":"+config.AppConfig.Port, nil)
	if err != nil {
		panic(err)
	}
}
