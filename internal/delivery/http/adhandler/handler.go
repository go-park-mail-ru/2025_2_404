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
)

type adUsecaseI interface {
	Create(ctx context.Context, ad modelad.Ads) (int, error)
	FindByUserID(ctx context.Context, userID modeluser.ID) (modelad.Ads, error)
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

	adID, err := h.adUsecase.Create(r.Context(), ads)
	if err != nil{
		http.Error(w, fmt.Sprintf("Ad not created: %v", err), http.StatusInternalServerError)
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "Successful create ad", map[string]interface{}{
		"adID":	adID,
		"ads":	ads,
	})
}