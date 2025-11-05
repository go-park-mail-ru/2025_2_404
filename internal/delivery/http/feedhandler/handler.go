package feedhandler

import (
	modelad "2025_2_404/internal/domain/models/ad"
	"2025_2_404/pkg"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type feedUsecaseI interface{
	Give(ctx context.Context, platformName string) ([]modelad.Ads, error)
}

type feedHandler struct{
	feedUsecase feedUsecaseI
}

func New(feedUsecase feedUsecaseI) *feedHandler{
	return &feedHandler{
		feedUsecase: feedUsecase,
	}
}

func (h *feedHandler) GetAdFeedHandler(w http.ResponseWriter, r *http.Request) {
    platformName := mux.Vars(r)["platform_name"]

    banners, err := h.feedUsecase.Give(r.Context(), platformName)
    if err != nil {
        log.Printf("Ошибка при получении баннеров для feed_id=%s: %v", platformName, err)
        banners = []modelad.Ads{} 
    }

	pkg.JSONResponse(w, http.StatusOK, "feed ads", map[string]interface{}{
		"platform_name": platformName,
		"banners": banners,
	})
}