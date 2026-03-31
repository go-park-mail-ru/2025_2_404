// Package ad defines the Advertisement domain model and related types.
package ad

import (
	modeluser "2025_2_404/internal/service/ad/domain/user"
	"time"

	"github.com/google/uuid"
)

type ID = uuid.UUID

type Ads struct {
	ID        ID `json:"add_id"`
	ClientID  modeluser.ID
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Budget    uint32    `json:"budget"`
	ImagePath string    `json:"img_path"`
	TargetURL string    `json:"target_url"`
	Status    string    `json:"status"`
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
}
