package adhandler

import (
	modelad "2025_2_404/internal/domain/models/ad"
	modelfullad "2025_2_404/internal/domain/models/ad_full_info"
	modeluser "2025_2_404/internal/domain/models/user"
	"2025_2_404/internal/modules"
	"2025_2_404/pkg"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type adUsecaseI interface {
	Create(ctx context.Context, ad modelad.Ads, file io.Reader, ext string) (error)
	FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelad.Ads, error)
	GetOneAd(ctx context.Context, adID int64) (modelfullad.AdFullInfo, int, []byte, error)
	Update(ctx context.Context, ad modelad.Ads, file io.Reader, ext string) (error)
	Delete(ctx context.Context, adID int64) (error)
}

type Handler struct {
	adUsecase adUsecaseI
}

func New(adUsecase adUsecaseI) *Handler {
	return &Handler{
		adUsecase: adUsecase,
	}
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {

	userID, error := modules.Get(r.Context())
	if error != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ads, err := h.adUsecase.FindByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Don't have ads this user", http.StatusInternalServerError)
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "Successful authorization", map[string]interface{}{
		"ads":	ads,
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

	var ads modelad.Ads
	ads.ClientID = clientID

	ads.Title = r.FormValue("title")
	ads.Content = r.FormValue("content")
	ads.TargetUrl = r.FormValue("target_url")

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
		
		if err = h.adUsecase.Create(r.Context(), ads, imgFile, ext); err != nil {
			http.Error(w, "ads create with image failed", http.StatusUnprocessableEntity)
			return
		}

	} else { 
		if err = h.adUsecase.Create(r.Context(), ads, nil, ""); err != nil {
			http.Error(w, "ads create failed:", http.StatusUnprocessableEntity)
			return
		}
	}

	pkg.JSONResponse(w, http.StatusCreated, "Successful create ad", map[string]interface{}{})
}

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request){
	r.Body = http.MaxBytesReader(w, r.Body, 10 * 1024 * 1024)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File is too big.", http.StatusBadRequest)
		return
	}

	var ads modelad.Ads
	vars := mux.Vars(r)
	adID, err := strconv.ParseInt(vars["ad_id"], 10, 64)
	if err != nil {
		http.Error(w, "id ad not valid:", http.StatusBadRequest)
		return
	}

	ads.Title = r.FormValue("title")
	ads.Content = r.FormValue("content")
	ads.TargetUrl = r.FormValue("target_url")
	ads.ID = modelad.ID(adID)

	ads.ClientID, err = modules.Get(r.Context())
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
		
		if err = h.adUsecase.Update(r.Context(), ads, imgFile, ext); err != nil {
			http.Error(w, "ads update with image failed", http.StatusUnprocessableEntity)
			return
		}

	} else { 
		if err = h.adUsecase.Update(r.Context(), ads, nil, ""); err != nil {
			http.Error(w, "ads update failed:", http.StatusUnprocessableEntity)
			return
		}
	}

	pkg.JSONResponse(w, http.StatusOK, "Successful update ad", map[string]interface{}{})
}

func (h * Handler) DeleteHandler(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	adID, err := strconv.ParseInt(vars["ad_id"], 10, 64)
	if err != nil {
		http.Error(w, "id ad not valid: ", http.StatusBadRequest)
	}

	err = h.adUsecase.Delete(r.Context(), adID)
	if err != nil {
		http.Error(w, "ad not delete:", http.StatusInternalServerError)
		return
	}

	pkg.JSONResponse(w, http.StatusNoContent, "Successful deleted ad", map[string]interface{}{})	
}

func (h *Handler) GetOneAd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	adID, err := strconv.ParseInt(vars["ad_id"], 10, 64)
	fmt.Println(adID)
	if err != nil {
		http.Error(w, "id ad not valid:", http.StatusBadRequest)
	}
	ad, conversion, bytes, err := h.adUsecase.GetOneAd(r.Context(), adID)
	if err != nil {
		http.Error(w, "Don't have ads this user", http.StatusInternalServerError)
		return
	}
	


	pkg.JSONResponse(w, http.StatusOK, "Successful search", map[string]interface{}{
		"ad": ad,
		"conversion": conversion,
		"image": bytes,
	})
}