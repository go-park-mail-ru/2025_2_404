// Package budget provides use case implementations for managing advertisement budgets.
package budget

import (
	modelad "2025_2_404/internal/service/ad/domain/ad"
	modeluser "2025_2_404/internal/service/ad/domain/user"
	"context"
)

type budgetRepositoryI interface {
	UpdateBudget(ctx context.Context, adID modelad.ID, clientID modeluser.ID, newBudget uint32) error
}

type UseCase struct {
	budgetRepo budgetRepositoryI
}

func New(budgetRepo budgetRepositoryI) *UseCase {
	return &UseCase{
		budgetRepo: budgetRepo,
	}
}

func (u *UseCase) UpdateBudget(ctx context.Context, adID modelad.ID, clientID modeluser.ID, newBudget uint32) error {
	return u.budgetRepo.UpdateBudget(ctx, adID, clientID, newBudget)
}
