package http

import (
	"context"
	"io"
	"net/http"
	"time"

	"2025_2_404/pkg"
	"2025_2_404/pkg/utils"
	pbProfile "2025_2_404/protos/profile"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ProfileHandler struct {
	client pbProfile.ProfileClient
}

func NewProfileHandler(client pbProfile.ProfileClient) *ProfileHandler {
	return &ProfileHandler{client: client}
}

func (h *ProfileHandler) Show(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	resp, err := h.client.Show(ctx, &pbProfile.ShowRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	pkg.JSONResponse(w, http.StatusOK, resp)
}

func parseMultipartForm(r *http.Request) error {
	const maxMemory = 32 << 20 // 32 MB
	return r.ParseMultipartForm(maxMemory)
}

func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	if err := parseMultipartForm(r); err != nil {
		http.Error(w, `{"error":"failed to parse form"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	req := &pbProfile.UpdateRequest{
		UserName:    r.FormValue("user_name"),
		Email:       r.FormValue("email"),
		Password:    r.FormValue("password"),
		FirstName:   r.FormValue("first_name"),
		LastName:    r.FormValue("last_name"),
		Company:     r.FormValue("company"),
		Phone:       r.FormValue("phone"),
		ProfileType: r.FormValue("profile_type"),
	}

	if _, fileHeader, err := r.FormFile("avatar"); err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, `{"error":"invalid avatar file"}`, http.StatusBadRequest)
			return
		}
		defer file.Close()

		bytesData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, `{"error":"failed to read avatar"}`, http.StatusInternalServerError)
			return
		}
		req.Avatar = bytesData
	}

	resp, err := h.client.Update(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	pkg.JSONResponse(w, http.StatusOK, resp)
}

func (h *ProfileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	_, err := h.client.Delete(ctx, &pbProfile.DeleteRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}