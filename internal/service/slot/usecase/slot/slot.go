// Package slot provides use case implementations for managing slots.
package slot

import (
	"2025_2_404/internal/service/slot/domain/slot"
	"context"
)

type slotRepository interface {
	Create(ctx context.Context, s slot.Slot) (slot.ID, error)
	GetByID(ctx context.Context, id slot.ID) (slot.Slot, error)
	ListByUserID(ctx context.Context, userID slot.UserID) ([]slot.Slot, error)
	Update(ctx context.Context, s slot.Slot) error
	Delete(ctx context.Context, id slot.ID, userID slot.UserID) error
}

type UseCase struct {
	repo slotRepository
}

func New(repo slotRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) Create(ctx context.Context, s slot.Slot) (slot.ID, error) {
	return u.repo.Create(ctx, s)
}

func (u *UseCase) GetByID(ctx context.Context, id slot.ID) (slot.Slot, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *UseCase) ListByUserID(ctx context.Context, userID slot.UserID) ([]slot.Slot, error) {
	return u.repo.ListByUserID(ctx, userID)
}

func (u *UseCase) Update(ctx context.Context, s slot.Slot) error {
	return u.repo.Update(ctx, s)
}

func (u *UseCase) Delete(ctx context.Context, id slot.ID, userID slot.UserID) error {
	return u.repo.Delete(ctx, id, userID)
}
