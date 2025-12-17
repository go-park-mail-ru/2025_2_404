// Package http provides HTTP delivery handlers for advertisement-related operations.
package http

import (
	adfullinfo "2025_2_404/internal/service/ad/domain/ad_full_info"
	"2025_2_404/pkg"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	pkgfile "2025_2_404/pkg/readerFile"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"2025_2_404/pkg/utils"
	pbAd "2025_2_404/protos/gen/go/ad"
	pbProfile "2025_2_404/protos/gen/go/profile"
	pbStorage "2025_2_404/protos/gen/go/storage"
)

type AdHandler struct {
	client        pbAd.AdServClient
	storageClient pbStorage.StorageClient
	profileClient pbProfile.ProfileClient
}

func NewAdHandler(client pbAd.AdServClient, storageClient pbStorage.StorageClient, profileClient pbProfile.ProfileClient) *AdHandler {
	return &AdHandler{client: client, storageClient: storageClient, profileClient: profileClient}
}

func (h *AdHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := parseMultipartForm(r); err != nil {
		http.Error(w, `{"error":"invalid form"}`, http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	targetURL := r.FormValue("target_url")
	budgetStr := r.FormValue("budget")

	if title == "" || content == "" || targetURL == "" || budgetStr == "" {
		http.Error(w, `{"error":"title, content, target_url and budget are required"}`, http.StatusBadRequest)
		return
	}

	budget, err := strconv.ParseUint(budgetStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"invalid budget format"}`, http.StatusBadRequest)
		return
	}

	fileBytes, imageFilename, err := pkgfile.ExtractImage(r, "ad/", "image")
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
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
			ImgPath:   imageFilename,
			Budget:    uint32(budget),
		},
	}

	if len(fileBytes) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := h.storageClient.Create(ctx, &pbStorage.CreateRequest{
			ImagePath: imageFilename,
			ImageData: fileBytes,
		})
		if err != nil {
			log.Printf("ERROR: async image upload failed for %s: %v", imageFilename, err)
		} else {
			log.Printf("INFO: image uploaded successfully: %s", imageFilename)
		}
	}

	resp, err := h.client.Create(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	pkg.JSONResponse(w, http.StatusCreated, "Ad created successfully", resp)
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

	pkg.JSONResponse(w, http.StatusOK, "Ads retrieved successfully", resp.Ads)
}

func (h *AdHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	adProto, err := h.client.GetAd(ctx, &pbAd.GetAdRequest{Id: id})
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

	adID, err := uuid.Parse(adProto.Ad.GetId())
	if err != nil {
		http.Error(w, `{"error":"invalid adId format"}`, http.StatusBadRequest)
	}

	startAt, err := time.Parse(time.RFC3339, adProto.Ad.GetStartAt())
	if err != nil {
		http.Error(w, `{"error":"invalid startAt format"}`, http.StatusBadRequest)
	}

	endAt, err := time.Parse(time.RFC3339, adProto.Ad.GetEndAt())
	if err != nil {
		http.Error(w, `{"error":"invalid endAt format"}`, http.StatusBadRequest)
	}

	adFullResp := adfullinfo.AdFullInfo{
		ID:          adID,
		Title:       adProto.Ad.GetTitle(),
		Content:     adProto.Ad.GetContent(),
		ImgPath:     adProto.Ad.GetImgPath(),
		TargetURL:   adProto.Ad.GetTargeturl(),
		Budget:      adProto.Ad.GetBudget(),
		Status:      adProto.Ad.GetStatus(),
		StartAt:     startAt,
		EndAt:       endAt,
		Clicks:      int(adProto.Ad.GetClicks()),
		Impressions: int(adProto.Ad.GetImpressions()),
	}

	imgPath := adProto.Ad.ImgPath
	var resIMG *pbStorage.GetResponse
	if imgPath != "" {
		resIMG, err = h.storageClient.Get(ctx, &pbStorage.GetRequest{
			ImagePath: imgPath,
		})
		if err != nil {
			log.Printf("ERROR: failed to get presigned URL for %s: %v", imgPath, err)
		}
	}

	pkg.JSONResponse(w, http.StatusOK, "Ad retrieved successfully", map[string]interface{}{
		"ad":        adFullResp,
		"imageData": resIMG,
	})
}

func (h *AdHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}
	if err := parseMultipartForm(r); err != nil {
		http.Error(w, `{"error":"invalid form"}`, http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	targetURL := r.FormValue("target_url")
	statusAd := r.FormValue("status")

	if title == "" || content == "" || targetURL == "" || statusAd == "" {
		http.Error(w, `{"error":"title, content, target_url and budget are required"}`, http.StatusBadRequest)
		return
	}

	var fileBytes []byte
	var newImageFilename string
	var err error

	if r.FormValue("image") != "" {
		fileBytes, newImageFilename, err = pkgfile.ExtractImage(r, "ad/", "image")
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
	}

	// ———— ШАГ 2: Загружаем изображение (если есть) ————
	if len(fileBytes) > 0 {
		uploadCtx, uploadCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer uploadCancel()

		_, err := h.storageClient.Create(uploadCtx, &pbStorage.CreateRequest{
			ImagePath: newImageFilename,
			ImageData: fileBytes,
		})
		if err != nil {
			log.Printf("ERROR: async image upload failed for %s: %v", newImageFilename, err)
		} else {
			log.Printf("INFO: image uploaded successfully: %s", newImageFilename)
		}
	}

	// ———— ШАГ 3: Вызываем gRPC Update ————
	req := &pbAd.UpdateRequest{
		Ad: &pbAd.Ad{
			Id:        id,
			Title:     title,
			Content:   content,
			Targeturl: targetURL,
			ImgPath:   newImageFilename,
			Status:    statusAd,
		},
	}

	updateResp, err := h.client.Update(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "Ad updated successfully", updateResp)
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

func (h *AdHandler) UpdateBudget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}

	budgetStr := r.FormValue("budget")
	if budgetStr == "" {
		http.Error(w, `{"error":"budget is required"}`, http.StatusBadRequest)
		return
	}

	budget, err := strconv.ParseUint(budgetStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"invalid budget format"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	_, err = h.profileClient.SubtractBalance(ctx, &pbProfile.SubtractBalanceRequest{
		SubAmount: uint32(budget),
	})
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	req := &pbAd.UpdateBudgetRequest{
		Id:     id,
		Budget: uint32(budget),
	}

	resp, err := h.client.UpdateAdBudget(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	pkg.JSONResponse(w, http.StatusOK, "Budget updated successfully", resp)
}
