package service

import (
	modeluser "2025_2_404/internal/service/auth/internal/domain"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type repositoryI interface {
	Create(ctx context.Context, user *modeluser.User) (modeluser.ID, error)
	FindByEmail(ctx context.Context, email string) (modeluser.User, error)
}

type tokenUsecaseI interface {
	GenerateToken(userID modeluser.ID) (string, error)
	InvalidateToken(tokenString string) (string, error)
}

type UseCase struct {
	repo repositoryI
	tokenUsecase tokenUsecaseI
}

func New(repo repositoryI, tokenUsecase tokenUsecaseI) *UseCase {
	return &UseCase{
		repo: repo,
		tokenUsecase: tokenUsecase,
	}
}

func (r *UseCase) Register(ctx context.Context, email, password, userName string) (string, error) {
	user, err := modeluser.RegisterUser(userName, email, password)
	if err != nil {
		log.Printf("ОШИБКААА ПИЗДЦ")
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

func (u *UseCase) Check(ctx context.Context, email string, password string) (modeluser.ID, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return uuid.Nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return uuid.Nil, errors.New("invalid password")
	}
	return user.ID, nil
}

func (u *UseCase) Login(ctx context.Context, email string, password string) (string, error) {
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
