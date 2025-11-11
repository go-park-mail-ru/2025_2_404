package feed

import (
	modelad "2025_2_404/internal/domain/models/ad"
	"context"
)

type repositoryI interface{
	Give(ctx context.Context, platformName string) ([]modelad.Ads, error)
}

type UseCase struct{
	repo repositoryI
}

func New(repo repositoryI) *UseCase{
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) Give(ctx context.Context, platformName string) ([]modelad.Ads, error) {
	return u.repo.Give(ctx, platformName)
}
