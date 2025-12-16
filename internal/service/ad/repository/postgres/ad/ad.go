package postgres

import (
	"database/sql"
	"context"
	"fmt"
	"time"
	modelad "2025_2_404/internal/service/ad/domain/ad"
	modelfullad "2025_2_404/internal/service/ad/domain/ad_full_info"
	modeluser "2025_2_404/internal/service/ad/domain/user"
)

const(
	sqlTextForSelectAds = "SELECT ad.id, ad.title, ad.content, ad.img_path, ad.target_url, COALESCE(ad_detail.budget, 0), COALESCE(ad_detail.status, 'non-active'), ad_detail.start_at, ad_detail.end_at, COALESCE(statistic.clicks, 0), COALESCE(statistic.impressions, 0) FROM ad JOIN ad_detail ON ad_detail.ad_id = ad.id LEFT JOIN statistic ON statistic.ad_detail_id = ad_detail.id WHERE ad.client_id = $1"
	sqlTextForInsertAds = "INSERT INTO ad (client_id, title, content, img_path, target_url) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	// sqlTextForUpdateAds = "UPDATE ad SET title = $1, content = $2, img_path = $3, target_url = $4, budget = $5, status = $6 WHERE id = $7 AND client_id = $8"
	sqlTextForUpdateAds = `UPDATE ad SET title = $1, content = $2, img_path = $3, target_url = $4 WHERE id = $5 AND client_id = $6`
	sqlTextForSaveBudget = "INSERT INTO ad_detail (ad_id, budget, status, start_at, end_at) VALUES ($1, $2, $3, $4, $5)"
	sqlTextForDeleteAds = "DELETE FROM ad WHERE id = $1 AND client_id = $2"
	sqlTextForFullAdInfo = "SELECT ad.id, ad.title, ad.content, ad.img_path, ad.target_url, COALESCE(ad_detail.budget, 0), COALESCE(ad_detail.status, 'non-active'), ad_detail.start_at, ad_detail.end_at, COALESCE(statistic.clicks, 0), COALESCE(statistic.impressions, 0) FROM ad LEFT JOIN ad_detail ON ad_detail.ad_id = ad.id LEFT JOIN statistic ON statistic.ad_detail_id = ad_detail.id WHERE ad.id = $1 AND client_id = $2"
	sqlTextForGetAdDetailID = "UPDATE ad_detail SET budget = ad_detail.budget - 3 WHERE ad_id = $1 RETURNING id "
	sqlTextForGetAdSlot = `
	SELECT id, title, content, img_path, target_url 
	FROM ad 
	WHERE id = (
	SELECT ad_id 
	FROM ad_detail
	WHERE budget >= $1
	AND status = 'active'
	ORDER BY RANDOM()
	LIMIT 1
	)`
	// sqlTextForUpdateAdDetail = `UPDATE ad_detail SET status = $1 WHERE ad_id = $2`
	sqlTextForUpdateAdDetail = `UPDATE ad_detail SET status = $1, start_at = $2, end_at = $3 WHERE ad_id = $4`
	sqlTextForCountAds = "SELECT COUNT(*) FROM ad WHERE client_id = $1"
	sqlTextForUpdateStatistic = "UPDATE statistic SET clicks = statistic.clicks + $1, impressions = statistic.impressions + $2 WHERE ad_detail_id = $3"
	sqlTextForGetPathImage = "SELECT img_path FROM ad WHERE id = $1"
)

type DB struct {
	sql *sql.DB
}

func New(sql *sql.DB) *DB {
	return &DB{
		sql: sql,
	}
}

func (r *DB) FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelfullad.AdFullInfo, error) {
    rows, err := r.sql.QueryContext(ctx, sqlTextForSelectAds, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to query ads: %w", err)
    }
    defer rows.Close()

    var ads []modelfullad.AdFullInfo

    for rows.Next() {
        var adInfo modelfullad.AdFullInfo
        var startAt, endAt sql.NullTime

        err := rows.Scan(
            &adInfo.ID,
            &adInfo.Title,
            &adInfo.Content,
            &adInfo.ImgPath,
            &adInfo.TargetUrl,
            &adInfo.Budget,
            &adInfo.Status,
            &startAt,
            &endAt,
            &adInfo.Clicks,
            &adInfo.Impressions,
        )
        if err != nil {
            return nil, fmt.Errorf("scan error: %w", err)
        }

        if startAt.Valid { adInfo.StartAt = startAt.Time }
        if endAt.Valid { adInfo.EndAt = endAt.Time }

        ads = append(ads, adInfo)
    }

    return ads, nil
}

func (r *DB) GetOneAd(ctx context.Context, adID modelad.ID, clientID modeluser.ID) (modelfullad.AdFullInfo, error) {
	var adInfo modelfullad.AdFullInfo
	row := r.sql.QueryRowContext(ctx, sqlTextForFullAdInfo, adID, clientID)

	err :=  row.Scan(
		&adInfo.ID,
		&adInfo.Title,
		&adInfo.Content,
		&adInfo.ImgPath,
		&adInfo.TargetUrl,
		&adInfo.Budget,
		&adInfo.Status,
		&adInfo.StartAt,
		&adInfo.EndAt,
		&adInfo.Clicks,
		&adInfo.Impressions,
	)

	if err != nil {
		return modelfullad.AdFullInfo{}, fmt.Errorf("failed to find an ad: %w", err)
	}

	return adInfo, nil
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

func (r *DB) Create(ctx context.Context, ad modelad.Ads) error {
	var newAdID modelad.ID
	err := r.sql.QueryRowContext(ctx, sqlTextForInsertAds, ad.ClientID, ad.Title, ad.Content, ad.ImagePath, ad.TargetUrl).Scan(&newAdID)
	if err != nil {
		return fmt.Errorf("failed to create ad: %w", err)
	}

	if ad.StartAt.IsZero() {
		ad.StartAt = time.Now()
	}
	
	if ad.EndAt.IsZero() {
		ad.EndAt = ad.StartAt.Add(time.Hour * 24 * 7)
	}

	_, err = r.sql.ExecContext(ctx, sqlTextForSaveBudget, newAdID, ad.Budget, ad.Status, ad.StartAt, ad.EndAt)
	if err != nil {
		r.Delete(ctx, ad.ID, ad.ClientID)
		return fmt.Errorf("failed to save ad budget: %w", err)
	}

	return nil
}

// func (r *DB) Update(ctx context.Context, ad modelad.Ads) error {
// 	tx, err := r.sql.BeginTx(ctx, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to begin transaction: %w", err)
// 	}
// 	defer tx.Rollback()

// 	_, err = tx.ExecContext(ctx,sqlTextForUpdateAds,
// 		ad.Title, ad.Content, ad.ImagePath, ad.TargetUrl, ad.ID, ad.ClientID,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to update ad: %w", err)
// 	}

// 	res, err := tx.ExecContext(ctx, sqlTextForUpdateAdDetail, ad.Status, ad.ID)
// 	if err != nil {
// 		return fmt.Errorf("failed to update ad_detail: %w", err)
// 	}

// 	rowsAffected, err := res.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to get rows affected: %w", err)
// 	}
// 	if rowsAffected == 0 {
// 		return fmt.Errorf("ad_detail for ad_id %v not found", ad.ID)
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		return fmt.Errorf("failed to commit transaction: %w", err)
// 	}

// 	return nil
// }


func (r *DB) Update(ctx context.Context, ad modelad.Ads) error {
	tx, err := r.sql.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	if ad.ImagePath == ""{
		err := tx.QueryRowContext(ctx, sqlTextForGetPathImage, ad.ID).Scan(&ad.ImagePath)
		if err != nil{
			return fmt.Errorf("failed to select image for ad: %w", err)
		}
	}
	res, err := tx.ExecContext(ctx, sqlTextForUpdateAds,
		ad.Title, ad.Content, ad.ImagePath, ad.TargetUrl, ad.ID, ad.ClientID,
	)
	if err != nil {
		return fmt.Errorf("failed to update ad: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ad with id %v not found or access denied", ad.ID)
	}
	
	_, err = tx.ExecContext(ctx, sqlTextForUpdateAdDetail, ad.Status, ad.StartAt, ad.EndAt, ad.ID)
	if err != nil {
		return fmt.Errorf("failed to update ad_detail: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *DB) GetAdDetailForSlot(ctx context.Context, id modelad.ID, click, impression int) (modelfullad.DetailID, error){
	var detail_id modelfullad.DetailID
	err := r.sql.QueryRowContext(ctx, sqlTextForGetAdDetailID, id).Scan(&detail_id)
	if err != nil{
		return modelfullad.DetailID{}, fmt.Errorf("not found ad_detail_id")
	}

	return detail_id, nil
}

func (r *DB) GetAdSlot(ctx context.Context, min_cost uint32) (modelad.Ads, error) {
	var adSlot modelad.Ads
	err := r.sql.QueryRowContext(ctx, sqlTextForGetAdSlot, min_cost).Scan(
		&adSlot.ID, 
		&adSlot.Title,
		&adSlot.Content,
		&adSlot.ImagePath,
		&adSlot.TargetUrl,
	)

	if err != nil{
		return modelad.Ads{}, fmt.Errorf("not found ad for slot")
	}

	return adSlot, nil
}

func (r *DB) GetAdCount(ctx context.Context, clientID modeluser.ID) (int64, error) {
	var count int64
	err := r.sql.QueryRowContext(ctx, sqlTextForCountAds, clientID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count ads: %w", err)
	}

	return count, nil
}