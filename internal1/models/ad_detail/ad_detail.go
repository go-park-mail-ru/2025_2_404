package addetail

import (
	modelad "2025_2_404/internal/domain/models/ad"
	modelPlatform "2025_2_404/internal/domain/models/platform"
)

type ID int

type AdDetail struct {
	ID       ID    `json:"ad_detail_id"`
	AdID modelad.ID
	PlatformID modelPlatform.ID
	AmountForAd int `json:"amount_for_ad"`
}