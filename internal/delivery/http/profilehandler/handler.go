package profilehandler

import (
	modeluser "2025_2_404/internal/domain/models/user"
	"2025_2_404/internal/modules"
	"2025_2_404/pkg"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type profileUsecaseI interface{
	Update(ctx context.Context, client modeluser.User) error
	Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, error)
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

	var err error
	var client modeluser.User
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusBadRequest)
		return
	}

	client.ID, err = modules.Get(r.Context())
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err = h.profileUsecase.Update(r.Context(), client); err != nil{
		http.Error(w, "client update faild", http.StatusUnprocessableEntity)
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "client update", map[string]interface{}{})
}

func (h *UserHandler) ShowHandler (w http.ResponseWriter, r *http.Request){
	clientID, error := modules.Get(r.Context())
	if error != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	client, err := h.profileUsecase.Show(r.Context(), clientID)
	if err != nil {
		http.Error(w, "Don't have ads this user", http.StatusInternalServerError)
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "Successful authorization", map[string]interface{}{
		"client":	client,
	})
}