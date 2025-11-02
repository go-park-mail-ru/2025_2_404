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
	sqlTextForUpdateAds = "UPDATE ad SET title = $1, content = $2, img_bin = $3, target_url = $4 WHERE id = $5"
	sqlTextForDeleteAds = "DELETE FROM ad WHERE id = $1"
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
		err := rows.Scan(&ad.ID, &ad.Title, &ad.Content, &ad.ImgBin, &ad.TargetUrl)
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

func (r *DB) Create(ctx context.Context, ad modelad.Ads) (error) {
	res, err := r.sql.ExecContext(ctx, sqlTextForInsertAds, ad.ClientID, ad.Title, ad.Content, ad.ImgBin, ad.TargetUrl)
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

	res, err := r.sql.ExecContext(ctx, sqlTextForUpdateAds, ad.Title, ad.Content, ad.ImgBin, ad.TargetUrl, ad.ID, ad.ClientID)
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
func (r *DB) Delete(ctx context.Context, adID int64) error {
	
	result, err := r.sql.ExecContext(ctx, sqlTextForDeleteAds, adID)
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
