// Package postgres provides a PostgreSQL implementation of the metric repository.
package postgres

import (
	user "2025_2_404/internal/service/profile/domain"
	"2025_2_404/internal/service/slot/domain/metric"
	"2025_2_404/pkg/globalerrors"
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"
)

const (
	sqlTextForCreateMetric = `
		INSERT INTO slot_event (slot_id, ad_detail_id, event_type)
		VALUES ($1, $2, $3)
	`

	sqlTextForGetClientID = `
		SELECT user_id 
		FROM slots 
		WHERE id = $1
	`

	sqlTextForGetMetric = `
		SELECT
			DATE(created_at) AS event_date,
			COUNT(*) FILTER (WHERE event_type = 'impression') AS impressions,
			COUNT(*) FILTER (WHERE event_type = 'click') AS clicks
		FROM
			slot_event
		WHERE
			slot_id = $1
		GROUP BY
			DATE(created_at)
		ORDER BY
			event_date
	`
)

type DB struct {
	sql *sql.DB
}

func New(sql *sql.DB) *DB {
	return &DB{sql: sql}
}

func (r *DB) CreateMetric(ctx context.Context, m metric.Metric) (user.ID, error) {
	_, err := r.sql.ExecContext(ctx, sqlTextForCreateMetric, m.SlotID, m.AdDetailID, m.EventType)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503": // foreign_key_violation (slot_id или ad_detail_id не существуют)
				return user.ID{}, globalerrors.ErrSlotNotFound
			case "23514": // check_violation (event_type не 'click'/'impression')
				return user.ID{}, globalerrors.ErrInvalidQuery
			}
		}
		return user.ID{}, globalerrors.ErrInternal
	}

	var userID user.ID
	err = r.sql.QueryRowContext(ctx, sqlTextForGetClientID, m.SlotID).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.ID{}, globalerrors.ErrSlotNotFound
		}
		return user.ID{}, globalerrors.ErrInternal
	}

	return userID, nil
}

func (r *DB) GetMetricForDay(ctx context.Context, slotID metric.SlotID) ([]metric.GetMetric, error) {
	rows, err := r.sql.QueryContext(ctx, sqlTextForGetMetric, slotID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Например, неверный UUID
			if pgErr.Code == "22P02" { // invalid_text_representation
				return nil, globalerrors.ErrInvalidQuery
			}
		}
		return nil, globalerrors.ErrInternal
	}
	defer func() {
		_ = rows.Close()
	}()

	var metrics []metric.GetMetric
	for rows.Next() {
		var m metric.GetMetric
		if err := rows.Scan(&m.EventDate, &m.Impressions, &m.Clicks); err != nil {
			return nil, globalerrors.ErrInternal
		}
		metrics = append(metrics, m)
	}

	if err = rows.Err(); err != nil {
		return nil, globalerrors.ErrInternal
	}

	return metrics, nil
}
