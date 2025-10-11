package handlers

import (
<<<<<<< HEAD
	"2025_2_404/internal/models"
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
=======
	modeluser "2025_2_404/internal/models/user"
	modelad "2025_2_404/internal/models/ad"
>>>>>>> a8230ea6cc45a4ef7d6d317222973fdc7959bd18
	"encoding/json"
	"net/http"
	"context"
)

type authUsecaseI interface{
	SessionGenerateAndSave(ctx context.Context, userID modeluser.ID) (string, error)
	RegisterUser(ctx context.Context, email, password, userName string) (modeluser.ID, error)
}

type adUsecaseI interface{
	CreateAd(ctx context.Context, ad modelad.Ads) (int, error)
	FindAdByUserID(ctx context.Context, userID modeluser.ID) (modelad.Ads, error)
}

type userUsecaseI interface{
	CheckUser(ctx context.Context, email string, password string) (modeluser.User, error)
	FindSessionByUserID(ctx context.Context, userID modeluser.ID) (string, error)
	FindUser(ctx context.Context, sessionID string) (modeluser.ID, error)
}

type FunctionHandler struct {
	authUsecase authUsecaseI
	adUsecase   adUsecaseI
	userUsecase userUsecaseI
}

func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}

func (h *FunctionHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

<<<<<<< HEAD
	var creds models.User
=======
	var creds modeluser.User
>>>>>>> a8230ea6cc45a4ef7d6d317222973fdc7959bd18
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Json not correct", http.StatusBadRequest)
		return
	}

<<<<<<< HEAD
	if validateErr := pkg.ValidateLoginUser(&creds); validateErr != nil {
		http.Error(w, validateErr.Error(), http.StatusUnprocessableEntity)
		return
	}

	passwordHash := sha1.Sum([]byte(creds.GetHashedPassword()))
	hexPasswordHash := hex.EncodeToString(passwordHash[:])

	returnUserID, err := h.foundUserByCredentialsDB(creds.GetEmail(), hexPasswordHash)
=======
	returnUser, err := h.userUsecase.CheckUser(r.Context(), creds.Email, creds.HashedPassword)
>>>>>>> a8230ea6cc45a4ef7d6d317222973fdc7959bd18
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	var sessionID string
	sessionID, err = h.userUsecase.FindSessionByUserID(r.Context(), returnUser.ID)
	if err == nil {
		sessionID, err = h.authUsecase.SessionGenerateAndSave(r.Context(), returnUser.ID)
		if err != nil {
			http.Error(w, "Session not generated", http.StatusInternalServerError)
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


func (h *FunctionHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

<<<<<<< HEAD
	var user models.User
=======
	var user modeluser.User
>>>>>>> a8230ea6cc45a4ef7d6d317222973fdc7959bd18

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Json not correct", http.StatusBadRequest)
		return
	}

	userID, err := h.authUsecase.RegisterUser(r.Context(), user.Email, user.HashedPassword, user.UserName)
	if err != nil {
		http.Error(w, "User not created", http.StatusInternalServerError)
		return
	}

<<<<<<< HEAD
	passwordHash := sha1.Sum([]byte(user.GetHashedPassword()))
	hexPasswordHash := hex.EncodeToString(passwordHash[:])

	var returnUserID int
	if err := h.DB.QueryRow(sqlTextForInsertUsers, user.GetEmail(), hexPasswordHash, user.GetUserName()).Scan(&returnUserID); err != nil{
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
=======
	sessionID, err := h.authUsecase.SessionGenerateAndSave(r.Context(), userID)
	if err != nil {
		http.Error(w, "Session not created", http.StatusInternalServerError)
>>>>>>> a8230ea6cc45a4ef7d6d317222973fdc7959bd18
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


func (h *FunctionHandler) Handle(w http.ResponseWriter, r *http.Request) {
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
	userID, err := h.userUsecase.FindUser(r.Context(), sessionID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

<<<<<<< HEAD
	var ads models.Ads
	if err := h.DB.QueryRow(sqlTextForSelectAds, userID).Scan(&ads.ID, &ads.FilePath, &ads.Title, &ads.Text); err != nil {
		http.Error(w, "Failed to retrieve ads", http.StatusInternalServerError)
		return
	}
	ads.SetCreatorID(userID)
	adsOut := []map[string]string{
		{
			"add_id":     string(ads.GetID()),
			"creater_id": string(ads.GetCreatorID()),
			"file_path":  ads.GetFilePath(),
			"title":      ads.GetTitle(),
			"text":       ads.GetText(),
		},
	}
=======
>>>>>>> a8230ea6cc45a4ef7d6d317222973fdc7959bd18
	JSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "Successful authorization",
		"ads":     userID,
	})
}	
