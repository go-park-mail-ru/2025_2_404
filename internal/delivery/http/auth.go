package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"

	"2025_2_404/pkg/utils"
	pbAuth "2025_2_404/protos/auth"
)

type AuthHandler struct {
	client pbAuth.AuthClient
}

func NewAuthHandler(client pbAuth.AuthClient) *AuthHandler {
	return &AuthHandler{client: client}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/auth")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)
	}
}

type registerDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	UserName string `json:"user_name" binding:"required"`
}

type loginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.Register(ctx, &pbAuth.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		UserName: req.UserName,
	})

	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.Login(ctx, &pbAuth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusOK, resp)
}