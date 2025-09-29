package handlers

import (
	"2025_2_404/models"
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
)


type Handlers struct {
	DB *sql.DB
}


func New(db *sql.DB) *Handlers {
	return &Handlers{DB: db}
}


func GenerateSession() (string, error){
	sessionID := make([]byte,32)
	if _, err := rand.Read(sessionID); err != nil{
		return "", err
	}
	return hex.EncodeToString(sessionID), nil
}

func (h *Handlers) foundUserBySessionDB(sessionID string) (string, error) {
	var userID string
	sqlTextForFoundUser := "SELECT user_id FROM session WHERE session_id = $1"
	err := h.DB.QueryRow(sqlTextForFoundUser, sessionID).Scan(&userID)
	if err != nil {
		return "", errors.New("session not found")
	}
	return userID, nil
}

func (h *Handlers) foundUserByCredentialsDB(email, password string) (string, error) {
	var userID string
	var hashedPassword string
	sqlTextForSelectUsers := "SELECT id, password FROM users WHERE email = $1 "
	err := h.DB.QueryRow(sqlTextForSelectUsers, email).Scan(&userID, &hashedPassword)
	if err != nil {
		return "", err
	}
	if hashedPassword != password {
		return "", errors.New("invalid password")
	}
	return userID, nil
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

	passwordHash := sha1.Sum([]byte(creds.Password))
	hexPasswordHash := hex.EncodeToString(passwordHash[:])

	returnUserID, err := h.foundUserByCredentialsDB(creds.Email, hexPasswordHash)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	var sessionID string
	sqlTextForCheckSession := "SELECT session_id FROM session WHERE user_id = $1"
	err = h.DB.QueryRow(sqlTextForCheckSession, returnUserID).Scan(&sessionID)
	if err != nil {
		sessionID, err = GenerateSession()
		if err != nil {
			http.Error(w, "Session not generated", http.StatusInternalServerError)
			return
		}

		sqlTextForInsertSession := "INSERT INTO session (user_id, session_id) VALUES ($1, $2)"
		_, err = h.DB.Exec(sqlTextForInsertSession, returnUserID, sessionID)
		if err != nil {
			http.Error(w, "Session token not created", http.StatusConflict)
			return
		}
	} 
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

	passwordHash := sha1.Sum([]byte(user.Password))
	hexPasswordHash := hex.EncodeToString(passwordHash[:])

	var returnUserID int
	sqlTextForInsertUsers := "INSERT INTO users (email, password, user_name) VALUES ( $1, $2, $3) RETURNING id"
	if err := h.DB.QueryRow(sqlTextForInsertUsers, user.Email, hexPasswordHash, user.UserName).Scan(&returnUserID); err != nil{
		http.Error(w, "User already registered", http.StatusUnprocessableEntity)
		return
	}

	sessionID, err := GenerateSession()
	if err != nil{
		http.Error(w, "Session not generated", http.StatusInternalServerError)
		return
	}

	sqlTextForInsertSession := "INSERT INTO session (user_id, session_id) VALUES ($1, $2)"
	_, err = h.DB.Exec(sqlTextForInsertSession, returnUserID, sessionID)
	if err != nil{
		http.Error(w,"Session token not created", http.StatusConflict)
		return
	}

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
	userID, err := h.foundUserBySessionDB(sessionID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	filePath := "images/1.jpg"
	title := "Новый набор студентов в технопарк!"
	text := "Этой весной наступает новый набор студентов в Технопарк по программе WEB-разработка. Обучим вас GO-lang! Студенты нашего курса разработали сайт ADS"
	sqlTextForInsertAds := "INSERT INTO ads (creator_id, file_path, title, text) VALUES ($1, $2, $3, $4)"
	_, err = h.DB.Exec(sqlTextForInsertAds, userID, filePath, title, text)
	if err != nil {
		http.Error(w, "Failed to insert ad", http.StatusInternalServerError)
		return
	}

	var ad models.Ads
	sqlTextForSelectAds := "SELECT id, file_path, title, text FROM ads WHERE creator_id = $1"
	err = h.DB.QueryRow(sqlTextForSelectAds, userID).Scan(&ad.AdID, &ad.FilePath, &ad.Title, &ad.Text)
	if err != nil {
		http.Error(w, "Failed to retrieve ads", http.StatusInternalServerError)
		return
	}

	ads := []map[string]string{
		{
			"add_id":     ad.AdID,
			"creater_id": userID,
			"file_path":  ad.FilePath,
			"title":      ad.Title,
			"text":       ad.Text,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ads": ads,
	})
}	
