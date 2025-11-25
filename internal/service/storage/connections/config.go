package connections

import (
	"2025_2_404/internal/service/storage/config"
	"2025_2_404/internal/service/storage/connections/postgres"
	"database/sql"
)

type Config struct {
	PostgresSQL	*sql.DB
}

func New(cfg *config.Config) (*Config, error) {
	postgresSQL, err := postgres.ConnectDB(cfg.DBConfig)
	if err != nil {
		return nil, err
	}
	return  &Config{
		PostgresSQL: postgresSQL,
	}, nil
}

func (c *Config) CloseAll() {
	c.PostgresSQL.Close()
}