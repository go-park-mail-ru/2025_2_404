package handler

import (
	"2025_2_404/internal/service/auth/service"
	"2025_2_404/protos/auth"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	auth.UnimplementedAuthServer 
	useCase *service.UseCase
}

func NewAuthServer(useCase *service.UseCase) *AuthServer {
	return &AuthServer{
		useCase: useCase,
	}
}

func (s *AuthServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	token, userID, err := s.useCase.Register(ctx, req.Email, req.Password, req.UserName)
	if err != nil {
		return nil, status.Error(codes.Internal, "registration failed")
	}

	return &auth.RegisterResponse{
		Token:  token,
		UserId: userID.String(),
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	token, userID, err := s.useCase.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}

	return &auth.LoginResponse{
		Token:  token,
		UserId: userID.String(),
	}, nil
}
