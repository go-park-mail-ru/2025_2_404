package usecase

import(
	modeluser "2025_2_404/internal/models/user"
	"context"
	"fmt"
	"2025_2_404/internal/utils"
)

type repositoryI interface {
	CreateUser(ctx context.Context, email, password, userName string) (*modeluser.User, error)
	CreateSession(ctx context.Context, userID, sessionID string) (string, error)
}

type AuthUseCase struct {
	repo repositoryI
}

func (r *AuthUseCase) RegisterUser(ctx context.Context, email, password, userName string) (*modeluser.ID, error) {
	user, err := modeluser.NewUser(userName, email, password)
	if err != nil {
		return nil, fmt.Errorf("not validate user: %w", err)
	}
	user, err = r.repo.CreateUser(ctx, user.Email, user.HashedPassword, user.UserName)
	if err != nil {
		return nil, fmt.Errorf("problem with repository CreateUser: %w", err)
	}
	return &user.ID, nil
}

func (r *AuthUseCase) SessionGenerateAndSave(ctx context.Context, userID string) (string, error) {
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
