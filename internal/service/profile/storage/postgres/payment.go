package profile

import (
	modelpayment "2025_2_404/internal/service/profile/domain"
	"context"
	"fmt"
	"log/slog"
	"time"
)

const (
	sqlTextForCreatePayment = `INSERT INTO wallet_top_up (client_wallet_id, amount, status, yoo_payment_id, payment_method)
	SELECT w.id, $2, $3, $4, $5
	FROM client_wallet w
	WHERE w.client_id = $1 `

	sqlTextForUpdatePaymentStatus = `
		UPDATE wallet_top_up
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		FROM client_wallet
		WHERE wallet_top_up.yoo_payment_id = $2
		AND wallet_top_up.client_wallet_id = client_wallet.id
		RETURNING client_wallet.client_id;`

	sqlTextForGetPaymentByClientID = `
		SELECT id, amount, status, payment_method, created_at
		FROM wallet_top_up
		WHERE client_wallet_id = (
			SELECT id
			FROM client_wallet
			WHERE client_id = $1
		) AND status = 'succeeded'`
)

func (r *DB) CreatePayment(ctx context.Context, payment modelpayment.Payment) error {
	slog.Info("Creating payment record",
		"client_id", payment.ClientID,
		"amount_rub", payment.AmountRub,
		"status", payment.Status,
		"yoo_payment_id", payment.YooPaymentID,
		"payment_method", payment.PaymentMethod,
	)

	_, err := r.sql.ExecContext(
		ctx,
		sqlTextForCreatePayment,
		payment.ClientID,
		payment.AmountRub,
		payment.Status,
		payment.YooPaymentID,
		payment.PaymentMethod,
	)
	if err != nil {
		slog.Error(" Failed to create payment",
			"yoo_payment_id", payment.YooPaymentID,
			"error", err,
		)
		return err
	}

	slog.Info(" Payment record created",
		"yoo_payment_id", payment.YooPaymentID,
		"client_id", payment.ClientID,
	)
	return nil
}

func (r *DB) UpdatePaymentStatus(ctx context.Context, yooPaymentID string, status modelpayment.PaymentStatus) (modelpayment.ID, error) {
	slog.Info(" Updating payment status",
		"yoo_payment_id", yooPaymentID,
		"new_status", status,
	)

	var clientID modelpayment.ID
	err := r.sql.QueryRowContext(
		ctx,
		sqlTextForUpdatePaymentStatus,
		status,
		yooPaymentID,
	).Scan(&clientID)
	if err != nil {
		slog.Error(" Failed to update payment status or get client_id",
			"yoo_payment_id", yooPaymentID,
			"status", status,
			"error", err,
		)
		return modelpayment.ID{}, err
	}

	slog.Info(" Payment status updated",
		"yoo_payment_id", yooPaymentID,
		"client_id", clientID,
		"status", status,
	)
	return clientID, nil
}

func (r *DB) GetPaymentsByClientID(ctx context.Context, clientID modelpayment.ID) ([]modelpayment.Payment, error) {
	slog.Info(" Fetching payments by client_id", "client_id", clientID)

	rows, err := r.sql.QueryContext(ctx, sqlTextForGetPaymentByClientID, clientID)
	if err != nil {
		slog.Error(" Failed to query payments by client_id",
			"client_id", clientID,
			"error", err,
		)
		return nil, fmt.Errorf("failed to query payments: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var payments []modelpayment.Payment
	for rows.Next() {
		var p modelpayment.Payment
		var createdTime time.Time
		err := rows.Scan(
			&p.ID,
			&p.AmountRub,
			&p.Status,
			&p.PaymentMethod,
			&createdTime,
		)
		if err != nil {
			slog.Error(" Failed to scan payment row",
				"client_id", clientID,
				"error", err,
			)
			return nil, fmt.Errorf("failed to scan payment row: %w", err)
		}
		p.CreatedTime = createdTime.Format(time.RFC3339)
		payments = append(payments, p)
	}

	if err = rows.Err(); err != nil {
		slog.Error(" Rows iteration error",
			"client_id", clientID,
			"error", err,
		)
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	slog.Info(" Retrieved payments",
		"client_id", clientID,
		"count", len(payments),
	)
	return payments, nil
}
