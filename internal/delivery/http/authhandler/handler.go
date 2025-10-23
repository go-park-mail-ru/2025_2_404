package authhandler

import (
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"2025_2_404/pkg"
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
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

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

	token, err := h.authUsecase.Register(ctx, user.Email, user.HashedPassword, user.UserName)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not created: %v", err), http.StatusInternalServerError)
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "User created", map[string]string{
		"token": token,
	})
}
