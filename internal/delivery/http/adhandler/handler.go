package adhandler

import (
	modelad "2025_2_404/internal/domain/models/ad"
	modeluser "2025_2_404/internal/domain/models/user"
	"2025_2_404/internal/modules"
	"2025_2_404/pkg"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type adUsecaseI interface {
	Create(ctx context.Context, ad modelad.Ads) (error)
	FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelad.Ads, error)
	Update(ctx context.Context, ad modelad.Ads) (error)
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

	userID, error := modules.Get(r.Context())
	if error != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var ads modelad.Ads
	ads.ClientID = userID
	if err := json.NewDecoder(r.Body).Decode(&ads); err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusBadRequest)
		return
	}

	err := h.adUsecase.Create(r.Context(), ads)
	if err != nil{
		http.Error(w, fmt.Sprintf("Ad not created: %v", err), http.StatusInternalServerError)
		return
	}

	pkg.JSONResponse(w, http.StatusCreated, "Successful create ad", map[string]interface{}{})
}

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request){

	var ads modelad.Ads
	vars := mux.Vars(r)
	adID, err := strconv.ParseInt(vars["ad_id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("id ad not valid: %v", err), http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&ads); err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusBadRequest)
		return
	}

	ads.ID = modelad.ID(adID)

	ads.ClientID, err = modules.Get(r.Context())
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = h.adUsecase.Update(r.Context(), ads)
	if err != nil {
		http.Error(w, fmt.Sprintf("ad not update: %v", err), http.StatusInternalServerError)
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "Successful update ad", map[string]interface{}{})
}