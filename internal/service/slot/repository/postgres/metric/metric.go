package postgres

import (
	user "2025_2_404/internal/service/profile/domain"
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

	sqlTextForGetClientID = `
		SELECT user_id 
		FROM slots 
		WHERE id = $1
	`

	sqlTextForGetMetric = `
		SELECT
		slot_id,
		DATE(created_at) AS event_date,
		COUNT(*) FILTER (WHERE event_type = 'impression') AS impressions,
		COUNT(*) FILTER (WHERE event_type = 'click') AS clicks
	FROM
		slot_event
	GROUP BY
		slot_id,
		DATE(created_at)
	ORDER BY
		slot_id,
		event_date;
	`
)

type DB struct {
	sql *sql.DB
}

func New(sql *sql.DB) *DB {
	return &DB{sql: sql}
}

func (r *DB) CreateMetric(ctx context.Context, metric metric.Metric) (user.ID ,error){
	_, err := r.sql.ExecContext(ctx, sqlTextForCreateMetric, metric.SlotID, metric.AdDetailID, metric.EventType)
	if err != nil{
		return user.ID{}, fmt.Errorf("failed to insert metric: %w", err)
	}
	var id user.ID
	err = r.sql.QueryRowContext(ctx, sqlTextForGetClientID, metric.SlotID).Scan(&id)
	if err != nil {
		return user.ID{}, fmt.Errorf("failed to select userID: %w", err)
	}
	return id, nil
}

func (r *DB) GetMetricForDay(ctx context.Context, slotID metric.SlotID)  ([]metric.GetMetric, error){
	
	var metrics []metric.GetMetric
	rows, err := r.sql.QueryContext(ctx, sqlTextForGetMetric)
	if err != nil{
		return nil, fmt.Errorf("failed to select metric: %w", err)
	}
	defer rows.Close()

	for rows.Next(){
		var metric metric.GetMetric
		if err := rows.Scan(&metric.SlotID, &metric.EventDate, &metric.Impressions, &metric.Clicks); err != nil {
			return nil, fmt.Errorf("failed to read rows: %w", err)
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}