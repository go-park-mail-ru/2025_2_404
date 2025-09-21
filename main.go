package main

import (
	"encoding/json"
	"net/http"
)

// loginHandler handles user login requests. It processes incoming HTTP requests,
// validates user credentials, and manages user authentication flow.
// The function writes the appropriate HTTP response based on the authentication result.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Не тот метод", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Json not correct", http.StatusBadRequest)
		return
	}

	// Пример проверки логина и пароля (заглушка)
	if creds.UserName != "admin" || creds.Password != "password1234" {
		http.Error(w, "Неверные имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Успешная авторизация",
		"token":   "Auth token",
	})
}

// registerHandler handles user registration requests.
// It expects a POST request with a JSON body containing "user_name", "email", and "password" fields.
// The handler validates the request method, parses and validates the input data,
// and responds with a JSON message and a placeholder token upon successful registration.
// Returns appropriate HTTP error codes for invalid methods, malformed JSON, or validation failures.
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Не тот метод", http.StatusMethodNotAllowed)
		return
	}

	var user struct {
		UserName string `json:"user_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Json not correct", http.StatusBadRequest)
		return
	}

	if len(user.Password) < 10 || len(user.UserName) < 4 {
		http.Error(w, "Name or password no validate", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Пользователь создался",
		"token":   "Create token",
	})
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Не тот метод", http.StatusMethodNotAllowed)
		return
	}

	var Add struct{
		AddId		string	`json:"add_id"`
		CreaterID	string	`json:"creater_id"`
		FilePath	string	`json:"file_path"`
		Title		string	`json:"title"`
		Text		string	`json:"text"`
	}
	
	var Creater struct{
		CreaterID	string	`json:"creater_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&Creater); err != nil{
		http.Error(w,"Json not correct", http.StatusBadRequest)
		return
	}
}

func main() {
	http.HandleFunc("/", handle)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
