package adhandler

import (
	modelad "2025_2_404/internal/domain/models/ad"
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
	"net/http"
	"2025_2_404/pkg"
	"2025_2_404/internal/modules"
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
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	userID, error := modules.Get(r.Context())
	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ads, err := h.adUsecase.FindByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to retrieve ads", http.StatusInternalServerError)
		return
	}

	

	pkg.JSONResponse(w, http.StatusOK, "Successful authorization", map[string]interface{}{
		"ads":	ads,
	})
}	