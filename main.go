package main

import (
	"2025_2_404/db"
	"2025_2_404/handlers"
	"net/http"
	"2025_2_404/config"
)

func main() {

	postgresql, err := db.ConnectDB(config.GetPostgresConfig())
	if err != nil {
		panic(err)
	}
	defer postgresql.Close()
	var handlers = handlers.New(postgresql)
	http.HandleFunc("/", handlers.Handle)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	err = http.ListenAndServe(":"+config.GetAppPort(), nil)
	if err != nil {
		panic(err)
	}
}
