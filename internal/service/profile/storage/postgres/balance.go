// Package profile provides a PostgreSQL implementation of the profile repository.
package profile

import (
	modelclient "2025_2_404/internal/service/profile/domain"
	"context"
)

const (
	sqlTextForShowBalance   = "SELECT balance FROM client_wallet WHERE client_id = $1"
	sqlTextForUpdateBalance = "UPDATE client_wallet SET balance = $1 WHERE client_id = $2"
)

func (r *DB) ShowBalance(ctx context.Context, clientID modelclient.ID) (uint32, error) {
	var balance uint32

	err := r.sql.QueryRowContext(ctx, sqlTextForShowBalance, clientID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (r *DB) AddBalance(ctx context.Context, clientID modelclient.ID, addAmount uint32) error {
	balance, err := r.ShowBalance(ctx, clientID)
	if err != nil {
		return err
	}

	_, err = r.sql.ExecContext(ctx, sqlTextForUpdateBalance, balance+addAmount, clientID)
	if err != nil {
		return err
	}

	return nil
}

func (r *DB) SubtractBalance(ctx context.Context, clientID modelclient.ID, subAmount uint32) error {
	balance, err := r.ShowBalance(ctx, clientID)
	if err != nil {
		return err
	}

	_, err = r.sql.ExecContext(ctx, sqlTextForUpdateBalance, balance-subAmount, clientID)
	if err != nil {
		return err
	}

	return nil
}
