package handlers

import (
	modelad "2025_2_404/internal/domain/models/ad"
	modeluser "2025_2_404/internal/domain/models/user"
	"2025_2_404/pkg"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type authUsecaseI interface{
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
	JwtPrivateKey *ecdsa.PrivateKey
	JwtPublicKey  *ecdsa.PublicKey
}

func New(authUsecase authUsecaseI, adUsecase adUsecaseI, jwtPrivateKey *ecdsa.PrivateKey, jwtPublicKey *ecdsa.PublicKey) *FunctionHandler {
	return &FunctionHandler{
		authUsecase: authUsecase,
		adUsecase:   adUsecase,
		JwtPrivateKey: jwtPrivateKey,
		JwtPublicKey:  jwtPublicKey,
	}
}

func JSONResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
		"data":    data,
	})
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

	tokenString, err := pkg.GenerateToken(h.JwtPrivateKey, returnUser.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate token: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, http.StatusOK, "Successful authorization", map[string]string{
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

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	userID, err := h.authUsecase.RegisterUser(ctx, user.Email, user.HashedPassword, user.UserName)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not created: %v", err), http.StatusInternalServerError)
		return
	}

	tokenString, err := pkg.GenerateToken(h.JwtPrivateKey, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate token: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, http.StatusOK, "User created", map[string]string{
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

	

	JSONResponse(w, http.StatusOK, "Successful authorization", map[string]interface{}{
		"ads":	ads,
	})
}	
