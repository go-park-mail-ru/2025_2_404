package balance

import (
	modelwallet "2025_2_404/internal/domain/models/client_wallet"
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
)

type repositoryI interface{
	Show(ctx context.Context, clientID modeluser.ID) (modelwallet.Balance, error)
}

type UseCase struct{
	repo repositoryI
}

func New(repo repositoryI) *UseCase{
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) Show(ctx context.Context, clientID modeluser.ID) (modelwallet.Balance, error){
	return u.repo.Show(ctx, clientID)
}