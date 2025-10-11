package usecase

import (
	"context"
	modeluser "2025_2_404/internal/domain/models/user"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

type userRepository interface {
	FindUserByEmail(ctx context.Context, email string) (modeluser.User, error)
	FindSessionByUserID(ctx context.Context, userID modeluser.ID) (string, error)
	FindUserBySessionID(ctx context.Context, sessionID int) (modeluser.ID, error)
}

type userUsecase struct {
	repo userRepository
}

func (u *userUsecase) CheckUser(ctx context.Context, email string, password string) (modeluser.User, error) {
	user, err := u.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return modeluser.User{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return modeluser.User{}, err
	}
	if user.HashedPassword != string(hashedPassword) {
		return modeluser.User{}, errors.New("invalid password")
	}
	return user, nil
}

func (u *userUsecase) FindSession(ctx context.Context, userID modeluser.ID) (string, error) {
	session, err := u.repo.FindSessionByUserID(ctx, userID)
	if err != nil {
		return "", err
	}
	return session, nil
}

func (u *userUsecase) FindUser(ctx context.Context, sessionID int) (modeluser.ID, error) {
	userID, err := u.repo.FindUserBySessionID(ctx, sessionID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}