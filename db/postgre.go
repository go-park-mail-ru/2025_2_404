package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx"
)

func ConnectDB(user, password, dbname, host, port string) (*sql.DB, error) {
	connectString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	db, err := sql.Open("pgx", connectString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CloseDB(db *sql.DB) error {
	return db.Close()
}



