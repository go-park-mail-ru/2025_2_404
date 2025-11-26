package http

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"2025_2_404/pkg/utils"
	pbAd "2025_2_404/protos/gen/go/ad"
	pbStorage "2025_2_404/protos/gen/go/storage"
)

type AdHandler struct {
	client pbAd.AdServClient
	storageClient pbStorage.StorageClient
}

func NewAdHandler(client pbAd.AdServClient, storageClient pbStorage.StorageClient) *AdHandler {
	return &AdHandler{client: client, storageClient: storageClient}
}

func (h *AdHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/ads")
	{
		api.POST("", h.Create)
		api.GET("", h.GetAll)
		api.GET("/:id", h.GetOne)
		api.PUT("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

type adDTO struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	TargetURL string `json:"target_url"`
}

func (h *AdHandler) Create(c *gin.Context) {
	var fileBytes []byte
	var imageFilename string

	fileHeader, err := c.FormFile("image")
	if err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot open image"})
			return
		}
		defer file.Close()

		fileBytes, err = io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read image"})
			return
		}

		extension := filepath.Ext(fileHeader.Filename)
		if extension == "" {
			extension = ".jpg" // по умолчанию
		}
		uuid := uuid.New().String()
		imageFilename = "storage/ad/" + uuid + extension
	} else if !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image"})
		return
	}
	
	title := c.PostForm("title")
	content := c.PostForm("content")
	targetURL := c.PostForm("target_url")

	if title == "" || content == "" || targetURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title, content and target_url are required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
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
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	if len(fileBytes) > 0 {
		storageCtx, storageCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer storageCancel()

		_, err := h.storageClient.Create(storageCtx, &pbStorage.CreateRequest{
			ImagePath: imageFilename,
			ImageData: fileBytes,
		})
		if err != nil {
			log.Printf("Warning: ad created but image upload failed: %v", err)
			c.JSON(http.StatusCreated, gin.H{
				"ad":       resp,
				"warning":  "image upload failed",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AdHandler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	}
	
	req := &pbAd.GetAllAdsRequest{}

	resp, err := h.client.GetAllAds(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusOK, resp.Ads)
}

func (h *AdHandler) GetOne(c *gin.Context) {
	id := c.Param("id")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	}

	req := &pbAd.GetAdRequest{Id: id}
	resp, err := h.client.GetAd(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusOK, resp.Ad)
}

func (h *AdHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var dto adDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	}

	req := &pbAd.UpdateRequest{
		Ad: &pbAd.Ad{
			Id:        id,
			Title:     dto.Title,
			Content:   dto.Content,
			Targeturl: dto.TargetURL,
		},
	}

	resp, err := h.client.Update(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AdHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	}

	req := &pbAd.DeleteRequest{Id: id}
	_, err := h.client.Delete(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.Status(http.StatusNoContent)
}