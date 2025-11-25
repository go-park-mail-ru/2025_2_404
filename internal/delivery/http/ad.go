package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"2025_2_404/pkg/utils"
	pbAd "2025_2_404/protos/gen/go/ad" 
)

type AdHandler struct {
	client pbAd.AdServClient
}

func NewAdHandler(client pbAd.AdServClient) *AdHandler {
	return &AdHandler{client: client}
}

func (h *AdHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/ads")
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

	req := &pbAd.CreateRequest{
		Ad: &pbAd.Ad{
			Title:     dto.Title,
			Content:   dto.Content,
			Targeturl: dto.TargetURL,
		},
	}

	resp, err := h.client.Create(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
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