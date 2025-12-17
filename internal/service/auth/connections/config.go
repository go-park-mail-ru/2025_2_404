// Package connections manages database connections for the authentication service.
package connections

import (
	"2025_2_404/internal/service/auth/config"
	"2025_2_404/internal/service/auth/connections/postgres"
	"database/sql"
)

type Config struct {
	PostgresSQL *sql.DB
}

func New(cfg *config.Config) (*Config, error) {
	postgresSQL, err := postgres.ConnectDB(cfg.DBConfig)
	if err != nil {
		return nil, err
	}
	return &Config{
		PostgresSQL: postgresSQL,
	}, nil
}

func (c *Config) CloseAll() {
	_ = c.PostgresSQL.Close()
}
