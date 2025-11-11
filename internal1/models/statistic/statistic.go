package statistic

import	(
	modeldetailad "2025_2_404/internal/domain/models/ad_detail"
)

type ID int

type Statistic struct {
	ID       ID    `json:"statistic_id"`
	AdDetailID modeldetailad.ID
	Clicks int `json:"clicks"`
	Impressions    int `json:"impressions"`
}
