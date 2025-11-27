package http

import (
	"2025_2_404/pkg"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"2025_2_404/pkg/utils"
	pbAd "2025_2_404/protos/gen/go/ad"
	pbStorage "2025_2_404/protos/gen/go/storage"
)

type AdHandler struct {
	client        pbAd.AdServClient
	storageClient pbStorage.StorageClient
}

func NewAdHandler(client pbAd.AdServClient, storageClient pbStorage.StorageClient) *AdHandler {
	return &AdHandler{client: client, storageClient: storageClient}
}

type adDTO struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	TargetURL string `json:"target_url"`
}

func (h *AdHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := parseMultipartForm(r); err != nil {
		http.Error(w, `{"error":"invalid form"}`, http.StatusBadRequest)
		return
	}

	var fileBytes []byte
	var imageFilename string

	if _, fileHeader, err := r.FormFile("image"); err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, `{"error":"cannot open image"}`, http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileBytes, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, `{"error":"cannot read image"}`, http.StatusInternalServerError)
			return
		}

		ext := filepath.Ext(fileHeader.Filename)
		if ext == "" {
			ext = ".jpg"
		}
		imageFilename = "storage/ad/" + uuid.New().String() + ext
	} else if !isMissingFileError(err) {
		http.Error(w, `{"error":"invalid image"}`, http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	targetURL := r.FormValue("target_url")

	if title == "" || content == "" || targetURL == "" {
		http.Error(w, `{"error":"title, content and target_url are required"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	req := &pbAd.CreateRequest{
		Ad: &pbAd.Ad{
			Title:     title,
			Content:   content,
			Targeturl: targetURL,
		},
	}

	resp, err := h.client.Create(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	if len(fileBytes) > 0 {
		storageCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := h.storageClient.Create(storageCtx, &pbStorage.CreateRequest{
			ImagePath: imageFilename,
			ImageData: fileBytes,
		})
		if err != nil {
			log.Printf("Warning: ad created but image upload failed: %v", err)
			pkg.JSONResponse(w, http.StatusCreated, map[string]interface{}{
				"ad": resp,
			})
			return
		}
	}

	pkg.JSONResponse(w, http.StatusCreated, resp)
}

func isMissingFileError(err error) bool {
	return err == http.ErrMissingFile
}

func (h *AdHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	resp, err := h.client.GetAllAds(ctx, &pbAd.GetAllAdsRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	pkg.JSONResponse(w, http.StatusOK, resp.Ads)
}

func (h *AdHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	resp, err := h.client.GetAd(ctx, &pbAd.GetAdRequest{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		code := utils.HTTPStatusFromCode(st.Code())
		if code == http.StatusNotFound {
			http.Error(w, `{"error":"ad not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"`+st.Message()+`"}`, code)
		}
		return
	}

	pkg.JSONResponse(w, http.StatusOK, resp.Ad)
}

func (h *AdHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var dto adDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	req := &pbAd.UpdateRequest{
		Ad: &pbAd.Ad{
			Id:        id,
			Title:     dto.Title,
			Content:   dto.Content,
			Targeturl: dto.TargetURL,
		},
	}

	_, err := h.client.Update(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AdHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	_, err := h.client.Delete(ctx, &pbAd.DeleteRequest{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}