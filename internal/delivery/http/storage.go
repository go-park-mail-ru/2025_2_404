package http

import (
	"context"
	"io"
	"net/http"
	"time"

	"2025_2_404/pkg"
	"2025_2_404/pkg/utils"
	pbStorage "2025_2_404/protos/gen/go/storage"

	"google.golang.org/grpc/status"
)

type StorageHandler struct {
	client pbStorage.StorageClient
}

func NewStorageHandler(client pbStorage.StorageClient) *StorageHandler {
	return &StorageHandler{client: client}
}

func (h *StorageHandler) Get(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, `{"error":"path query param is required"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.Get(ctx, &pbStorage.GetRequest{ImagePath: path})
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	if len(resp.ImageData) == 0 {
		http.Error(w, `{"error":"image not found"}`, http.StatusNotFound)
		return
	}

	contentType := resp.ContentType
	if contentType == "" {
		contentType = http.DetectContentType(resp.ImageData)
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(resp.ImageData)
}

func (h *StorageHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := parseMultipartForm(r); err != nil {
		http.Error(w, `{"error":"invalid form"}`, http.StatusBadRequest)
		return
	}

	path := r.FormValue("path")
	if path == "" {
		http.Error(w, `{"error":"path is required"}`, http.StatusBadRequest)
		return
	}

	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"file is required"}`, http.StatusBadRequest)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		http.Error(w, `{"error":"cannot open file"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, `{"error":"cannot read file"}`, http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = h.client.Create(ctx, &pbStorage.CreateRequest{
		ImagePath: path,
		ImageData: fileBytes,
	})
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	pkg.JSONResponse(w, http.StatusCreated, nil)
}

func (h *StorageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, `{"error":"path is required"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.client.Delete(ctx, &pbStorage.DeleteRequest{ImagePath: path})
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}