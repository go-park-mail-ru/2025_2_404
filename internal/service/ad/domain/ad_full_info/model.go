package adfullinfo

import "github.com/google/uuid"

type ID = uuid.UUID
type DetailID = uuid.UUID

type AdFullInfo struct {
	ID        ID     `json:"add_id"`
	AddDetailID	DetailID `json:"add_detail_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	ImgPath    string `json:"img_bin"`
	TargetUrl string `json:"target_url"`

	Budget uint32 `json:"budget"`

	Clicks      int `json:"clicks"`
	Impressions int `json:"impressions"`
}