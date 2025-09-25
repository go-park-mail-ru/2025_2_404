package handlers

import(
	"encoding/json"
	"net/http"
	"2025_2_404/models"
	"database/sql"
)

type Handlers struct {
	DB *sql.DB
}

func New(db *sql.DB) *Handlers {
	return &Handlers{DB: db}
}

func foundUserBySessionDB(sessionID string) string {
	if sessionID == "valid_session_id" {
		return "user_id"
	}
	return ""
}

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	sessionID := "Create session ID"

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   8080,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successful authorization",
	})
}

func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	sessionID := "Create session ID"

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   8080,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created",
	})
}

func (h *Handlers) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID := sessionCookie.Value
	userID := foundUserBySessionDB(sessionID)

	ads := []map[string]string{
		{
			"add_id":     "1",
			"creater_id": userID,
			"file_path":  "/files/ad1.jpg",
			"title":      "Реклама 1",
			"text":       "Текст рекламы 1",
		},
		{
			"add_id":     "2",
			"creater_id": userID,
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
