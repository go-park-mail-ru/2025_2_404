package adfullinfo

type ID int

type AdFullInfo struct {
	ID        ID     `json:"add_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	ImgBin    string `json:"img_bin"`
	TargetUrl string `json:"target_url"`

	AmountForAd int `json:"amount_for_ad"`

	Clicks      int `json:"clicks"`
	Impressions int `json:"impressions"`
}