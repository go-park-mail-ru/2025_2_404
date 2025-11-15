package support

import (
	modelsup "2025_2_404/internal/domain/models/support"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

const(
    sqlTextForInsertSupport = "INSERT INTO support (client_id, sup_status, category, sup_description, contact_name, contact_email) VALUES ($1, $2, $3, $4, $5) RETURNING id"
)

type DB struct{
	sql *sql.DB
}

func New(sql *sql.DB) *DB{
	return &DB{
		sql: sql,
	}
}

func (r *DB) Create(ctx context.Context, sup modelsup.Support) (error) {
	res, err := r.sql.ExecContext(ctx, sqlTextForInsertSupport, )
	if err != nil {
		return fmt.Errorf("failed to create ad: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ad with id %v not found", ad.ID)
	}

	return nil
}