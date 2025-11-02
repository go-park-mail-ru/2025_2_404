package ad

import (
	modelad "2025_2_404/internal/domain/models/ad"
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
)

type adRepositoryI interface {
	FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelad.Ads, error)
	Create(ctx context.Context, ad modelad.Ads) (modelad.Ads, error)
	Update(ctx context.Context, ad modelad.Ads) error
	Delete(ctx context.Context, adID int64) error
}

type UseCase struct {
	adRepo adRepositoryI
}

func New(adRepo adRepositoryI) *UseCase {
	return &UseCase{
		adRepo: adRepo,
	}
}

func (u *UseCase) FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelad.Ads, error) {
	return u.adRepo.FindByUserID(ctx, userID)
}

func (u *UseCase) Create(ctx context.Context, ad modelad.Ads) (modelad.Ads, error) {
	return u.adRepo.Create(ctx, ad)
}

func (u *UseCase) Update(ctx context.Context, ad modelad.Ads) error{
	return u.adRepo.Update(ctx, ad)
}

func (u *UseCase) Delete(ctx context.Context, adID int64) error{
	return u.adRepo.Delete(ctx, adID)
}

