package db

import (
	"database/sql"
	_ "github.com/jackc/pgx"
)

func ConnectDB() (*sql.DB, error) {
	connectString := "postgresql://user:password@localhost:5432/database"
	db, err := sql.Open("pgx", connectString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CloseDB(db *sql.DB) error {
	return db.Close()
}

func PingDB(db *sql.DB) error {
	return db.Ping()
}


