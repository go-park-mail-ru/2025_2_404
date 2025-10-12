package usecase

import(
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"fmt"
)

type repositoryI interface {
	CreateUser(ctx context.Context, user *modeluser.User) (modeluser.ID, error)
	CreateSession(ctx context.Context, userID modeluser.ID, sessionID string) (string, error)
	FindUserByEmail(ctx context.Context, email string) (modeluser.User, error)
	FindSessionByUserID(ctx context.Context, userID modeluser.ID) (string, error)
	FindUserBySessionID(ctx context.Context, sessionID string) (modeluser.ID, error)
}

type AuthUseCase struct {
	repo repositoryI
}

func New(repo repositoryI) *AuthUseCase {
	return &AuthUseCase{
		repo: repo,
	}
}

func (r *AuthUseCase) RegisterUser(ctx context.Context, email, password, userName string) (modeluser.ID, error) {
	user, err := modeluser.NewUser(userName, email, password)
	if err != nil {
		return 0, fmt.Errorf("not validate user: %w", err)
	}
	userID, err := r.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("problem with repository CreateUser: %w", err)
	}
	return userID, nil
}

// func (r *AuthUseCase) SessionGenerateAndSave(ctx context.Context, userID modeluser.ID) (string, error) {
// 	sessionID, err := utils.GenerateSession()
// 	if err != nil {
// 		return "", fmt.Errorf("problem with session generation: %w", err)
// 	}
// 	sessionID, err = r.repo.CreateSession(ctx, userID, sessionID)
// 	if err != nil {
// 		return "", fmt.Errorf("problem with repository CreateSession: %w", err)
// 	}
// 	return sessionID, nil
// }

func (u *AuthUseCase) CheckUser(ctx context.Context, email string, password string) (modeluser.User, error) {
	user, err := u.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return modeluser.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return modeluser.User{}, errors.New("invalid password")
	}
	return user, nil
}

func (u *AuthUseCase) FindSession(ctx context.Context, userID modeluser.ID) (string, error) {
	session, err := u.repo.FindSessionByUserID(ctx, userID)
	if err != nil {
		return "", err
	}
	return session, nil
}

func (u *AuthUseCase) FindUser(ctx context.Context, sessionID string) (modeluser.ID, error) {
	userID, err := u.repo.FindUserBySessionID(ctx, sessionID)
	if err != nil {
		return 0, err
	}
	return modeluser.ID(userID), nil
}
