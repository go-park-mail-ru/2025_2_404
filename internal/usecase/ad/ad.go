package ad

import (
	modelad "2025_2_404/internal/models/ad"
	modeluser "2025_2_404/internal/models/user"
	"context"
)

type adRepositoryI interface {
	FindAdByUserID(ctx context.Context, userID modeluser.ID) (modelad.Ads, error)
	CreateAd(ctx context.Context, ad modelad.Ads) (int, error)
}

type AdUseCase struct {
	adRepo adRepositoryI
}

func New(adRepo adRepositoryI) *AdUseCase {
	return &AdUseCase{
		adRepo: adRepo,
	}
}

func (u *AdUseCase) FindAdByUserID(ctx context.Context, userID modeluser.ID) (modelad.Ads, error) {
	return u.adRepo.FindAdByUserID(ctx, userID)
}

func (u *AdUseCase) CreateAd(ctx context.Context, ad modelad.Ads) (int, error) {
	return u.adRepo.CreateAd(ctx, ad)
}

func (u *AdUseCase) UpdateAd(ctx context.Context, ad modelad.Ads) error {
	u.adRepo.FindAdByUserID(ctx, ad.)
	return u.adRepo.UpdateAd(ctx, ad)
}
