package balancehandler

import (
	modelwallet "2025_2_404/internal/domain/models/client_wallet"
	modeluser "2025_2_404/internal/domain/models/user"
	"2025_2_404/internal/modules"
	"2025_2_404/pkg"
	"context"
	"fmt"
	"net/http"
)

type balanceUsecaseI interface{
	Show(ctx context.Context, clientID modeluser.ID) (modelwallet.Balance, error)
}

type balanceHandler struct{
	balanceUsecase balanceUsecaseI
}

func New(balanceUsecase balanceUsecaseI) *balanceHandler{
	return &balanceHandler{
		balanceUsecase: balanceUsecase,
	}
}

func (h *balanceHandler) Show(w http.ResponseWriter, r *http.Request){
	
	clientID, err := modules.Get(r.Context())
	if err != nil{
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	balance, err := h.balanceUsecase.Show(r.Context(), clientID)
	if err != nil {
		http.Error(w, fmt.Sprintf("internal server error: %s", err), http.StatusInternalServerError)
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "Successful authorization", map[string]interface{}{
		"balance": balance,
	})
}