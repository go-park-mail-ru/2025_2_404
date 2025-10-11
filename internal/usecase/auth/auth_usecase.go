package usecase

<<<<<<< HEAD
import (
	"2025_2_404/internal/models"
	modeluser "2025_2_404/internal/models"
	"errors"
=======
import(
	modeluser "2025_2_404/internal/models/user"
	"context"
	"fmt"
	"2025_2_404/internal/utils"
>>>>>>> a8230ea6cc45a4ef7d6d317222973fdc7959bd18
)

type repositoryI interface {
	CreateUser(ctx context.Context, email, password, userName string) (*modeluser.User, error)
	CreateSession(ctx context.Context, userID modeluser.ID, sessionID string) (string, error)
}

type AuthUseCase struct {
	repo repositoryI
}

<<<<<<< HEAD
type UserUseCase struct {
	userRepo modeluser.UserRepository
}

// dto - data transport object 

func (uc *UserUseCase) Register(dto RegisterUser) (*models.User, error) {
	user, err := models.NewUser(dto.UserName, dto.Email, dto.Password)
=======
func (r *AuthUseCase) RegisterUser(ctx context.Context, email, password, userName string) (*modeluser.ID, error) {
	user, err := modeluser.NewUser(userName, email, password)
>>>>>>> a8230ea6cc45a4ef7d6d317222973fdc7959bd18
	if err != nil {
		return nil, fmt.Errorf("not validate user: %w", err)
	}
<<<<<<< HEAD

	return user, nil
}

func (uc *UserUseCase) Login(dto BaseUser) (*models.User, error) {
	user, err := models.LoginUser(dto.Email, dto.Password)
=======
	user, err = r.repo.CreateUser(ctx, user.Email, user.HashedPassword, user.UserName)
>>>>>>> a8230ea6cc45a4ef7d6d317222973fdc7959bd18
	if err != nil {
		return nil, fmt.Errorf("problem with repository CreateUser: %w", err)
	}
	return &user.ID, nil
}

func (r *AuthUseCase) SessionGenerateAndSave(ctx context.Context, userID modeluser.ID) (string, error) {
	sessionID, err := utils.GenerateSession()
	if err != nil {
		return "", fmt.Errorf("problem with session generation: %w", err)
	}
	sessionID, err = r.repo.CreateSession(ctx, userID, sessionID)
	if err != nil {
		return "", fmt.Errorf("problem with repository CreateSession: %w", err)
	}
	return sessionID, nil
}
