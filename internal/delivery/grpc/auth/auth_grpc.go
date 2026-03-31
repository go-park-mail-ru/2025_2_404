// Package handler provides gRPC delivery handlers for authentication-related operations.
package handler

import (
	modeluser "2025_2_404/internal/service/auth/domain"
	"2025_2_404/pkg/utils"
	"2025_2_404/protos/gen/go/auth"
	"context"
	"errors"
	"log"

	"go.uber.org/zap"
)

var (
	ErrAlreadyExists = errors.New("user already exists")
	ErrNotFound      = errors.New("user not found")
)

type UseCase interface {
	Register(ctx context.Context, email, password, userName string) (string, modeluser.ID, error)
	Login(ctx context.Context, email string, password string) (string, modeluser.ID, error)
	// ValidateToken(ctx context.Context, tokenString string) (modeluser.ID, error)
}

type UseCaseJWT interface {
	ValidateToken(ctx context.Context, tokenString string) (modeluser.ID, error)
}

type AuthServer struct {
	auth.UnimplementedAuthServer
	useCase    UseCase
	useCaseJWT UseCaseJWT
	logger     *zap.Logger
}

func NewAuthServer(useCase UseCase, useCaseJWT UseCaseJWT, logger *zap.Logger) *AuthServer {
	return &AuthServer{
		useCase:    useCase,
		useCaseJWT: useCaseJWT,
		logger:     logger,
	}
}

func (s *AuthServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	token, userID, err := s.useCase.Register(ctx, req.Email, req.Password, req.UserName)
	if err != nil {
		log.Printf("Register error: %v", err)
		return nil, utils.ToGRPCError(err) // ← преобразуем в gRPC-статус
	}
	return &auth.RegisterResponse{Token: token, UserId: userID.String()}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	token, userID, err := s.useCase.Login(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return nil, utils.ToGRPCError(err)
	}
	return &auth.LoginResponse{Token: token, UserId: userID.String()}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *auth.TokenRequest) (*auth.TokenResponse, error) {
	userID, err := s.useCaseJWT.ValidateToken(ctx, req.Token)
	if err != nil {
		log.Printf("Token validation error: %v", err)
		return nil, utils.ToGRPCError(err)
	}
	return &auth.TokenResponse{Valid: true, UserId: userID.String()}, nil
}
