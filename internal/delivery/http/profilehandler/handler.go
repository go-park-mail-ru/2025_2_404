package profilehandler

import (
	modeluser "2025_2_404/internal/domain/models/user"
	"2025_2_404/internal/modules"
	"2025_2_404/pkg"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

type profileUsecaseI interface{
	Update(ctx context.Context, client modeluser.User, file io.Reader, ext string) error
	Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, []byte, error)
	Delete(ctx context.Context, clientID modeluser.ID) error
}

type UserHandler struct{
	profileUsecase	profileUsecaseI
}

func New(profileUsecase profileUsecaseI) *UserHandler{
	return &UserHandler{
		profileUsecase: profileUsecase,
	}
}

func (h *UserHandler) UpdateHandler (w http.ResponseWriter, r *http.Request){

	var client modeluser.User
	var err error
	r.Body = http.MaxBytesReader(w, r.Body, 10 * 1024 * 1024)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File is too big.", http.StatusBadRequest)
		return
	}

	client.ID, err = modules.Get(r.Context())
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	client.Email = r.FormValue("email")
	client.UserName = r.FormValue("user_name")
	client.HashedPassword = r.FormValue("password")
	
	imgFile, header, err := r.FormFile("img")
	if err != nil && err != http.ErrMissingFile{
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	if header != nil {
		if imgFile != nil {
			defer imgFile.Close()
		}

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext != ".jpg" && ext != ".png" && ext != ".gif" {
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			return
		}
		
		if err = h.profileUsecase.Update(r.Context(), client, imgFile, ext); err != nil {
			http.Error(w, "client update with image failed:", http.StatusUnprocessableEntity)
			return
		}

	} else { 
		if err = h.profileUsecase.Update(r.Context(), client, nil, ""); err != nil {
			http.Error(w, "client update failed:", http.StatusUnprocessableEntity)
			return
		}
	}
    
	pkg.JSONResponse(w, http.StatusOK, "client update", map[string]interface{}{})
}

func (h *UserHandler) ShowHandler (w http.ResponseWriter, r *http.Request){
	clientID, error := modules.Get(r.Context())
	if error != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	client, bytes, err := h.profileUsecase.Show(r.Context(), clientID)
	if err != nil {
		http.Error(w, "Don't have ads this user", http.StatusInternalServerError)
		return
	}

	avatarBase64 := base64.StdEncoding.EncodeToString(bytes)

	pkg.JSONResponse(w, http.StatusOK, "Successful authorization", map[string]interface{}{
		"client":	client,
		"img":		avatarBase64,
	})
}

func (h *UserHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	clientID, err := modules.Get(r.Context())
	if err != nil {
		http.Error(w, "Bad Request: Unauthorized", http.StatusBadRequest)
		return
	}

	err = h.profileUsecase.Delete(r.Context(), clientID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	pkg.JSONResponse(w, http.StatusNoContent, "user deleted successfully", map[string]interface{}{})
}