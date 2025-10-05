package domain

type Ads struct {
	ID       string `json:"add_id"`
	CreatorID string `json:"creater_id"`
	FilePath string `json:"file_path"`
	Title    string `json:"title"`
	Text     string `json:"text"`
}
