package infra

import (
	"database/sql"
	modeluser "2025_2_404/internal/models"
	"context"
	"fmt"
	modelad "2025_2_404/internal/models"
)

const(
	sqlTextForSelectAds = "SELECT id, file_path, title, text_ad FROM ad WHERE creator_id = $1"
	sqlTextForInsertAds = "INSERT INTO ad (creator_id, file_path, title, text_ad) VALUES ($1, $2, $3, $4)"
)

type sqlI interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type DB struct {
	sql sqlI
}

func New(sql sqlI) *DB {
	return &DB{
		sql: sql,
	}
}

func (r *DB) FindAdByUserID(ctx context.Context, userID modeluser.ID) (modelad.Ads, error) {
	var ad modelad.Ads
	err := r.sql.QueryRowContext(ctx, sqlTextForSelectAds, userID).Scan(&ad.ID, &ad.FilePath, &ad.Title, &ad.Text)
	if err != nil {
		return modelad.Ads{}, fmt.Errorf("failed to find ad by user ID: %w", err)
	}
	return ad, nil
}

func (r *DB) CreateAd(ctx context.Context, ad modelad.Ads) (int, error) {
	var adID int
	err := r.sql.QueryRowContext(ctx, sqlTextForInsertAds, ad.CreatorID, ad.FilePath, ad.Title, ad.Text).Scan(&adID)
	if err != nil {
		return 0, fmt.Errorf("failed to create ad: %w", err)
	}
	return adID, nil
}

const(
	sqlTextForInsertSession = "INSERT INTO session (user_id, session_id) VALUES ($1, $2)"
	sqlTextForFoundUser = "SELECT user_id FROM session WHERE session_id = $1"
	sqlTextForSelectUsers = "SELECT id, password FROM app_user WHERE email = $1 "
	sqlTextForInsertUsers = "INSERT INTO app_user (email, password, user_name) VALUES ( $1, $2, $3) RETURNING id"
)


func (r *DB) CreateUser(ctx context.Context, user modeluser.User) (int, error) {
	var userID int
	err := r.sql.QueryRowContext(ctx, sqlTextForInsertUsers, user.GetEmail(), user.GetHashedPassword(), user.GetUserName()).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return userID, nil
}

func (r *DB) FindUserByEmail(ctx context.Context, email string) (modeluser.User, error) {
	var user modeluser.User
	err := r.sql.QueryRowContext(ctx, sqlTextForSelectUsers, email).Scan(&user.ID, &user.Email)
	if err != nil {
		return user, fmt.Errorf("failed to find user by email: %w", err)
	}
	return user, nil
}

func (r *DB) CreateSession(ctx context.Context, userID, sessionID string) (string, error) {
	_, err := r.sql.ExecContext(ctx, sqlTextForInsertSession, userID, sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	return sessionID, nil
}

func (r *DB) FindSession(ctx context.Context, sessionID int) (string, error) {
	var userID string
	err := r.sql.QueryRowContext(ctx, sqlTextForFoundUser, sessionID).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("failed to find session by user ID: %w", err)
	}
	return userID, nil
}



