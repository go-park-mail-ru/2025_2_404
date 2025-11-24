package http

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	pb "2025_2_404/protos/profile"
)

type ProfileHandler struct {
	client pb.ProfileClient
}

func NewProfileHandler(client pb.ProfileClient) *ProfileHandler {
	return &ProfileHandler{
		client: client,
	}
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

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	}

	resp, err := h.client.Show(ctx, &pb.ShowRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
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

	req := &pb.UpdateRequest{
		UserName:  c.PostForm("user_name"),
		Email:     c.PostForm("email"),
		FirstName: c.PostForm("user_first_name"),
		LastName:  c.PostForm("user_second_name"),
		Company:   c.PostForm("company"),
		Phone:     c.PostForm("phone_number"),
	}

	fileHeader, err := c.FormFile("img")
	if err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad image file"})
			return
		}
		defer file.Close()

		imgBytes, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read image"})
			return
		}
		
		req.Avatar = imgBytes
	}

	resp, err := h.client.Update(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
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

	_, err := h.client.Delete(ctx, &pb.DeleteRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.Status(http.StatusNoContent)
}

func HTTPStatusFromCode(code interface{}) int {
	return http.StatusInternalServerError
} 