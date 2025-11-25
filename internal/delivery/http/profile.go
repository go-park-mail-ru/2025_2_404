package http

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"2025_2_404/pkg/utils"
	pbProfile "2025_2_404/protos/profile"
)

type ProfileHandler struct {
	client pbProfile.ProfileClient
}

func NewProfileHandler(client pbProfile.ProfileClient) *ProfileHandler {
	return &ProfileHandler{client: client}
}

func (h *ProfileHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/profile")
	{
		api.GET("", h.Show)
		api.POST("/update", h.Update) 
		api.DELETE("", h.Delete)
	}
}

func (h *ProfileHandler) Show(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проброс токена
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	}

	resp, err := h.client.Show(ctx, &pbProfile.ShowRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ProfileHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	}

	req := &pbProfile.UpdateRequest{
		UserName:    c.PostForm("user_name"),
		Email:       c.PostForm("email"),
		Password:    c.PostForm("password"),
		FirstName:   c.PostForm("first_name"),
		LastName:    c.PostForm("last_name"),
		Company:     c.PostForm("company"),
		Phone:       c.PostForm("phone"),
		ProfileType: c.PostForm("profile_type"),
	}

	// Обработка файла (поле "avatar")
	fileHeader, err := c.FormFile("avatar")
	if err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid avatar file"})
			return
		}
		defer file.Close()

		bytesData, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read avatar"})
			return
		}
		req.Avatar = bytesData
	}

	resp, err := h.client.Update(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ProfileHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	}

	_, err := h.client.Delete(ctx, &pbProfile.DeleteRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.Status(http.StatusNoContent)
}