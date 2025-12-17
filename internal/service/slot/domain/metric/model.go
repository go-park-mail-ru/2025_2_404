// Package metric defines the Metric domain model and related types.
package metric

import (
	"time"

	"github.com/google/uuid"
)

type AdDetailID uuid.UUID
type SlotID uuid.UUID

type Metric struct {
	SlotID      SlotID
	AdDetailID  AdDetailID
	EventType   string
	CreateeTime time.Time
}

type GetMetric struct {
	SlotID      string    `db:"slot_id"`
	EventDate   time.Time `db:"event_date"`
	Impressions int       `db:"impressions"`
	Clicks      int       `db:"clicks"`
}

type MetricsForDay struct {
	SlotID      string `json:"slot_id"`
	Clicks      int32  `json:"clicks"`
	Impressions int32  `json:"impressions"`
	EventDate   string `json:"date"`
}

type GetMetricsResponse struct {
	SlotID           string          `json:"slot_id"`
	TotalClicks      int32           `json:"total_clicks"`
	TotalImpressions int32           `json:"total_impressions"`
	Metrics          []MetricsForDay `json:"daily_stats"`
}
