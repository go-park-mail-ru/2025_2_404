package auth

import (
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type repositoryI interface {
	Create(ctx context.Context, user *modeluser.User) (modeluser.ID, error)
	FindByEmail(ctx context.Context, email string) (modeluser.User, error)
}

type tokenUsecaseI interface {
	GenerateToken(userID modeluser.ID) (string, error)
}

type AuthUseCase struct {
	repo repositoryI
	tokenUsecase tokenUsecaseI
}

func New(repo repositoryI, tokenUsecase tokenUsecaseI) *AuthUseCase {
	return &AuthUseCase{
		repo: repo,
		tokenUsecase: tokenUsecase,
	}
}

func (r *AuthUseCase) Register(ctx context.Context, email, password, userName string) (string, error) {
	user, err := modeluser.NewUser(userName, email, password)
	if err != nil {
		return "", fmt.Errorf("not validate user: %w", err)
	}

	userID, err := r.repo.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("problem with repository CreateUser: %w", err)
	}

	token, err := r.tokenUsecase.GenerateToken(userID)
	if err != nil {
		return "", fmt.Errorf("auth_login : %w", err)
	}
	return token, nil
}

func (u *AuthUseCase) Check(ctx context.Context, email string, password string) (modeluser.ID, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return modeluser.ID(0), err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return modeluser.ID(0), errors.New("invalid password")
	}
	return user.ID, nil
}

func (u *AuthUseCase) Login(ctx context.Context, email string, password string) (string, error) {
	userID, err := u.Check(ctx, email, password)
	if err != nil {
		return "", err
	}

	token, err := u.tokenUsecase.GenerateToken(userID)
	if err != nil {
		return  "", fmt.Errorf("auth_login : %w", err)
	}

	return token, nil
}
