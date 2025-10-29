package ad

import	(
	modeluser "2025_2_404/internal/domain/models/user"
)

type ID int

type Ads struct {
	ID       ID    `json:"add_id"`
	ClientID modeluser.ID
	Title string `json:"title"`
	Content    string `json:"content"`
	ImgBin     string `json:"img_bin"`
	TargetUrl	string	`json:"target_url"`
}
