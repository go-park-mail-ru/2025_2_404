package supporthandler

import (
	modeluved "2025_2_404/internal/domain/models/support"
	modeluser "2025_2_404/internal/domain/models/user"
	"2025_2_404/internal/modules"
	"2025_2_404/pkg"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type uvedUsecaseI interface {
	Create(ctx context.Context, uved modeluved.Support, file io.Reader, ext string) (error)
	FindByUserID(ctx context.Context, clientID modeluser.ID) ([]modeluved.Support, error)
	GetOneUved(ctx context.Context, uvedID modeluved.ID, userID modeluser.ID) (modeluved.Support, []byte, error)
	Update(ctx context.Context, uved modeluved.Support, file io.Reader, ext string) (error)
	Delete(ctx context.Context, uvedID modeluved.ID, clientID modeluser.ID) (error)
	GetAllSups(ctx context.Context, superID int64) ([]modeluved.Support, error)
}

type Handler struct {
	uvedUsecase	uvedUsecaseI
}

func New(uvedUsecase uvedUsecaseI) *Handler {
	return &Handler{
		uvedUsecase: uvedUsecase,
	}
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {

	clientID, error := modules.Get(r.Context())
	log.Printf("ОШИБКА ПИЗДА %w", clientID)
	if error != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	uved, err := h.uvedUsecase.FindByUserID(r.Context(), clientID)
	if err != nil {
		http.Error(w, "Don't have ads this user", http.StatusInternalServerError)
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "Successful authorization", map[string]interface{}{
		"uved":	uved,
	})
}	

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request){

	r.Body = http.MaxBytesReader(w, r.Body, 10 * 1024 * 1024)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File is too big.", http.StatusBadRequest)
		return
	}


	clientID, error := modules.Get(r.Context())
	if error != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var uved modeluved.Support
	uved.UserID = clientID

	uved.Status = r.FormValue("status")
	uved.Category = r.FormValue("category")
	uved.Description = r.FormValue("description")
	uved.ContactName = r.FormValue("contact_name")
	uved.ContactEmail = r.FormValue("contact_email")

	imgFile, header, err := r.FormFile("image")
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
		
		if err = h.uvedUsecase.Create(r.Context(), uved, imgFile, ext); err != nil {
			http.Error(w, "uved create with image failed", http.StatusUnprocessableEntity)
			return
		}

	} else { 
		if err = h.uvedUsecase.Create(r.Context(), uved, nil, ""); err != nil {
			http.Error(w, "uved create failed:", http.StatusUnprocessableEntity)
			return
		}
	}

	pkg.JSONResponse(w, http.StatusCreated, "Successful create uved", map[string]interface{}{})
}

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request){
	r.Body = http.MaxBytesReader(w, r.Body, 10 * 1024 * 1024)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File is too big.", http.StatusBadRequest)
		return
	}

	var uved modeluved.Support
	vars := mux.Vars(r)
	uvedID, err := strconv.ParseInt(vars["uved_id"], 10, 64)
	if err != nil {
		http.Error(w, "id ad not valid:", http.StatusBadRequest)
		return
	}

	uved.Status = r.FormValue("status")
	uved.Category = r.FormValue("category")
	uved.Description = r.FormValue("description")
	uved.ContactName = r.FormValue("contact_name")
	uved.ContactEmail = r.FormValue("contact_email")
	uved.ID = modeluved.ID(uvedID)

	uved.UserID, err = modules.Get(r.Context())
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	imgFile, header, err := r.FormFile("image")
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
		
		if err = h.uvedUsecase.Update(r.Context(), uved, imgFile, ext); err != nil {
			http.Error(w, "ads update with image failed", http.StatusUnprocessableEntity)
			return
		}

	} else { 
		if err = h.uvedUsecase.Update(r.Context(), uved, nil, ""); err != nil {
			http.Error(w, "ads update failed:", http.StatusUnprocessableEntity)
			return
		}
	}

	pkg.JSONResponse(w, http.StatusOK, "Successful update uved", map[string]interface{}{})
}

func (h * Handler) DeleteHandler(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	uvedID, err := strconv.ParseInt(vars["uved_id"], 10, 64)
	if err != nil {
		http.Error(w, "id ad not valid: ", http.StatusBadRequest)
	}

	clientID, err := modules.Get(r.Context())
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = h.uvedUsecase.Delete(r.Context(), modeluved.ID(uvedID), modeluser.ID(clientID))
	if err != nil {
		http.Error(w, "ad not delete:", http.StatusInternalServerError)
		return
	}

	pkg.JSONResponse(w, http.StatusNoContent, "Successful deleted ad", map[string]interface{}{})	
}

func (h *Handler) GetOneUved(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uvedID, err := strconv.ParseInt(vars["uved_id"], 10, 64)
	fmt.Println(uvedID)
	if err != nil {
		http.Error(w, "id ad not valid:", http.StatusBadRequest)
	}
	clientID, err := modules.Get(r.Context())
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	uved, bytes, err := h.uvedUsecase.GetOneUved(r.Context(), modeluved.ID(uvedID), modeluser.ID(clientID))
	if err != nil {
		http.Error(w, "Don't have ads this user", http.StatusInternalServerError)
		return
	}
	


	pkg.JSONResponse(w, http.StatusOK, "Successful search", map[string]interface{}{
		"uved": uved,
		"image": bytes,
	})
}

func (h *Handler) HandlerGetAll(w http.ResponseWriter, r *http.Request) {

	clientID, error := modules.Get(r.Context())
	if error != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	uved, err := h.uvedUsecase.GetAllSups(r.Context(), int64(clientID))
	if err != nil {
		http.Error(w, "Don't have ads this user", http.StatusInternalServerError)
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "Successful authorization", map[string]interface{}{
		"uved":	uved,
	})
}	