package authhandler

import (
	modeluser "2025_2_404/internal/domain/models/user"
	"2025_2_404/pkg"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type authUsecaseI interface {
	Register(ctx context.Context, email, password, userName string) (string, error)
	Login(ctx context.Context, email string, password string) (string, error)
}

type AuthHandler struct {
	authUsecase authUsecaseI
}

func New(authUsecase authUsecaseI) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	
	}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var creds modeluser.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Json not correct", http.StatusBadRequest)
		return
	}

	token, err := h.authUsecase.Login(r.Context(), creds.Email, creds.HashedPassword)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "Successful authorization", map[string]string{
		"token": token,
	})
}


func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var user modeluser.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Json not correct", http.StatusBadRequest)
		return
	}

	token, err := h.authUsecase.Register(r.Context(), user.Email, user.HashedPassword, user.UserName)
	if err != nil {
		log.Printf("!!! REGISTRATION ERROR: %v", err)
		http.Error(w, "User not created:", http.StatusConflict)
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "User created", map[string]string{
		"token": token,
	})
}


