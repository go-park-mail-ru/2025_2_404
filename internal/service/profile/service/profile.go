package profile

import (
	modeluser "2025_2_404/internal/service/profile/domain"
	"context"
	"fmt"
)

type repositoryI interface{
	Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, error)
	Update(ctx context.Context, client modeluser.User) error
	Delete(ctx context.Context, clientID modeluser.ID) error
}

type UseCase struct{
	repo repositoryI
}

func New(repo repositoryI) *UseCase{
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) Update(ctx context.Context, client modeluser.User) error {
	return u.repo.Update(ctx, client)
}

func (u *UseCase) Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, error){
	return u.repo.Show(ctx, clientID)
}

func (u *UseCase) Delete(ctx context.Context, clientID modeluser.ID) error {

	if err := u.repo.Delete(ctx, clientID); err != nil {
		return fmt.Errorf("failed to delete user from repository: %w", err)
	}

	return nil
}