package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"2025_2_404/config"
)

func ConnectDB(config *config.PostgresConfig) (*sql.DB, error) {
	connectString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.User, config.Password, config.Host, config.Port, config.DB)
	db, err := sql.Open("pgx", connectString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func CloseDB(db *sql.DB) error {
	return db.Close()
}

