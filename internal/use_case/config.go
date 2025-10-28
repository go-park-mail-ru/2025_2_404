package usecase

import (
	"2025_2_404/internal/use_case/ad"
	"2025_2_404/internal/use_case/auth"
	"2025_2_404/internal/use_case/token"
	"2025_2_404/internal/config"
	"2025_2_404/internal/repository/postgres"
)

type Config struct	{
	AdUsecase *ad.UseCase
	AuthUsecase *auth.UseCase
	TokenUsecase *token.UseCase
}

func New(cfg *config.Config , configRepo *postgres.Config) *Config {
	tokenUsecase := token.New(cfg)
	authUsecase := auth.New(configRepo.AuthRepo, tokenUsecase)
	adUsecase := ad.New(configRepo.AdRepo)
	return &Config {
		AdUsecase:	adUsecase,
		AuthUsecase:	authUsecase,
		TokenUsecase:	tokenUsecase,
	}
}