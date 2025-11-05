package postgres

import (
	"2025_2_404/internal/repository/postgres/ad"
	"2025_2_404/internal/repository/postgres/auth"
	"2025_2_404/internal/repository/postgres/profile"
	"2025_2_404/internal/repository/postgres/balance"
	"2025_2_404/internal/connections"
)

type Config struct {
	AdRepo		*ad.DB
	AuthRepo	*auth.DB
	ProfileRepo	*profile.DB
	BalanceRepo	*balance.DB
}

func New(connCFG *connections.Config) *Config {
	adRepo := ad.New(connCFG.PostgresSQL)
	authRepo := auth.New(connCFG.PostgresSQL)
	profileRepo := profile.New(connCFG.PostgresSQL)
	balanceRepo := balance.New(connCFG.PostgresSQL)
	return &Config{
		AdRepo:	adRepo,
		AuthRepo:	authRepo,
		ProfileRepo: profileRepo,
		BalanceRepo: balanceRepo,
	}
}

