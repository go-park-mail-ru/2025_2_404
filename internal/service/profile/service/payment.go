package profile

import (
	modelpayment "2025_2_404/internal/service/profile/domain"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

func (u *UseCase) CreatePayment(ctx context.Context, payment modelpayment.Payment) (string, error) {
	// Генерация ID платежа
	yooKassaID := uuid.New().String()
	payment.YooPaymentID = yooKassaID

	slog.Debug("Создание платежа: начало",
		"user_id", payment.ClientID,
		"amount", payment.AmountRub,
		"yoo_payment_id", yooKassaID,
	)

	// Вызов внешнего платежного сервиса (YooKassa)
	yooKassaResp, err := u.ext.CreatePayment(ctx, payment)
	if err != nil {
		slog.Error("Ошибка при создании платежа во внешнем сервисе",
			"yoo_payment_id", yooKassaID,
			"user_id", payment.ClientID,
			"amount", payment.AmountRub,
			"error", err,
		)
		return "", fmt.Errorf("failed to create payment in external service: %w", err)
	}

	slog.Debug("Платеж успешно создан во внешнем сервисе",
		"yoo_payment_id", yooKassaResp.ID,
		"payment_link", yooKassaResp.Confirmation.URL,
	)

	payment.YooPaymentID = yooKassaResp.ID
	payment.Status = modelpayment.PaymentPending
	yooKassaLink := yooKassaResp.Confirmation.URL

	// Сохраняем в репозиторий
	if err := u.repo.CreatePayment(ctx, payment); err != nil {
		slog.Error("Ошибка при сохранении платежа в БД",
			"yoo_payment_id", yooKassaID,
			"user_id", payment.ClientID,
			"error", err,
		)
		return "", fmt.Errorf("failed to store payment in repository: %w", err)
	}

	slog.Info("Платёж успешно создан и сохранён",
		"yoo_payment_id", yooKassaID,
		"user_id", payment.ClientID,
		"amount", payment.AmountRub,
		"status", payment.Status,
	)

	return yooKassaLink, nil
}

func (u *UseCase) UpdatePaymentStatus(ctx context.Context, yooPaymentID string, status modelpayment.PaymentStatus) (modelpayment.ID, error) {
	return u.repo.UpdatePaymentStatus(ctx, yooPaymentID, status)
}

func (u *UseCase) GetPaymentsByClientID(ctx context.Context, clientID modelpayment.ID) ([]modelpayment.Payment, error) {
	return u.repo.GetPaymentsByClientID(ctx, clientID)
}
