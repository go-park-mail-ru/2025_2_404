package usecase

import (
	"2025_2_404/internal/use_case/ad"
	"2025_2_404/internal/use_case/auth"
	"2025_2_404/internal/use_case/token"
	"2025_2_404/internal/use_case/profile"
	"2025_2_404/internal/config"
	"2025_2_404/internal/repository/postgres"
	"2025_2_404/internal/use_case/filestorage"
	"2025_2_404/internal/use_case/balance"
	"2025_2_404/internal/use_case/feed"
)

type Config struct	{
	AdUsecase *ad.UseCase
	AuthUsecase *auth.UseCase
	ProfileUsecase *profile.UseCase
	TokenUsecase *token.UseCase
	StorageUsecase	*filestorage.UseCase
	BalanceUsecase 	*balance.UseCase
	FeedUsecase *feed.UseCase
}

func New(cfg *config.Config , configRepo *postgres.Config) *Config {
	tokenUsecase := token.New(cfg)
	storageUsecase := filestorage.New(cfg)
	authUsecase := auth.New(configRepo.AuthRepo, tokenUsecase)
	adUsecase := ad.New(configRepo.AdRepo)
	profileUsecase := profile.New(configRepo.ProfileRepo)
	balanceUsecase := balance.New(configRepo.BalanceRepo)
	feedUsecase := feed.New(configRepo.FeedRepo)
	return &Config {
		AdUsecase:	adUsecase,
		AuthUsecase:	authUsecase,
		ProfileUsecase: profileUsecase,
		TokenUsecase:	tokenUsecase,
		BalanceUsecase: balanceUsecase,
		FeedUsecase: feedUsecase,
		StorageUsecase: storageUsecase,
	}
}