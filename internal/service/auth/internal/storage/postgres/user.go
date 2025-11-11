package postgres

import (
	modeluser "2025_2_404/internal/service/auth/internal/domain"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

const(
	sqlTextForSelectUsers = "SELECT id, password_hash FROM client WHERE email = $1"
	sqlTextForInsertBalance = "INSERT INTO client_wallet (client_id, balance) VALUES ($1, $2)"
	sqlTextForInsertUsers = "INSERT INTO client (id, email, password_hash, name) VALUES ($1, $2, $3, $4)"
)

type DB struct {
	sql *sql.DB
}

func New(sql *sql.DB) *DB {
	return &DB{
		sql: sql,
	}
}

func (r *DB) Create(ctx context.Context, user *modeluser.User) (modeluser.ID, error) {
	err := r.sql.QueryRowContext(ctx, sqlTextForInsertUsers, user.ID, user.Email, user.HashedPassword, user.UserName).Scan(&user.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}
	_, err = r.sql.ExecContext(ctx, sqlTextForInsertBalance, user.ID, 0)
	if err != nil {
		return user.ID, fmt.Errorf("balance not added: %w", err)
	}
	return user.ID, nil
}

func (r *DB) FindByEmail(ctx context.Context, email string) (modeluser.User, error) {
	var user modeluser.User
	err := r.sql.QueryRowContext(ctx, sqlTextForSelectUsers, email).Scan(&user.ID, &user.HashedPassword)
	if err != nil {
		return user, fmt.Errorf("failed to find user by email: %w", err)
	}
	return user, nil
}