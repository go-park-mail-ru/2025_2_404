package ad

import (
	"database/sql"
	"context"
	"fmt"
	modelad "2025_2_404/internal/domain/models/ad"
	modeluser "2025_2_404/internal/domain/models/user"
)

const(
	sqlTextForSelectAds = "SELECT id, file_path, title, text_ad FROM ad WHERE creator_id = $1"
	sqlTextForInsertAds = "INSERT INTO ad (creator_id, file_path, title, text_ad) VALUES ($1, $2, $3, $4)"
)

type DB struct {
	sql *sql.DB
}

func New(sql *sql.DB) *DB {
	return &DB{
		sql: sql,
	}
}

func (r *DB) FindAdByUserID(ctx context.Context, userID modeluser.ID) (modelad.Ads, error) {
	var ad modelad.Ads
	err := r.sql.QueryRowContext(ctx, sqlTextForSelectAds, userID).Scan(&ad.ID, &ad.FilePath, &ad.Title, &ad.Text)
	if err != nil {
		return modelad.Ads{}, fmt.Errorf("failed to find ad by user ID: %w", err)
	}
	return ad, nil
}

func (r *DB) CreateAd(ctx context.Context, ad modelad.Ads) (int, error) {
	var adID int
	err := r.sql.QueryRowContext(ctx, sqlTextForInsertAds, ad.CreatorID, ad.FilePath, ad.Title, ad.Text).Scan(&adID)
	if err != nil {
		return 0, fmt.Errorf("failed to create ad: %w", err)
	}
	return adID, nil
}
