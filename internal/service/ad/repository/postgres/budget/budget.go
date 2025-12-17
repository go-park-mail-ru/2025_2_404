// Package budget provides a PostgreSQL implementation of the budget repository.
package budget

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"
	"go.uber.org/zap"

	modelad "2025_2_404/internal/service/ad/domain/ad"
	modeluser "2025_2_404/internal/service/ad/domain/user"
	"2025_2_404/pkg/globalerrors"
)

const (
	// 	sqlTextForSelectBudget = "SELECT COALESCE(ad_detail.budget, 0) FROM ad LEFT JOIN ad_detail ON ad_detail.ad_id = ad.id WHERE ad.id = $1 AND ad.client_id = $2"
	sqlTextForUpdateBudget = "UPDATE ad_detail SET budget = ad_detail.budget + $1 FROM ad WHERE ad_detail.ad_id = ad.id AND ad.id = $2 AND ad.client_id = $3"
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

func (r *DB) UpdateBudget(ctx context.Context, adID modelad.ID, clientID modeluser.ID, newBudget uint32) error {
	r.logger.Debug("executing UpdateBudget query",
		zap.String("ad_id", adID.String()),
		zap.String("client_id", clientID.String()),
		zap.Uint32("new_budget", newBudget),
	)
	res, err := r.sql.ExecContext(ctx, sqlTextForUpdateBudget, newBudget, adID, clientID)
	if err != nil {
		r.logger.Error("failed to execute UpdateBudget query",
			zap.String("ad_id", adID.String()),
			zap.String("client_id", clientID.String()),
			zap.Uint32("new_budget", newBudget),
			zap.Error(err),
		)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23514":
				return globalerrors.ErrInvalidQuery
			case "23503":
				return globalerrors.ErrAdNotFound
			}
		}
		return globalerrors.ErrInternal
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.logger.Error("failed to get rows affected in UpdateBudget",
			zap.String("ad_id", adID.String()),
			zap.String("client_id", clientID.String()),
			zap.Error(err),
		)
		return globalerrors.ErrInternal
	}

	if rowsAffected == 0 {
		r.logger.Debug("UpdateBudget affected 0 rows — ad not found or access denied",
			zap.String("ad_id", adID.String()),
			zap.String("client_id", clientID.String()),
		)
		return globalerrors.ErrAccessDenied
	}

	r.logger.Debug("budget updated successfully in DB",
		zap.String("ad_id", adID.String()),
		zap.String("client_id", clientID.String()),
		zap.Uint32("new_budget", newBudget),
		zap.Int64("rows_affected", rowsAffected),
	)

	return nil
}
