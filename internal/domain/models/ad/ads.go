package ad

import	(
	modeluser "2025_2_404/internal/domain/models/user"
)

type ID int

type Ads struct {
	ID       ID    `json:"add_id"`
	CreatorID modeluser.ID    `json:"creater_id"`
	FilePath string `json:"file_path"`
	Title    string `json:"title"`
	Text     string `json:"text"`
}
