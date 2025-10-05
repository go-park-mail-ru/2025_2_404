package ad

type ID int64

type Ads struct {
	ID       ID    `json:"add_id"`
	CreatorID int    `json:"creater_id"`
	FilePath string `json:"file_path"`
	Title    string `json:"title"`
	Text     string `json:"text"`
}
