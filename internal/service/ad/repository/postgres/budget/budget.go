package budget

import (
	modelad "2025_2_404/internal/service/ad/domain/ad"
	modeluser "2025_2_404/internal/service/ad/domain/user"
	"context"
	"database/sql"
	"fmt"
)

const(
	// 	sqlTextForSelectBudget = "SELECT COALESCE(ad_detail.budget, 0) FROM ad LEFT JOIN ad_detail ON ad_detail.ad_id = ad.id WHERE ad.id = $1 AND ad.client_id = $2"
	sqlTextForUpdateBudget = "UPDATE ad_detail SET budget =$1 FROM ad WHERE ad_detail.ad_id = ad.id AND ad.id = $2 AND ad.client_id = $3"
)

type DB struct {
	sql *sql.DB
}

func New(sql *sql.DB) *DB {
	return &DB{
		sql: sql,
	}
}

func (r *DB) UpdateBudget(ctx context.Context, adID modelad.ID, clientID modeluser.ID, newBudget uint32) (error) {
	res, err := r.sql.ExecContext(ctx, sqlTextForUpdateBudget, newBudget, adID, clientID)
    if err != nil {
        return fmt.Errorf("failed to execute update budget query: %w", err)
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("ad not found or access denied")
    }

    return nil
}

