// Package ad provides use case implementations for advertisement management.
package ad

import (
	modelad "2025_2_404/internal/service/ad/domain/ad"
	modelfullad "2025_2_404/internal/service/ad/domain/ad_full_info"
	modeluser "2025_2_404/internal/service/ad/domain/user"
	"context"

	"go.uber.org/zap"
)

type adRepositoryI interface {
	FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelfullad.AdFullInfo, error)
	Create(ctx context.Context, ad modelad.Ads) error
	GetOneAd(ctx context.Context, adID modelad.ID, clientID modeluser.ID) (modelfullad.AdFullInfo, error)
	Update(ctx context.Context, ad modelad.Ads) error
	Delete(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error
	GetAdDetailForSlot(ctx context.Context, id modelad.ID, click, impression int) (modelfullad.DetailID, error)
	GetAdSlot(ctx context.Context, minCost uint32) (modelad.Ads, error)
	GetAdCount(ctx context.Context, clientID modeluser.ID) (int64, error)
}

type UseCase struct {
	adRepo adRepositoryI
	logger *zap.Logger
}

func New(adRepo adRepositoryI, logger *zap.Logger) *UseCase {
	return &UseCase{
		adRepo: adRepo,
		logger: logger,
	}
}

func (u *UseCase) FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelfullad.AdFullInfo, error) {
	return u.adRepo.FindByUserID(ctx, userID)
}

func (u *UseCase) Create(ctx context.Context, ad modelad.Ads) error {
	// Если бюджет меньше 100 делаем неактивной
	if ad.Budget < 100 {
		ad.Status = "non-active"
	} else {
		ad.Status = "active"
	}

	u.logger.Info("creating ad",
		zap.String("ad_id", ad.ID.String()),
		zap.String("client_id", ad.ClientID.String()),
		zap.String("title", ad.Title),
		zap.Uint32("budget", ad.Budget),
		zap.String("status", ad.Status),
	)

	if err := u.adRepo.Create(ctx, ad); err != nil {
		u.logger.Error("failed to create ad in repo",
			zap.String("client_id", ad.ClientID.String()),
			zap.Error(err),
		)
		return err
	}

	u.logger.Info("ad created successfully",
		zap.String("ad_id", ad.ID.String()),
		zap.String("client_id", ad.ClientID.String()),
	)

	return nil
}

func (u *UseCase) Update(ctx context.Context, ad modelad.Ads) error {
	u.logger.Info("updating ad",
		zap.String("ad_id", ad.ID.String()),
		zap.String("client_id", ad.ClientID.String()),
	)

	if err := u.adRepo.Update(ctx, ad); err != nil {
		u.logger.Error("failed to update ad", zap.Error(err))
		return err
	}

	u.logger.Info("ad updated successfully",
		zap.String("ad_id", ad.ID.String()),
		zap.String("client_id", ad.ClientID.String()),
	)
	return nil
}

func (u *UseCase) Delete(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error {
	u.logger.Info("deleting ad",
		zap.String("ad_id", adID.String()),
		zap.String("client_id", clientID.String()),
	)

	if err := u.adRepo.Delete(ctx, adID, clientID); err != nil {
		u.logger.Error("failed to delete ad", zap.Error(err))
		return err
	}

	u.logger.Info("ad deleted successfully",
		zap.String("ad_id", adID.String()),
		zap.String("client_id", clientID.String()),
	)
	return nil
}

func (u *UseCase) GetOneAd(ctx context.Context, adID modelad.ID, clientID modeluser.ID) (modelfullad.AdFullInfo, int, error) {
	adInfo, err := u.adRepo.GetOneAd(ctx, adID, clientID)
	if err != nil {
		// Репозиторий уже возвращает globalerrors — просто прокидываем
		return modelfullad.AdFullInfo{}, -1, err
	}

	// Вычисляем конверсию (целочисленное деление → может быть 0!)
	conversion := -1
	if adInfo.Impressions > 0 {
		conversion = adInfo.Clicks * 100 / adInfo.Impressions // ← проценты? или оставить как есть?
		// Если нужен float — возвращать float64, но в твоём случае int — ок.
	}

	// Обновляем статус на лету, если бюджет < 100
	if adInfo.Budget < 100 {
		adInfo.Status = "non-active"
	}

	return adInfo, conversion, nil
}

func (u *UseCase) GetAdDetailForSlot(ctx context.Context, id modelad.ID, eventType string) (modelfullad.DetailID, error) {
	click := 0
	impression := 0
	if eventType == "impression" {
		impression = 1
	} else if eventType == "click" {
		click = 1
	}

	return u.adRepo.GetAdDetailForSlot(ctx, id, click, impression)
}

func (u *UseCase) GetAdSlot(ctx context.Context, minCost uint32) (modelad.Ads, error) {
	return u.adRepo.GetAdSlot(ctx, minCost)
}

func (u *UseCase) GetAdCount(ctx context.Context, clientID modeluser.ID) (int64, error) {
	return u.adRepo.GetAdCount(ctx, clientID)
}
