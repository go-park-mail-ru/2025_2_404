// Package postgres provides PostgreSQL storage implementation for user authentication.
package postgres

import (
	modeluser "2025_2_404/internal/service/auth/domain"
	"2025_2_404/pkg/globalerrors"
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgconn"
	"go.uber.org/zap"
)

const (
	sqlTextForSelectUsers = "SELECT id, password_hash FROM client WHERE email = $1"
	// sqlTextForInsertBalance = "INSERT INTO client_wallet (client_id, balance) VALUES ($1, $2)"
	sqlTextForInsertUsers = "INSERT INTO client (email, password_hash, name) VALUES ($1, $2, $3) RETURNING id"
)

type DB struct {
	sql    *sql.DB
	logger *zap.Logger
}

func New(sql *sql.DB, logger *zap.Logger) *DB {
	return &DB{
		sql:    sql,
		logger: logger,
	}
}

func (r *DB) Create(ctx context.Context, user *modeluser.User) (modeluser.ID, error) {
	err := r.sql.QueryRowContext(ctx, sqlTextForInsertUsers, user.Email, user.HashedPassword, user.UserName).Scan(&user.ID)
	// if err == sql.ErrNoRows {
	// 	fmt.Println("Пользователь с таким ID не существует")
	// 	return uuid.Nil, fmt.Errorf("Пользователь с таким ID не существует")
	// }
	if err == nil {
		return user.ID, nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return modeluser.ID{}, globalerrors.ErrUserAlreadyExists

		case "23514":
			if strings.Contains(strings.ToLower(pgErr.Message), "email") ||
				strings.Contains(strings.ToLower(pgErr.ConstraintName), "email") {
				return modeluser.ID{}, globalerrors.ErrNonValidEmail
			}
			return modeluser.ID{}, globalerrors.ErrInvalidQuery

		case "23502":
			return modeluser.ID{}, globalerrors.ErrInvalidQuery

		default:
			return modeluser.ID{}, globalerrors.ErrInternal
		}
	}
	// _, err = r.sql.ExecContext(ctx, sqlTextForInsertBalance, user.ID, 0)
	// if err != nil {
	// 	log.Println("Не удалось создать баланс пользователя, причина: %w", err)
	// 	return user.ID, fmt.Errorf("balance not added: %w", err)user already exists
	// }
	return modeluser.ID{}, globalerrors.ErrorContextTimeout
}

func (r *DB) FindByEmail(ctx context.Context, email string) (modeluser.User, error) {
	var user modeluser.User
	err := r.sql.QueryRowContext(ctx, sqlTextForSelectUsers, email).Scan(&user.ID, &user.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return modeluser.User{}, globalerrors.ErrWrongEmailOrPassword
		}
		return modeluser.User{}, globalerrors.ErrInternal
	}
	return user, nil
}
