package profile

import (
	modelpayment "2025_2_404/internal/service/profile/domain"
	"context"
	"fmt"
)

const (
	sqlTextForCreatePayment = `INSERT INTO wallet_top_up (client_wallet_id, amount, status, yoo_payment_id, payment_method)
	SELECT w.id, $2, $3, $4, $5
	FROM client_wallet w
	WHERE w.client_id = $1 `

	sqlTextForUpdatePaymentStatus = `
		UPDATE wallet_top_up
		SET status = $1
		WHERE yoo_payment_id = $2`

	sqlTextForGetPaymentByClientID = `
		SELECT id, amount, status, payment_method
		FROM wallet_top_up
		WHERE client_wallet_id = (
			SELECT id
			FROM client_wallet
			WHERE client_id = $1
		)`
)

func (r *DB) CreatePayment(ctx context.Context, payment modelpayment.Payment) error {
	_, err := r.sql.ExecContext(
		ctx,
		sqlTextForCreatePayment,
		payment.ClientID,
		payment.AmountRub,
		payment.Status,
		payment.YooPaymentID,
		payment.PaymentMethod,
	)
	return err
}

func (r *DB) UpdatePaymentStatus(ctx context.Context, yooPaymentID string, status modelpayment.PaymentStatus) error {
	_, err := r.sql.ExecContext(
		ctx,
		sqlTextForUpdatePaymentStatus,
		status,
		yooPaymentID,
	)
	return err
}

func (r *DB) GetPaymentsByClientID(ctx context.Context, clientID modelpayment.ID) ([]modelpayment.Payment, error) {
	rows, err := r.sql.QueryContext(ctx, sqlTextForGetPaymentByClientID, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to query payments: %w", err)
	}
	defer rows.Close()

	var payments []modelpayment.Payment
	for rows.Next() {
		var p modelpayment.Payment
		err := rows.Scan(
			&p.ID,
			&p.AmountRub,
			&p.Status,
			&p.PaymentMethod,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment row: %w", err)
		}
		payments = append(payments, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return payments, nil
}