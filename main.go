package main

import (
	"encoding/json"
	"net/http"
	"2025_2_404/models"
)

// loginHandler handles user login requests. It processes incoming HTTP requests,
// validates user credentials, and manages user authentication flow.
// The function writes the appropriate HTTP response based on the authentication result.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	var creds models.BaseUser
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Json not correct", http.StatusBadRequest)
		return
	}

	// Пример проверки логина и пароля (заглушка)
	if creds.UserName != "admin" || creds.Password != "password1234" {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
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
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	var user models.RegisterUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Json not correct", http.StatusBadRequest)
		return
	}

	if len(user.Password) < 10 || len(user.UserName) < 4 {
		http.Error(w, "Name or password no validate", http.StatusUnprocessableEntity)
		return
	}

	// session.token = "Create token"
	// session.user_id = "Create user_id"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created",
		"token":   "Create token",
	})
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	var req models.Session
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Json not correct", http.StatusBadRequest)
		return
	}

	req.UserID = "Extracted UserID"

	ads := []map[string]string{
		{
			"add_id":     "1",
			"creater_id": req.UserID,
			"file_path":  "/files/ad1.jpg",
			"title":      "Реклама 1",
			"text":       "Текст рекламы 1",
		},
		{
			"add_id":     "2",
			"creater_id": req.UserID,
			"file_path":  "/files/ad2.jpg",
			"title":      "Реклама 2",
			"text":       "Текст рекламы 2",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ads": ads,
	})
}	

func main() {
	http.HandleFunc("/", handle)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
