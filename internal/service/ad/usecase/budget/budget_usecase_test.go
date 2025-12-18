package budget

import (
	modelad "2025_2_404/internal/service/ad/domain/ad"
	modeluser "2025_2_404/internal/service/ad/domain/user"
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"
)

type mockBudgetRepository struct {
	updateBudgetFunc func(ctx context.Context, adID modelad.ID, clientID modeluser.ID, newBudget uint32) error
}

func (m *mockBudgetRepository) UpdateBudget(ctx context.Context, adID modelad.ID, clientID modeluser.ID, newBudget uint32) error {
	if m.updateBudgetFunc != nil {
		return m.updateBudgetFunc(ctx, adID, clientID, newBudget)
	}
	return nil
}

func TestNew(t *testing.T) {
	mockRepo := &mockBudgetRepository{}
	logger := zap.NewNop()

	useCase := New(mockRepo, logger)

	if useCase == nil {
		t.Fatal("expected useCase to be created, got nil")
	}

	if useCase.budgetRepo == nil {
		t.Error("expected budgetRepo to be set")
	}

	if useCase.logger == nil {
		t.Error("expected logger to be set")
	}
}

func TestUpdateBudget_Success(t *testing.T) {
	called := false
	mockRepo := &mockBudgetRepository{
		updateBudgetFunc: func(ctx context.Context, adID modelad.ID, clientID modeluser.ID, newBudget uint32) error {
			called = true
			if newBudget != 5000 {
				t.Errorf("expected newBudget 5000, got %d", newBudget)
			}
			return nil
		},
	}

	logger := zap.NewNop()
	useCase := New(mockRepo, logger)

	var adID modelad.ID
	var clientID modeluser.ID
	err := useCase.UpdateBudget(context.Background(), adID, clientID, 5000)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !called {
		t.Error("expected repository UpdateBudget to be called")
	}
}

func TestUpdateBudget_RepositoryError(t *testing.T) {
	expectedError := errors.New("database error")
	mockRepo := &mockBudgetRepository{
		updateBudgetFunc: func(ctx context.Context, adID modelad.ID, clientID modeluser.ID, newBudget uint32) error {
			return expectedError
		},
	}

	logger := zap.NewNop()
	useCase := New(mockRepo, logger)

	var adID modelad.ID
	var clientID modeluser.ID
	err := useCase.UpdateBudget(context.Background(), adID, clientID, 1000)

	if err != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}
