package profile

import (
	modelpayment "2025_2_404/internal/service/profile/domain"
	"context"

	"github.com/google/uuid"
)

func (u *UseCase) CreatePayment(ctx context.Context, payment modelpayment.Payment) (string, error) {
	yooKassaID := uuid.New().String()
	payment.YooPaymentID = yooKassaID
	yooKassaLink, err := u.ext.CreatePayment(ctx, payment)
	if err != nil{
		return "", err
	}

	payment.Status = modelpayment.PaymentPending

	err = u.repo.CreatePayment(ctx, payment)
	if err != nil{
		return "", err
	}

	return yooKassaLink, nil
}

func (u *UseCase) UpdatePaymentStatus(ctx context.Context, yooPaymentID string, status modelpayment.PaymentStatus) error {
	return u.repo.UpdatePaymentStatus(ctx, yooPaymentID, status)
}

func (u *UseCase) GetPaymentsByClientID(ctx context.Context, clientID modelpayment.ID) ([]modelpayment.Payment, error) {
	return u.repo.GetPaymentsByClientID(ctx, clientID)
}
