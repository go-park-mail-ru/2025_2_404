package auth

import (
	"database/sql"
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
	"fmt"
)

const(
	sqlTextForInsertSession = "INSERT INTO session (user_id, session_id) VALUES ($1, $2)"
	sqlTextForFoundUser = "SELECT user_id FROM session WHERE session_id = $1"
	sqlTextForSelectUsers = "SELECT id, password FROM app_user WHERE email = $1 "
	sqlTextForInsertUsers = "INSERT INTO app_user (email, password, user_name) VALUES ( $1, $2, $3) RETURNING id"
	sqlTextForFoundSession = "SELECT session_id FROM session WHERE user_id = $1"
)

type sqlI interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type DB struct {
	sql sqlI
}

func (r *DB) CreateUser(ctx context.Context, user *modeluser.User) (modeluser.ID, error) {
	err := r.sql.QueryRowContext(ctx, sqlTextForInsertUsers, user.Email, user.HashedPassword, user.UserName).Scan(&user.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return user.ID, nil
}

func (r *DB) FindUserByEmail(ctx context.Context, email string) (modeluser.User, error) {
	var user modeluser.User
	err := r.sql.QueryRowContext(ctx, sqlTextForSelectUsers, email).Scan(&user.ID, &user.HashedPassword)
	if err != nil {
		return user, fmt.Errorf("failed to find user by email: %w", err)
	}
	return user, nil
}

func (r *DB) CreateSession(ctx context.Context, userID modeluser.ID, sessionID string) (string, error) {
	_, err := r.sql.ExecContext(ctx, sqlTextForInsertSession, userID, sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	return sessionID, nil
}

func (r *DB) FindUserBySessionID(ctx context.Context, sessionID string) (string, error) {
	var userID string
	err := r.sql.QueryRowContext(ctx, sqlTextForFoundUser, sessionID).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("failed to find session by user ID: %w", err)
	}
	return userID, nil
}

func (r *DB) FindSessionByUserID(ctx context.Context, userID modeluser.ID) (string, error) {
	var sessionID string
	err := r.sql.QueryRowContext(ctx, sqlTextForFoundSession, userID).Scan(&sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to find session by user ID: %w", err)
	}
	return sessionID, nil
}