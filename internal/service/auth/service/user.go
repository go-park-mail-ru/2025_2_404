package service

import (
	modeluser "2025_2_404/internal/service/auth/domain"
	"2025_2_404/pkg/globalerrors"
	"context"
	"errors"
	"log"

	// "log"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type repositoryI interface {
	Create(ctx context.Context, user *modeluser.User) (modeluser.ID, error)
	FindByEmail(ctx context.Context, email string) (modeluser.User, error)
}

type tokenUsecaseI interface {
	GenerateToken(userID modeluser.ID) (string, error)
	// InvalidateToken(tokenString string) (string, error)
}

type UseCase struct {
	repo         repositoryI
	tokenUsecase tokenUsecaseI
	logger       *zap.Logger
}

func New(repo repositoryI, tokenUsecase tokenUsecaseI, logger *zap.Logger) *UseCase {
	return &UseCase{
		repo:         repo,
		tokenUsecase: tokenUsecase,
		logger:       logger,
	}
}

func (r *UseCase) Register(ctx context.Context, email, password, userName string) (string, modeluser.ID, error) {
	user, err := modeluser.ValidateRegisterUser(userName, email, password)
	if err != nil {
		log.Println("Не валидированный пользователь %w", err)
		return "", uuid.Nil, err
	}

	userID, err := r.repo.Create(ctx, user)
	if err != nil {
		log.Println("Траблы с созданием пользвоателя %w", err)
		return "", uuid.Nil, err
	}

	token, err := r.tokenUsecase.GenerateToken(userID)
	if err != nil {
		log.Println("Не получилось создать токен, ошибка %w", err)
		return "", uuid.Nil, err
	}
	return token, userID, nil
}

func (u *UseCase) Check(ctx context.Context, email, password string) (modeluser.ID, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, globalerrors.ErrUserNotFound) {
			return uuid.Nil, globalerrors.ErrWrongEmailOrPassword
		}
		return uuid.Nil, globalerrors.ErrInternal
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return uuid.Nil, globalerrors.ErrWrongEmailOrPassword
	}

	return user.ID, nil
}

func (u *UseCase) Login(ctx context.Context, email string, password string) (string, modeluser.ID, error) {
	err := modeluser.ValidateLoginUser(email, password)
	if err != nil {
		log.Println("Валидация пароля или emaik не прошла, ошибка валидейт логин  %w", err)
		return "", uuid.Nil, globalerrors.ErrInvalidCredentials
	}

	userID, err := u.Check(ctx, email, password)
	if err != nil {
		log.Println("Валидация пароля или emaik не прошла, ошибка чек %w", err)
		return "", uuid.Nil, globalerrors.ErrWrongEmailOrPassword
	}

	token, err := u.tokenUsecase.GenerateToken(userID)
	if err != nil {
		log.Println("Токен не сгенерировался, %w", err)
		return "", uuid.Nil, err
	}

	return token, userID, nil
}

// func (u *UseCase) Logout(ctx context.Context, token string) (error) {

// }
