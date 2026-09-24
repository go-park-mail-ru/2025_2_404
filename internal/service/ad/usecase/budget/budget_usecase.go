package budget

import (
	modelad "2025_2_404/internal/service/ad/domain/ad"
	modeluser "2025_2_404/internal/service/ad/domain/user"
	"context"
	"fmt"
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
	err := u.budgetRepo.UpdateBudget(ctx, adID, clientID, newBudget)
    if err != nil {
        return fmt.Errorf("failed to update budget in usecase: %w", err)
    }

    return nil
}