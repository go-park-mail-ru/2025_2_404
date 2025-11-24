package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
	pbAuth "2025_2_404/protos/auth"
)

type AuthHandler struct {
	client pbAuth.AuthClient
}

func NewAuthHandler(client pbAuth.AuthClient) *AuthHandler {
	return &AuthHandler{
		client: client,
	}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/auth")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)
	}
}

type loginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type registerDTO struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	UserName  string `json:"user_name" binding:"required"` 
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input registerDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.Register(ctx, &pbAuth.RegisterRequest{
		Email:    input.Email,
		Password: input.Password,
		UserName: input.UserName,
	})

	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input loginDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.Login(ctx, &pbAuth.LoginRequest{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func HTTPStatusFromCode(code interface{}) int {
	return http.StatusInternalServerError 
}