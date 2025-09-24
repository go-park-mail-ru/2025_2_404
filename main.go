package main

import (
	"net/http"
	"2025_2_404/handlers"
	"2025_2_404/db"
)

func main() {
	
	postgresql, err := db.ConnectDB()
	if err != nil{
		panic(err)
	}
	defer postgresql.Close()
	var handlers = handlers.New()
	http.HandleFunc("/", handlers.Handle)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	erri := http.ListenAndServe(":8080", nil)
	if erri != nil {
		panic(erri)
	}
}
