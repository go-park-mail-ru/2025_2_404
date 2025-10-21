package ad

import (
	modelad "2025_2_404/internal/domain/models/ad"
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
)

type adRepositoryI interface {
	FindAdByUserID(ctx context.Context, userID modeluser.ID) (modelad.Ads, error)
	CreateAd(ctx context.Context, ad modelad.Ads) (int, error)
}

type UseCase struct {
	adRepo adRepositoryI
}

func New(adRepo adRepositoryI) *UseCase {
	return &UseCase{
		adRepo: adRepo,
	}
}

func (u *UseCase) FindAdByUserID(ctx context.Context, userID modeluser.ID) (modelad.Ads, error) {
	return u.adRepo.FindAdByUserID(ctx, userID)
}

func (u *UseCase) CreateAd(ctx context.Context, ad modelad.Ads) (int, error) {
	return u.adRepo.CreateAd(ctx, ad)
}


