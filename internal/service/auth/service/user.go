package service

import (
	modeluser "2025_2_404/internal/service/auth/domain"
	"context"
	"errors"
	"fmt"
	"log"

	// "log"

	"github.com/google/uuid"
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
	repo repositoryI
	tokenUsecase tokenUsecaseI
}

func New(repo repositoryI, tokenUsecase tokenUsecaseI) *UseCase {
	return &UseCase{
		repo: repo,
		tokenUsecase: tokenUsecase,
	}
}

func (r *UseCase) Register(ctx context.Context, email, password, userName string) (string, modeluser.ID, error) {
	user, err := modeluser.ValidateRegisterUser(userName, email, password)
	if err != nil {
		// log.Printf("ОШИБКААА ПИЗДЦ")
		log.Println("Не валидированный пользователь %w", err)
		return "", uuid.Nil, fmt.Errorf("not validate user: %w", err)
	}

	userID, err := r.repo.Create(ctx, user)
	if err != nil {
		log.Println("Траблы с созданием пользвоателя %w", err)
		return "", uuid.Nil, fmt.Errorf("problem with repository CreateUser: %w", err)
	}

	token, err := r.tokenUsecase.GenerateToken(userID)
	if err != nil {
		log.Println("Не получилось создать токен, ошибка %w", err)
		return "", uuid.Nil,fmt.Errorf("auth_login : %w", err)
	}
	return token, userID, nil
}

func (u *UseCase) Check(ctx context.Context, email string, password string) (modeluser.ID, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		log.Println("Не валидированный пользователь %w", err)
		return uuid.Nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		log.Println("Неправильный пароль %w", err)
		return uuid.Nil, errors.New("invalid password")
	}
	return user.ID, nil
}

func (u *UseCase) Login(ctx context.Context, email string, password string) (string, modeluser.ID, error) {
	err := modeluser.ValidateLoginUser(email, password)
	if err != nil {
		log.Println("Валидация пароля или emaik не прошла, ошибка валидейт логин  %w", err)
		return "", uuid.Nil, fmt.Errorf("wrong validation of data: %w", err)
	}
	
	userID, err := u.Check(ctx, email, password)
	if err != nil {
		log.Println("Валидация пароля или emaik не прошла, ошибка чек %w", err)
		return "",uuid.Nil, err
	}

	token, err := u.tokenUsecase.GenerateToken(userID)
	if err != nil {
		log.Println("Токен не сгенерировался, %w", err)
		return  "",uuid.Nil, fmt.Errorf("auth_login : %w", err)
	}

	return token, userID, nil
}


// func (u *UseCase) Logout(ctx context.Context, token string) (error) {

// }