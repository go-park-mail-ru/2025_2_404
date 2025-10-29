package ad

import (
	"database/sql"
	"context"
	"fmt"
	modelad "2025_2_404/internal/domain/models/ad"
	modeluser "2025_2_404/internal/domain/models/user"
)

const(
	sqlTextForSelectAds = "SELECT id, title, content, img_bin, target_url FROM ad WHERE client_id = $1"
	sqlTextForInsertAds = "INSERT INTO ad (client_id, title, content, img_bin, target_url) VALUES ($1, $2, $3, $4, $5) RETURNING id"
)

type DB struct {
	sql *sql.DB
}

func New(sql *sql.DB) *DB {
	return &DB{
		sql: sql,
	}
}

func (r *DB) FindByUserID(ctx context.Context, userID modeluser.ID) (modelad.Ads, error) {
	var ad modelad.Ads
	err := r.sql.QueryRowContext(ctx, sqlTextForSelectAds, userID).Scan(&ad.ID, &ad.Title, &ad.Content, &ad.ImgBin, &ad.TargetUrl)
	if err != nil {
		return modelad.Ads{}, fmt.Errorf("failed to find ad by user ID: %w", err)
	}
	return ad, nil
}

func (r *DB) Create(ctx context.Context, ad modelad.Ads) (int, error) {
	var adID int
	err := r.sql.QueryRowContext(ctx, sqlTextForInsertAds, ad.ClientID, ad.Title, ad.Content, ad.ImgBin, ad.TargetUrl).Scan(&adID)
	if err != nil {
		return 0, fmt.Errorf("failed to create ad: %w", err)
	}
	return adID, nil
}
