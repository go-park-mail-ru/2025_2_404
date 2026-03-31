// Package adfullinfo defines the AdFullInfo domain model and related types.
package adfullinfo

import (
	"time"

	"github.com/google/uuid"
)

type ID = uuid.UUID
type DetailID = uuid.UUID

type AdFullInfo struct {
	ID          ID       `json:"add_id"`
	AddDetailID DetailID `json:"add_detail_id"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	ImgPath     string   `json:"img_bin"`
	TargetURL   string   `json:"target_url"`

	Budget uint32 `json:"budget"`
	Status string `json:"status"`

	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
	CreatedAt time.Time `json:"created_at"`

	Clicks      int `json:"clicks"`
	Impressions int `json:"impressions"`
}
