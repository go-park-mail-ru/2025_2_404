package handlers

import (
	modeluser "2025_2_404/internal/domain/models/user"
	modelad "2025_2_404/internal/domain/models/ad"
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

	var creds modeluser.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Json not correct", http.StatusBadRequest)
		return
	}

	returnUser, err := h.userUsecase.CheckUser(r.Context(), creds.Email, creds.HashedPassword)
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

	var user modeluser.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Json not correct", http.StatusBadRequest)
		return
	}

	userID, err := h.authUsecase.RegisterUser(r.Context(), user.Email, user.HashedPassword, user.UserName)
	if err != nil {
		http.Error(w, "User not created", http.StatusInternalServerError)
		return
	}

	sessionID, err := h.authUsecase.SessionGenerateAndSave(r.Context(), userID)
	if err != nil {
		http.Error(w, "Session not created", http.StatusInternalServerError)
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

	JSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "Successful authorization",
		"ads":     userID,
	})
}	
