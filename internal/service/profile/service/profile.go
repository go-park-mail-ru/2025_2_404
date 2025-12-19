package profile

import (
	modeluser "2025_2_404/internal/service/profile/domain"
	"context"
	"fmt"
)

type repositoryI interface {
	Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, error)
	Update(ctx context.Context, client modeluser.User) error
	Delete(ctx context.Context, clientID modeluser.ID) error
	ShowBalance(ctx context.Context, clientID modeluser.ID) (uint32, error)
	AddBalance(ctx context.Context, clientID modeluser.ID, addAmount uint32) error
	SubtractBalance(ctx context.Context, clientID modeluser.ID, subAmount uint32) error
	CreatePayment(ctx context.Context, payment modeluser.Payment) error
	UpdatePaymentStatus(ctx context.Context, yooPaymentID string, status modeluser.PaymentStatus) (modeluser.ID, error)
	GetPaymentsByClientID(ctx context.Context, clientID modeluser.ID) ([]modeluser.Payment, error)
}

type externalYooKassaHTTPI interface {
	CreatePayment(ctx context.Context, payment modeluser.Payment) (modeluser.PaymentResponse, error)
}

type UseCase struct {
	repo repositoryI
	ext  externalYooKassaHTTPI
}

func New(repo repositoryI, ext externalYooKassaHTTPI) *UseCase {
	return &UseCase{
		repo: repo,
		ext:  ext,
	}
}

func (u *UseCase) Update(ctx context.Context, client modeluser.User) error {
	return u.repo.Update(ctx, client)
}

func (u *UseCase) Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, error) {
	return u.repo.Show(ctx, clientID)
}

func (u *UseCase) Delete(ctx context.Context, clientID modeluser.ID) error {

	if err := u.repo.Delete(ctx, clientID); err != nil {
		return fmt.Errorf("failed to delete user from repository: %w", err)
	}

	return nil
}
