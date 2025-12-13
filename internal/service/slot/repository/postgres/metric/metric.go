package postgres

import (
	"2025_2_404/internal/service/slot/domain/metric"
	"context"
	"database/sql"
	"fmt"
)

const(
	sqlTextForCreateMetric = `
		INSERT INTO slot_event (slot_id, ad_detail_id, event_type)
		VALUES ($1, $2, $3)
	`
)

type DB struct {
	sql *sql.DB
}

func New(sql *sql.DB) *DB {
	return &DB{sql: sql}
}

func (r *DB) CreateMetric(ctx context.Context, metric metric.Metric) error{
	_, err := r.sql.ExecContext(ctx, sqlTextForCreateMetric, metric.SlotID, metric.AdDetailID, metric.EventType)
	if err != nil{
		return fmt.Errorf("failed to insert metric: %w", err)
	}
	return nil
}