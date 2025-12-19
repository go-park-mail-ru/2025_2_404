package http

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	"2025_2_404/pkg/utils"
	pbAuth "2025_2_404/protos/gen/go/auth"

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

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(errorResponse{Error: msg}); err != nil {
		slog.Warn("Failed to write HTTP error response", "error", err, "message", msg, "code", code)
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid JSON", http.StatusBadRequest)
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
		httpCode := utils.HTTPStatusFromCode(st.Code())
		log.Printf("gRPC Register error: %v (gRPC code: %v, HTTP code: %d)", err, st.Code(), httpCode)
		writeError(w, st.Message(), httpCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Warn("Failed to encode response in Register handler", "error", err)
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid JSON", http.StatusBadRequest)
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
		httpCode := utils.HTTPStatusFromCode(st.Code())
		log.Printf("gRPC Login error: %v (gRPC code: %v, HTTP code: %d)", err, st.Code(), httpCode)
		writeError(w, st.Message(), httpCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Warn("Failed to encode response in Login handler", "error", err)
	}
}
