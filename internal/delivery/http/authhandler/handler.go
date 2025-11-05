package authhandler

import (
	"2025_2_404/internal/config"
	modeluser "2025_2_404/internal/domain/models/user"
	"2025_2_404/internal/modules"
	"2025_2_404/pkg"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type authUsecaseI interface {
	Register(ctx context.Context, email, password, userName string) (string, error)
	Login(ctx context.Context, email string, password string) (string, error)
	AddImage(ctx context.Context, userID int64, imageUrl string) error
}

type AuthHandler struct {
	authUsecase authUsecaseI
	config      *config.Config
}

func New(authUsecase authUsecaseI) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		config:      config.GetConfig(),
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
		http.Error(w, fmt.Sprintf("User not created: %v", err), http.StatusInternalServerError)
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "User created", map[string]string{
		"token": token,
	})
}

func (h *AuthHandler) AddImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File is too big.", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid file.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	uniqueFilename := uuid.New().String() + ext

	if ext != ".jpg" && ext != ".png" && ext != ".gif" {
		http.Error(w, "Invalid type of file", http.StatusBadRequest)
		return
	}

	uploadPath := h.config.AppConfig.StoragePath
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		http.Error(w, "Failed to create upload directory.", http.StatusInternalServerError)
		return
	}

	fullPath := filepath.Join(uploadPath, uniqueFilename)
	dst, err := os.Create(fullPath)
	if err != nil {
		http.Error(w, "Failed to create file on server.", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file.", http.StatusInternalServerError)
		return
	}

	userID, err := modules.Get(r.Context())
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = h.authUsecase.AddImage(r.Context(), int64(userID), fullPath)
	if err != nil {
		http.Error(w, "Failed to update user profile.", http.StatusInternalServerError)
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "Avatar uploaded successfully", map[string]string{})
}
