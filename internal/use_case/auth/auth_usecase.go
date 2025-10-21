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
	FindUserByEmail(ctx context.Context, email string) (modeluser.User, error)
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
