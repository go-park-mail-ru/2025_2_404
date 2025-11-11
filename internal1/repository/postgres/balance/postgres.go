package balance

import (
	modelwallet "2025_2_404/internal/domain/models/client_wallet"
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
	"database/sql"
	"fmt"
)

const(
	sqlTextForShowBalance = "SELECT balance FROM client_wallet WHERE client_id = $1"
)

type DB struct{
	sql *sql.DB
}

func New(sql *sql.DB) *DB{
	return &DB{
		sql: sql,
	}
}

func (r *DB) Show(ctx context.Context, clientID modeluser.ID) (modelwallet.Balance, error){

	var balance modelwallet.Balance
	err := r.sql.QueryRowContext(ctx, sqlTextForShowBalance, clientID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("error with DB: %w",err)
	}

	return balance, nil
}