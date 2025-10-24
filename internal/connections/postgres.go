package db

import (
	"2025_2_404/internal/config"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func ConnectDB(config *config.PostgresConfig) (*sql.DB, error) {
	connectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DB)
	db, err := sql.Open("pgx", connectString)
	if err != nil {
		return nil, err
	}
	i := 0
	for err := db.Ping(); err != nil; err = db.Ping() {
		i++
		if i >= 10 {
			return nil, err
		}
		time.Sleep(1 * time.Second)
	}
	return db, nil
}

func CloseDB(db *sql.DB) error {
	return db.Close()
}

