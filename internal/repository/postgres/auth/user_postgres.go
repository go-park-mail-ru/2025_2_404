package auth

import (
	"database/sql"
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
	"fmt"
)

const(
	sqlTextForSelectUsers = "SELECT id, password_hash FROM client WHERE email = $1 "
	sqlTextForInsertUsers = "INSERT INTO client (email, password_hash, name) VALUES ( $1, $2, $3) RETURNING id"
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
	err := r.sql.QueryRowContext(ctx, sqlTextForInsertUsers, user.Email, user.HashedPassword, user.UserName).Scan(&user.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
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