package handlers

import (
	modeluser "2025_2_404/internal/domain/models/user"
	modelad "2025_2_404/internal/domain/models/ad"
	"encoding/json"
	"net/http"
	"context"
	"2025_2_404/pkg"
)

type authUsecaseI interface{
	//SessionGenerateAndSave(ctx context.Context, userID modeluser.ID) (string, error)
	RegisterUser(ctx context.Context, email, password, userName string) (modeluser.ID, error)
	CheckUser(ctx context.Context, email string, password string) (modeluser.User, error)
}

type adUsecaseI interface{
	CreateAd(ctx context.Context, ad modelad.Ads) (int, error)
	FindAdByUserID(ctx context.Context, userID modeluser.ID) (modelad.Ads, error)
}

type FunctionHandler struct {
	authUsecase authUsecaseI
	adUsecase   adUsecaseI
}

func New(authUsecase authUsecaseI, adUsecase adUsecaseI) *FunctionHandler {
	return &FunctionHandler{
		authUsecase: authUsecase,
		adUsecase:   adUsecase,
	}
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

	returnUser, err := h.authUsecase.CheckUser(r.Context(), creds.Email, creds.HashedPassword)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// var sessionID string
	// sessionID, err = h.userUsecase.FindSessionByUserID(r.Context(), returnUser.ID)
	// if err == nil {
	// 	sessionID, err = h.authUsecase.SessionGenerateAndSave(r.Context(), returnUser.ID)
	// 	if err != nil {
	// 		http.Error(w, "Session not generated", http.StatusInternalServerError)
	// 		return
	// 	}
	// } 

	tokenString, err := pkg.GenerateToken(int64(returnUser.ID))
	if err != nil {
		http.Error(w, "Не получилось сгенерить токен", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		// Value:    sessionID,
		Path:     "/",
		MaxAge:   8080,
		HttpOnly: true,
	})

	JSONResponse(w, http.StatusCreated, map[string]string {  
		"message": "Successful authorization",
	})

	JSONResponse(w, http.StatusOK, map[string]string {
		"token": tokenString,
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

	// sessionID, err := h.authUsecase.SessionGenerateAndSave(r.Context(), userID)
	// if err != nil {
	// 	http.Error(w, "Session not created", http.StatusInternalServerError)
	// 	return
	// }

	tokenString, err := pkg.GenerateToken(int64(userID))
	if err != nil {
		http.Error(w, "Не удалось создать токен", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		// Value:    sessionID,
		Path:     "/",
		MaxAge:   8080,
		HttpOnly: true,
	})

	JSONResponse(w, http.StatusCreated, map[string]string {
		"message": "User created",
	})

	JSONResponse(w, http.StatusOK, map[string]string {
		"token": tokenString,
	})
}


func (h *FunctionHandler) AdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(UserIDKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ads, err := h.adUsecase.FindAdByUserID(r.Context(), modeluser.ID(userID))
	if err != nil {
		http.Error(w, "Failed to retrieve ads", http.StatusInternalServerError)
		return
	}

	

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Successful authorization",
		"ads":     ads,
	})
}	
