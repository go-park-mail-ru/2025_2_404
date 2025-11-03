package profile

import (
	"context"
	modeluser "2025_2_404/internal/domain/models/user"
)

type repositoryI interface{
	Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, error)
	Update(ctx context.Context, client modeluser.User) error
}

type UseCase struct{
	repo repositoryI
}

func New(repo repositoryI) *UseCase{
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) Update(ctx context.Context, client modeluser.User) error{
	return u.repo.Update(ctx, client)
}

func (u *UseCase) Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, error){
	return u.repo.Show(ctx, clientID)
}

