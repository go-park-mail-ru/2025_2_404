package postgres

import (
	"database/sql"
	"context"
	"fmt"
	modelad "2025_2_404/internal/service/ad/domain/ad"
	modelfullad "2025_2_404/internal/service/ad/domain/ad_full_info"
	modeluser "2025_2_404/internal/service/ad/domain/user"
)

const(
	sqlTextForSelectAds = "SELECT id, title, content, img_path, target_url FROM ad WHERE client_id = $1"
	sqlTextForInsertAds = "INSERT INTO ad (client_id, title, content, img_path, target_url) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	sqlTextForUpdateAds = "UPDATE ad SET title = $1, content = $2, img_path = $3, target_url = $4 WHERE id = $5 AND client_id = $6"
	sqlTextForDeleteAds = "DELETE FROM ad WHERE id = $1 AND client_id = $2"
	sqlTextForFullAdInfo = "SELECT ad.id, ad.title, ad.content, ad.img_path, ad.target_url, COALESCE(ad_detail.amount_for_ad, 0), COALESCE(statistic.clicks, 0), COALESCE(statistic.impressions, 0) FROM ad LEFT JOIN ad_detail ON ad_detail.ad_id = ad.id LEFT JOIN statistic ON statistic.ad_detail_id = ad_detail.id WHERE ad.id = $1"
)

type DB struct {
	sql *sql.DB
}

func New(sql *sql.DB) *DB {
	return &DB{
		sql: sql,
	}
}

func (r *DB) FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelad.Ads, error) {
	rows, err := r.sql.QueryContext(ctx, sqlTextForSelectAds, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query ads by user ID: %w", err)
	}
	defer rows.Close()

	var ads []modelad.Ads
	for rows.Next() {
		var ad modelad.Ads
		err := rows.Scan(&ad.ID, &ad.Title, &ad.Content, &ad.ImagePath, &ad.TargetUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ad: %w", err)
		}
		ads = append(ads, ad)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return ads, nil
}

func (r *DB) GetOneAd(ctx context.Context, adID modelad.ID) (modelfullad.AdFullInfo, error) {
	var adInfo modelfullad.AdFullInfo
	row := r.sql.QueryRowContext(ctx, sqlTextForFullAdInfo, adID)

	err :=  row.Scan(
		&adInfo.ID,
		&adInfo.Title,
		&adInfo.Content,
		&adInfo.ImgPath,
		&adInfo.TargetUrl,
		&adInfo.AmountForAd,
		&adInfo.Clicks,
		&adInfo.Impressions,
	)

	if err != nil {
		return modelfullad.AdFullInfo{}, fmt.Errorf("failed to find an ad: %w", err)
	}

	return adInfo, nil
}

func (r *DB) Create(ctx context.Context, ad modelad.Ads) (error) {
	res, err := r.sql.ExecContext(ctx, sqlTextForInsertAds, ad.ClientID, ad.Title, ad.Content, ad.ImagePath, ad.TargetUrl)
	if err != nil {
		return fmt.Errorf("failed to create ad: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ad with id %v not found", ad.ID)
	}

	return nil
}

func (r *DB) Update(ctx context.Context, ad modelad.Ads) error {

	res, err := r.sql.ExecContext(ctx, sqlTextForUpdateAds, ad.Title, ad.Content, ad.ImagePath, ad.TargetUrl, ad.ID, ad.ClientID)
	if err != nil {
		return fmt.Errorf("failed to update ad: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ad with id %v and client_id %v not found", ad.ID,ad.ClientID)
	}

	return nil
}
func (r *DB) Delete(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error {
	
	result, err := r.sql.ExecContext(ctx, sqlTextForDeleteAds, adID, clientID)
	if err != nil {
		return fmt.Errorf("failed to delete ad: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get a rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Ad with ID %d not found", adID)
	}
	fmt.Printf("Пользователь с ID %d успешно удален. Затронуто строк: %d", adID, rowsAffected)
	return nil
}