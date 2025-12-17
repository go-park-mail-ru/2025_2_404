// Package postgres provides a PostgreSQL implementation of the advertisement repository.
package postgres

import (
	modelad "2025_2_404/internal/service/ad/domain/ad"
	modelfullad "2025_2_404/internal/service/ad/domain/ad_full_info"
	modeluser "2025_2_404/internal/service/ad/domain/user"
	"2025_2_404/pkg/globalerrors"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"go.uber.org/zap"
)

const (
	sqlTextForSelectAds = `
		SELECT 
			ad.id, 
			ad.title, 
			COALESCE(ad_detail.status, 'non-active') AS status,
			ad.created_at
		FROM ad 
		JOIN ad_detail ON ad_detail.ad_id = ad.id 
		WHERE ad.client_id = $1`
	sqlTextForInsertAds     = "INSERT INTO ad (client_id, title, content, img_path, target_url) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	sqlTextForUpdateAds     = `UPDATE ad SET title = $1, content = $2, img_path = $3, target_url = $4 WHERE id = $5 AND client_id = $6`
	sqlTextForSaveBudget    = "INSERT INTO ad_detail (ad_id, budget, status, start_at, end_at) VALUES ($1, $2, $3, $4, $5)"
	sqlTextForDeleteAds     = "DELETE FROM ad WHERE id = $1 AND client_id = $2"
	sqlTextForFullAdInfo    = "SELECT ad.id, ad.title, ad.content, ad.img_path, ad.target_url, COALESCE(ad_detail.budget, 0), COALESCE(ad_detail.status, 'non-active'), ad_detail.start_at, ad_detail.end_at, COALESCE(statistic.clicks, 0), COALESCE(statistic.impressions, 0) FROM ad LEFT JOIN ad_detail ON ad_detail.ad_id = ad.id LEFT JOIN statistic ON statistic.ad_detail_id = ad_detail.id WHERE ad.id = $1 AND client_id = $2"
	sqlTextForGetAdDetailID = "UPDATE ad_detail SET budget = ad_detail.budget - 3 WHERE ad_id = $1 RETURNING id "
	sqlTextForGetAdSlot     = `
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
	sqlTextForUpdateAdDetail  = `UPDATE ad_detail SET status = $1, start_at = $2, end_at = $3 WHERE ad_id = $4`
	sqlTextForCountAds        = "SELECT COUNT(*) FROM ad WHERE client_id = $1"
	sqlTextForUpdateStatistic = "UPDATE statistic SET clicks = statistic.clicks + $1, impressions = statistic.impressions + $2 WHERE ad_detail_id = $3"
	sqlTextForGetPathImage    = "SELECT img_path FROM ad WHERE id = $1"
)

type DB struct {
	sql    *sql.DB
	logger *zap.Logger
}

func New(sql *sql.DB, logger *zap.Logger) *DB {
	return &DB{
		sql:    sql,
		logger: logger,
	}
}

func (r *DB) FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelfullad.AdFullInfo, error) {
	rows, err := r.sql.QueryContext(ctx, sqlTextForSelectAds, userID)
	if err != nil {
		r.logger.Error("failed to query ads by user ID", zap.Error(err))
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "42P01":
				return nil, globalerrors.ErrInternal
			case "42703":
				return nil, globalerrors.ErrInternal
			}
		}
		return nil, globalerrors.ErrInternal
	}
	defer func() {
		_ = rows.Close()
	}()

	var ads []modelfullad.AdFullInfo
	for rows.Next() {
		var adInfo modelfullad.AdFullInfo
		var createAt sql.NullTime

		if err := rows.Scan(&adInfo.ID, &adInfo.Title, &adInfo.Status, &createAt); err != nil {
			r.logger.Error("failed to scan ad row", zap.Error(err))
			return nil, globalerrors.ErrInternal
		}
		if createAt.Valid {
			adInfo.CreatedAt = createAt.Time
		}
		ads = append(ads, adInfo)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("rows iteration error", zap.Error(err))
		return nil, globalerrors.ErrInternal
	}

	return ads, nil
}

func (r *DB) GetOneAd(ctx context.Context, adID modelad.ID, clientID modeluser.ID) (modelfullad.AdFullInfo, error) {
	var adInfo modelfullad.AdFullInfo
	err := r.sql.QueryRowContext(ctx, sqlTextForFullAdInfo, adID, clientID).Scan(
		&adInfo.ID,
		&adInfo.Title,
		&adInfo.Content,
		&adInfo.ImgPath,
		&adInfo.TargetURL,
		&adInfo.Budget,
		&adInfo.Status,
		&adInfo.StartAt,
		&adInfo.EndAt,
		&adInfo.Clicks,
		&adInfo.Impressions,
	)
	if err != nil {
		r.logger.Error("failed to get ad by ID", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return modelfullad.AdFullInfo{}, globalerrors.ErrAdNotFound
		}
		return modelfullad.AdFullInfo{}, globalerrors.ErrInternal
	}
	r.logger.Info("ad found", zap.String("ad_id", adID.String()), zap.String("client_id", clientID.String()))
	return adInfo, nil
}

func (r *DB) Delete(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error {
	result, err := r.sql.ExecContext(ctx, sqlTextForDeleteAds, adID, clientID)
	if err != nil {
		r.logger.Error("ad not delete", zap.Error(err))
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return globalerrors.ErrForeignKeyViolation
			}
		}
		return globalerrors.ErrInternal
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("failed to get rows affected on delete", zap.Error(err))
		return globalerrors.ErrInternal
	}
	if rowsAffected == 0 {
		return globalerrors.ErrAdNotFound
	}
	return nil
}

func (r *DB) Create(ctx context.Context, ad modelad.Ads) error {
	var newAdID modelad.ID
	err := r.sql.QueryRowContext(ctx, sqlTextForInsertAds, ad.ClientID, ad.Title, ad.Content, ad.ImagePath, ad.TargetURL).Scan(&newAdID)
	if err != nil {
		r.logger.Error("failed to insert ad", zap.Error(err))
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return globalerrors.ErrForeignKeyViolation
			case "23514":
				return globalerrors.ErrInvalidQuery
			}
		}
		return globalerrors.ErrInternal
	}

	if ad.StartAt.IsZero() {
		ad.StartAt = time.Now()
	}
	if ad.EndAt.IsZero() {
		ad.EndAt = ad.StartAt.Add(7 * 24 * time.Hour)
	}

	_, err = r.sql.ExecContext(ctx, sqlTextForSaveBudget, newAdID, ad.Budget, ad.Status, ad.StartAt, ad.EndAt)
	if err != nil {
		r.logger.Warn("rolling back ad creation due to budget insert failure",
			zap.String("ad_id", newAdID.String()),
			zap.Error(err),
		)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return globalerrors.ErrInternal
			case "23514":
				return globalerrors.ErrInvalidQuery
			}
		}
		_ = r.Delete(ctx, newAdID, ad.ClientID)
		return globalerrors.ErrInternal
	}
	r.logger.Debug("ad created in DB",
		zap.String("ad_id", newAdID.String()),
		zap.String("client_id", ad.ClientID.String()),
	)

	return nil
}

func (r *DB) Update(ctx context.Context, ad modelad.Ads) error {
	tx, err := r.sql.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("failed to begin transaction for update", zap.Error(err))
		return globalerrors.ErrInternal
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if ad.ImagePath == "" {
		if err := tx.QueryRowContext(ctx, sqlTextForGetPathImage, ad.ID).Scan(&ad.ImagePath); err != nil {
			r.logger.Error("failed to get image path during update", zap.Error(err))
			if errors.Is(err, sql.ErrNoRows) {
				return globalerrors.ErrAdNotFound
			}
			return globalerrors.ErrInternal
		}
	}

	res, err := tx.ExecContext(ctx, sqlTextForUpdateAds, ad.Title, ad.Content, ad.ImagePath, ad.TargetURL, ad.ID, ad.ClientID)
	if err != nil {
		r.logger.Error("failed to update ad", zap.Error(err))
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23514":
				return globalerrors.ErrInvalidQuery
			}
		}
		return globalerrors.ErrInternal
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.logger.Error("failed to get rows affected on update", zap.Error(err))
		return globalerrors.ErrInternal
	}
	if rowsAffected == 0 {
		return globalerrors.ErrAccessDenied
	}

	_, err = tx.ExecContext(ctx, sqlTextForUpdateAdDetail, ad.Status, ad.StartAt, ad.EndAt, ad.ID)
	if err != nil {
		r.logger.Error("failed to update ad detail", zap.Error(err))
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23514":
				return globalerrors.ErrInvalidQuery
			}
		}
		return globalerrors.ErrInternal
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit update transaction", zap.Error(err))
		return globalerrors.ErrInternal
	}
	r.logger.Debug("ad updated successfully",
		zap.String("ad_id", ad.ID.String()),
		zap.String("client_id", ad.ClientID.String()),
	)
	return nil
}

func (r *DB) GetAdDetailForSlot(ctx context.Context, id modelad.ID, click, impression int) (modelfullad.DetailID, error) {
	tx, err := r.sql.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("failed to begin transaction for GetAdDetailForSlot", zap.Error(err))
		return modelfullad.DetailID{}, globalerrors.ErrInternal
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var detailID modelfullad.DetailID
	err = tx.QueryRowContext(ctx, sqlTextForGetAdDetailID, id).Scan(&detailID)
	if err != nil {
		r.logger.Error("failed to get ad detail ID", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return modelfullad.DetailID{}, globalerrors.ErrBudgetTooLow
		}
		return modelfullad.DetailID{}, globalerrors.ErrInternal
	}

	_, err = tx.ExecContext(ctx, sqlTextForUpdateStatistic, click, impression, detailID)
	if err != nil {
		r.logger.Error("failed to update statistic", zap.Error(err))
		return modelfullad.DetailID{}, globalerrors.ErrInternal
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit GetAdDetailForSlot transaction", zap.Error(err))
		return modelfullad.DetailID{}, globalerrors.ErrInternal
	}

	return detailID, nil
}

func (r *DB) GetAdSlot(ctx context.Context, minCost uint32) (modelad.Ads, error) {
	var adSlot modelad.Ads
	err := r.sql.QueryRowContext(ctx, sqlTextForGetAdSlot, minCost).Scan(
		&adSlot.ID,
		&adSlot.Title,
		&adSlot.Content,
		&adSlot.ImagePath,
		&adSlot.TargetURL,
	)
	if err != nil {
		r.logger.Error("failed to get ad slot", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return modelad.Ads{}, globalerrors.ErrAdNotFound
		}
		return modelad.Ads{}, globalerrors.ErrInternal
	}
	return adSlot, nil
}

func (r *DB) GetAdCount(ctx context.Context, clientID modeluser.ID) (int64, error) {
	var count int64
	if err := r.sql.QueryRowContext(ctx, sqlTextForCountAds, clientID).Scan(&count); err != nil {
		r.logger.Error("failed to get ad count", zap.Error(err))
		return 0, globalerrors.ErrInternal
	}
	return count, nil
}
