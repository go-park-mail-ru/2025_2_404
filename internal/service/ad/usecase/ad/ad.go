package ad

import (
	modelad "2025_2_404/internal/service/ad/domain/ad"
	modelfullad "2025_2_404/internal/service/ad/domain/ad_full_info"
	modeluser "2025_2_404/internal/service/ad/domain/user"
	"context"
	"fmt"
)

type adRepositoryI interface {
	FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelad.Ads, error)
	Create(ctx context.Context, ad modelad.Ads) (error)
	GetOneAd(ctx context.Context, adID modelad.ID) (modelfullad.AdFullInfo, error)
	Update(ctx context.Context, ad modelad.Ads) error
	Delete(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error
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

func (u *UseCase) Create(ctx context.Context, ad modelad.Ads) (error) {
	return u.adRepo.Create(ctx, ad)
}

func (u *UseCase) Update(ctx context.Context, ad modelad.Ads) error {
	return u.adRepo.Update(ctx, ad)
}

func (u *UseCase) Delete(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error {
	return u.adRepo.Delete(ctx, adID, clientID)
}

func (u *UseCase) GetOneAd(ctx context.Context, adID modelad.ID) (modelfullad.AdFullInfo, int, error) {
	adInfo, err := u.adRepo.GetOneAd(ctx, adID)
	conversion := -1
	if err != nil {
		return modelfullad.AdFullInfo{}, conversion, fmt.Errorf("Failed to get ad with id error %w", err)
	}
	if adInfo.Impressions != 0{
		conversion = adInfo.Clicks / adInfo.Impressions
	}

	return adInfo, conversion, nil
}
