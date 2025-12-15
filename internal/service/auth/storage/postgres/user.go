package postgres

import (
	modeluser "2025_2_404/internal/service/auth/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const(
	sqlTextForSelectUsers = "SELECT id, password_hash FROM client WHERE email = $1"
	// sqlTextForInsertBalance = "INSERT INTO client_wallet (client_id, balance) VALUES ($1, $2)"
	sqlTextForInsertUsers = "INSERT INTO client (email, password_hash, name) VALUES ($1, $2, $3) RETURNING id"
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
	// if err == sql.ErrNoRows {
	// 	fmt.Println("Пользователь с таким ID не существует")
	// 	return uuid.Nil, fmt.Errorf("Пользователь с таким ID не существует")
	// }
	if err != nil {
        var pgErr *pq.Error
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return uuid.Nil, errors.New("already exists")
        }
        return uuid.Nil, err
    }
    return user.ID, nil
	// _, err = r.sql.ExecContext(ctx, sqlTextForInsertBalance, user.ID, 0)
	// if err != nil {
	// 	log.Println("Не удалось создать баланс пользователя, причина: %w", err)
	// 	return user.ID, fmt.Errorf("balance not added: %w", err)
	// }
}

func (r *DB) FindByEmail(ctx context.Context, email string) (modeluser.User, error) {
	var user modeluser.User
	err := r.sql.QueryRowContext(ctx, sqlTextForSelectUsers, email).Scan(&user.ID, &user.HashedPassword)
	if err != nil {
		log.Println("Не удалось найти пользователя по email, причина: %w", err)
		return user, fmt.Errorf("failed to find user by email: %w", err)
	}
	return user, nil
}	