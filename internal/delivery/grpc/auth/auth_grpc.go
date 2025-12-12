package handler

import (
	modeluser "2025_2_404/internal/service/auth/domain"
	"2025_2_404/protos/gen/go/auth"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	useCase UseCase
	useCaseJWT UseCaseJWT
}

func NewAuthServer(useCase UseCase, useCaseJWT UseCaseJWT) *AuthServer {
	return &AuthServer{
		useCase: useCase,
		useCaseJWT: useCaseJWT,
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

func (s *AuthServer) ValidateToken(ctx context.Context, req *auth.TokenRequest) (*auth.TokenResponse, error) {
	userID, err := s.useCaseJWT.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return &auth.TokenResponse{
		Valid: true,
		UserId:  userID.String(),
		Error: "",
	}, nil
}