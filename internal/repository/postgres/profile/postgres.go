package profile

import (
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
	"database/sql"
	"fmt"
)

const(
	sqlTextForUpdateClient = "UPDATE client SET name = $1, email = $2, img_path = $3 WHERE id = $4"
	sqlTextForShowClient = "SELECT name, email FROM client WHERE id = $1"
)

type DB struct{
	sql *sql.DB
}

func New(sql *sql.DB) *DB{
	return &DB{
		sql: sql,
	}
}

func (r *DB) Update(ctx context.Context, client modeluser.User) error {
	res, err := r.sql.ExecContext(ctx, sqlTextForUpdateClient, client.UserName, client.Email, client.ImagePath, client.ID)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ad client id %v not found", client.ID)
	}

	return nil
}

func (r *DB) Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, error){
	var client modeluser.User
	err := r.sql.QueryRowContext(ctx, sqlTextForShowClient, clientID).Scan(&client.UserName, &client.Email)
	if err != nil {
		return modeluser.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return client, nil
}
