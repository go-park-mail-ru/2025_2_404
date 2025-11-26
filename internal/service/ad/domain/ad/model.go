package ad

import (
	modeluser "2025_2_404/internal/service/ad/domain/user"

	"github.com/google/uuid"
)

type ID = uuid.UUID

type Ads struct {
	ID       ID    `json:"add_id"`
	ClientID modeluser.ID
	Title string `json:"title"`
	Content    string `json:"content"`
	ImagePath     string `json:"img_path"`
	TargetUrl	string	`json:"target_url"`
}