package infra

import (
	"database/sql"
	modeluser "2025_2_404/internal/models/user"
	modelad "2025_2_404/internal/models/ad"
	"context"
	"fmt"
)

const(
	sqlTextForInsertSession = "INSERT INTO session (user_id, session_id) VALUES ($1, $2)"
	sqlTextForFoundUser = "SELECT user_id FROM session WHERE session_id = $1"
	sqlTextForSelectUsers = "SELECT id, password FROM app_user WHERE email = $1 "
	sqlTextForSelectAds = "SELECT id, file_path, title, text_ad FROM ad WHERE creator_id = $1"
	sqlTextForInsertUsers = "INSERT INTO app_user (email, password, user_name) VALUES ( $1, $2, $3) RETURNING id"
	sqlTextForInsertAds = "INSERT INTO ad (creator_id, file_path, title, text_ad) VALUES ($1, $2, $3, $4)"
)

type sqlI interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type DB struct {
	SQL sqlI
}

func (r *DB) CreateUser(ctx context.Context, user modeluser.RegisterUser) (int, error) {
	var userID int
	err := r.SQL.QueryRowContext(ctx, sqlTextForInsertUsers, user.Email, user.Password, user.UserName).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return userID, nil
}

func (r *DB) FindUserByEmail(ctx context.Context, email string) (modeluser.RegisterUser, error) {
	var user modeluser.RegisterUser
	err := r.SQL.QueryRowContext(ctx, sqlTextForSelectUsers, email).Scan(&user.ID, &user.Password)
	if err != nil {
		return user, fmt.Errorf("failed to find user by email: %w", err)
	}
	return user, nil
}

func (r *DB) CreateSession(ctx context.Context, userID, sessionID string) (string, error) {
	_, err := r.SQL.ExecContext(ctx, sqlTextForInsertSession, userID, sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	return sessionID, nil
}

func (r *DB) FindSession(ctx context.Context, sessionID int) (string, error) {
	var userID string
	err := r.SQL.QueryRowContext(ctx, sqlTextForFoundUser, sessionID).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("failed to find session by user ID: %w", err)
	}
	return userID, nil
}

func (r *DB) CreateAd(ctx context.Context, ad modelad.Ads) (int, error) {
	var adID int
	err := r.SQL.QueryRowContext(ctx, sqlTextForInsertAds, ad.CreatorID, ad.FilePath, ad.Title, ad.Text).Scan(&adID)
	if err != nil {
		return 0, fmt.Errorf("failed to create ad: %w", err)
	}
	return adID, nil
}

func (r *DB) FindAdByID(ctx context.Context, adID int) (modelad.Ads, error) {
	var ad modelad.Ads
	err := r.SQL.QueryRowContext(ctx, sqlTextForSelectAds, adID).Scan(&ad.ID, &ad.FilePath, &ad.Title, &ad.Text)
	if err != nil {
		return modelad.Ads{}, fmt.Errorf("failed to find ad by ID: %w", err)
	}
	return ad, nil
}
