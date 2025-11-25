package http

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"

	"2025_2_404/pkg/utils"
	pbStorage "2025_2_404/protos/gen/go/storage"
)

type StorageHandler struct {
	client pbStorage.StorageClient
}

func NewStorageHandler(client pbStorage.StorageClient) *StorageHandler {
	return &StorageHandler{client: client}
}

func (h *StorageHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/profile")
	{
		api.GET("", h.Get)
		api.POST("", h.Create) 
		api.DELETE("", h.Delete)
	}
}
func (h *StorageHandler) Get(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path query param is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.Get(ctx, &pbStorage.GetRequest{
		ImagePath: path,
	})

	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	if len(resp.ImageData) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	contentType := resp.ContentType
	if contentType == "" {
		contentType = http.DetectContentType(resp.ImageData)
	}

	c.Data(http.StatusOK, contentType, resp.ImageData)
}

func (h *StorageHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	path := c.PostForm("path") 
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot open file"})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read file"})
		return
	}

	_, err = h.client.Create(ctx, &pbStorage.CreateRequest{
		ImagePath: path,
		ImageData: fileBytes,
	})

	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "file uploaded"})
}

func (h *StorageHandler) Delete(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.client.Delete(ctx, &pbStorage.DeleteRequest{
		ImagePath: path,
	})

	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.Status(http.StatusNoContent)
}