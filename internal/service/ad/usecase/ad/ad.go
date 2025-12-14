package ad

import (
	modelad "2025_2_404/internal/service/ad/domain/ad"
	modelfullad "2025_2_404/internal/service/ad/domain/ad_full_info"
	modeluser "2025_2_404/internal/service/ad/domain/user"
	"context"
	"fmt"
)

type adRepositoryI interface {
	FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelfullad.AdFullInfo, error)
	Create(ctx context.Context, ad modelad.Ads) (error)
	GetOneAd(ctx context.Context, adID modelad.ID, clientID modeluser.ID) (modelfullad.AdFullInfo, error)
	Update(ctx context.Context, ad modelad.Ads) error
	Delete(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error
	GetAdDetailForSlot(ctx context.Context, id modelad.ID) (modelfullad.DetailID, error)
	GetAdSlot(ctx context.Context, min_cost uint32) (modelad.Ads, error)
	GetAdCount(ctx context.Context, clientID modeluser.ID) (int64, error)
}

type UseCase struct {
	adRepo adRepositoryI
}

func New(adRepo adRepositoryI) *UseCase {
	return &UseCase{
		adRepo: adRepo,
	}
}

func (u *UseCase) FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelfullad.AdFullInfo, error) {
	ads, err := u.adRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return ads, nil
}

func (u *UseCase) Create(ctx context.Context, ad modelad.Ads) (error) {
	if ad.Budget < 100{
		ad.Status = "non-active"
	} else {
		ad.Status = "active"
	}
	return u.adRepo.Create(ctx, ad)
}

func (u *UseCase) Update(ctx context.Context, ad modelad.Ads) error {
	return u.adRepo.Update(ctx, ad)
}

func (u *UseCase) Delete(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error {
	return u.adRepo.Delete(ctx, adID, clientID)
}

func (u *UseCase) GetOneAd(ctx context.Context, adID modelad.ID, clientID modeluser.ID) (modelfullad.AdFullInfo, int, error) {
	adInfo, err := u.adRepo.GetOneAd(ctx, adID, clientID)
	conversion := -1
	if err != nil {
		return modelfullad.AdFullInfo{}, conversion, fmt.Errorf("Failed to get ad with id error %w", err)
	}
	if adInfo.Impressions != 0{
		conversion = adInfo.Clicks / adInfo.Impressions
	}
	if adInfo.Budget < 100{
		adInfo.Status = "non-active"
	}

	return adInfo, conversion, nil
}

func (u *UseCase) GetAdDetailForSlot(ctx context.Context, id modelad.ID) (modelfullad.DetailID, error){
	return u.adRepo.GetAdDetailForSlot(ctx, id)
}

func (u *UseCase) GetAdSlot(ctx context.Context, min_cost uint32) (modelad.Ads, error){
	return u.adRepo.GetAdSlot(ctx, min_cost)
}

func (u *UseCase) GetAdCount(ctx context.Context, clientID modeluser.ID) (int64, error) {
	return u.adRepo.GetAdCount(ctx, clientID)
}