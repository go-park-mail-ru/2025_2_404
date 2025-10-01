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
	"2025_2_404/pkg"
)


type Handlers struct {
	DB *sql.DB
}


func New(db *sql.DB) *Handlers {
	return &Handlers{DB: db}
}

const(
	sqlTextForInsertSession = "INSERT INTO session (user_id, session_id) VALUES ($1, $2)"
	sqlTextForFoundUser = "SELECT user_id FROM session WHERE session_id = $1"
	sqlTextForSelectUsers = "SELECT id, password FROM app_user WHERE email = $1 "
	sqlTextForCheckSession = "SELECT session_id FROM session WHERE user_id = $1"
	sqlTextForSelectAds = "SELECT id, file_path, title, text_ad FROM ad WHERE creator_id = $1"
	sqlTextForInsertUsers = "INSERT INTO app_user (email, password, user_name) VALUES ( $1, $2, $3) RETURNING id"
	sqlTextForInsertAds = "INSERT INTO ad (creator_id, file_path, title, text_ad) VALUES ($1, $2, $3, $4)"
)

func GenerateSession() (string, error){
	sessionID := make([]byte,32)
	if _, err := rand.Read(sessionID); err != nil{
		return "", err
	}
	return hex.EncodeToString(sessionID), nil
}

func (h *Handlers) foundUserBySessionDB(sessionID string) (string, error) {
	var userID string
	err := h.DB.QueryRow(sqlTextForFoundUser, sessionID).Scan(&userID)
	if err != nil {
		return "", errors.New("session not found")
	}
	return userID, nil
}

func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}

func (h *Handlers) foundUserByCredentialsDB(email, password string) (string, error) {
	var userID string
	var hashedPassword string
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

	if validateErr := pkg.ValidateLoginUser(&creds); validateErr != nil {
		http.Error(w, validateErr.Error(), http.StatusUnprocessableEntity)
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
	err = h.DB.QueryRow(sqlTextForCheckSession, returnUserID).Scan(&sessionID)
	if err != nil {
		sessionID, err = GenerateSession()
		if err != nil {
			http.Error(w, "Session not generated", http.StatusInternalServerError)
			return
		}

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

	JSONResponse(w, http.StatusCreated, map[string]interface{}{
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

	if validateErr := pkg.ValidateRegisterUser(&user); validateErr != nil {
		http.Error(w, validateErr.Error(), http.StatusUnprocessableEntity)
		return
	}

	passwordHash := sha1.Sum([]byte(user.Password))
	hexPasswordHash := hex.EncodeToString(passwordHash[:])

	var returnUserID int
	if err := h.DB.QueryRow(sqlTextForInsertUsers, user.Email, hexPasswordHash, user.UserName).Scan(&returnUserID); err != nil{
		http.Error(w, "User already registered", http.StatusUnprocessableEntity)
		return
	}

	sessionID, err := GenerateSession()
	if err != nil{
		http.Error(w, "Session not generated", http.StatusInternalServerError)
		return
	}

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

	JSONResponse(w, http.StatusCreated, map[string]interface{}{
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

	var ads models.Ads
	if err := h.DB.QueryRow(sqlTextForSelectAds, userID).Scan(&ads.ID, &ads.FilePath, &ads.Title, &ads.Text); err != nil {
		http.Error(w, "Failed to retrieve ads", http.StatusInternalServerError)
		return
	}
	ads.CreatorID = userID
	adsOut := []map[string]string{
		{
			"add_id":     ads.ID,
			"creater_id": ads.CreatorID,
			"file_path":  ads.FilePath,
			"title":      ads.Title,
			"text":       ads.Text,
		},
	}
	JSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "Successful authorization",
		"ads":     adsOut,
	})
}	
