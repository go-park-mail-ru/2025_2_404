package feed

import (
	modelad "2025_2_404/internal/domain/models/ad"
	"context"
	"database/sql"
	"fmt"
)

const(
	sqlTextForSelectFeedAds = "SELECT a.id, a.title, a.content, a.target_url FROM ad a JOIN ad_detail ad_d ON a.id = ad_d.ad_id JOIN platform p ON ad_d.platform_id = p.id WHERE p.platform_name = $1;"
)

type DB struct{
	sql *sql.DB
}

func New(sql *sql.DB) *DB{
	return &DB{
		sql: sql,
	}
}

func (r *DB) Give(ctx context.Context, platformName string) ([]modelad.Ads, error) {
	rows, err := r.sql.QueryContext(ctx, sqlTextForSelectFeedAds, platformName)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer rows.Close()

	var ads []modelad.Ads
	for rows.Next() {
		var ad modelad.Ads
		err := rows.Scan(
			&ad.ID,
			&ad.Title,
			&ad.Content,
			&ad.TargetUrl,
		)
		if err != nil {
			return nil, fmt.Errorf(" failed to scan row: %w", err)
		}
		ads = append(ads, ad)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return ads, nil
}
