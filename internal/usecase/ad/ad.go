package ad

import (
	// modelad "2025_2_404/internal/models"
	// modeluser "2025_2_404/internal/models"
	"2025_2_404/internal/models"
	"context"
)

type adRepositoryI interface {
	FindAdByUserID(ctx context.Context, userID models.User) (models.Ads, error)
	CreateAd(ctx context.Context, ad models.Ads) (int, error)
	UpdateAd(ctx context.Context, ad models.Ads) (error)
}

type AdUseCase struct {
	adRepo adRepositoryI
}

func New(adRepo adRepositoryI) *AdUseCase {
	return &AdUseCase{
		adRepo: adRepo,
	}
}

func (u *AdUseCase) FindAdByUserID(ctx context.Context, userID int64) (models.Ads, error) {
	return u.adRepo.FindAdByUserID(ctx, userID)
}

func (u *AdUseCase) CreateAd(ctx context.Context, ad models.Ads) (int, error) {
	return u.adRepo.CreateAd(ctx, ad)
}

<<<<<<< HEAD
func (u *AdUseCase) UpdateAd(ctx context.Context, ad models.Ads) error {
	_, err := u.adRepo.FindAdByUserID(ctx, ad.CreatorID)
	if err != nil {
        // Хорошей практикой будет проверить, найдено ли объявление вообще
		return err 
	}
	return ctx.Err()
}
=======

>>>>>>> a8230ea6cc45a4ef7d6d317222973fdc7959bd18
