package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"2025_2_404/pkg/utils"
	pbAuth "2025_2_404/protos/auth"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	client pbAuth.AuthClient
}

func NewAuthHandler(client pbAuth.AuthClient) *AuthHandler {
	return &AuthHandler{client: client}
}

type registerDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserName string `json:"user_name"`
}

type loginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.UserName == "" {
		http.Error(w, `{"error":"email, password, and user_name are required"}`, http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		http.Error(w, `{"error":"password must be at least 6 characters"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.Register(ctx, &pbAuth.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		UserName: req.UserName,
	})
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"email and password are required"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.Login(ctx, &pbAuth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}