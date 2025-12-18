package profile

import (
	modelclient "2025_2_404/internal/service/profile/domain"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type mockRepository struct {
	showFunc                func(ctx context.Context, clientID modelclient.ID) (modelclient.User, error)
	updateFunc              func(ctx context.Context, client modelclient.User) error
	deleteFunc              func(ctx context.Context, clientID modelclient.ID) error
	showBalanceFunc         func(ctx context.Context, clientID modelclient.ID) (uint32, error)
	addBalanceFunc          func(ctx context.Context, clientID modelclient.ID, addAmount uint32) error
	subtractBalanceFunc     func(ctx context.Context, clientID modelclient.ID, subAmount uint32) error
	createPaymentFunc       func(ctx context.Context, payment modelclient.Payment) error
	updatePaymentStatusFunc func(ctx context.Context, yooPaymentID string, status modelclient.PaymentStatus) (modelclient.ID, error)
	getPaymentsByClientIDFunc func(ctx context.Context, clientID modelclient.ID) ([]modelclient.Payment, error)
}

func (m *mockRepository) Show(ctx context.Context, clientID modelclient.ID) (modelclient.User, error) {
	if m.showFunc != nil {
		return m.showFunc(ctx, clientID)
	}
	return modelclient.User{}, nil
}

func (m *mockRepository) Update(ctx context.Context, client modelclient.User) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, client)
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, clientID modelclient.ID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, clientID)
	}
	return nil
}

func (m *mockRepository) ShowBalance(ctx context.Context, clientID modelclient.ID) (uint32, error) {
	if m.showBalanceFunc != nil {
		return m.showBalanceFunc(ctx, clientID)
	}
	return 0, nil
}

func (m *mockRepository) AddBalance(ctx context.Context, clientID modelclient.ID, addAmount uint32) error {
	if m.addBalanceFunc != nil {
		return m.addBalanceFunc(ctx, clientID, addAmount)
	}
	return nil
}

func (m *mockRepository) SubtractBalance(ctx context.Context, clientID modelclient.ID, subAmount uint32) error {
	if m.subtractBalanceFunc != nil {
		return m.subtractBalanceFunc(ctx, clientID, subAmount)
	}
	return nil
}

func (m *mockRepository) CreatePayment(ctx context.Context, payment modelclient.Payment) error {
	if m.createPaymentFunc != nil {
		return m.createPaymentFunc(ctx, payment)
	}
	return nil
}

func (m *mockRepository) UpdatePaymentStatus(ctx context.Context, yooPaymentID string, status modelclient.PaymentStatus) (modelclient.ID, error) {
	if m.updatePaymentStatusFunc != nil {
		return m.updatePaymentStatusFunc(ctx, yooPaymentID, status)
	}
	return uuid.Nil, nil
}

func (m *mockRepository) GetPaymentsByClientID(ctx context.Context, clientID modelclient.ID) ([]modelclient.Payment, error) {
	if m.getPaymentsByClientIDFunc != nil {
		return m.getPaymentsByClientIDFunc(ctx, clientID)
	}
	return nil, nil
}

type mockExternalYooKassaHTTP struct {
	createPaymentFunc func(ctx context.Context, payment modelclient.Payment) (modelclient.PaymentResponse, error)
}

func (m *mockExternalYooKassaHTTP) CreatePayment(ctx context.Context, payment modelclient.Payment) (modelclient.PaymentResponse, error) {
	if m.createPaymentFunc != nil {
		return m.createPaymentFunc(ctx, payment)
	}
	return modelclient.PaymentResponse{}, nil
}

func TestShowBalance_Success(t *testing.T) {
	expectedBalance := uint32(500)
	mockRepo := &mockRepository{
		showBalanceFunc: func(ctx context.Context, clientID modelclient.ID) (uint32, error) {
			return expectedBalance, nil
		},
	}

	useCase := New(mockRepo, nil)
	clientID := uuid.New()

	balance, err := useCase.ShowBalance(context.Background(), clientID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if balance != expectedBalance {
		t.Errorf("expected balance %d, got %d", expectedBalance, balance)
	}
}

func TestAddBalance_Success(t *testing.T) {
	mockRepo := &mockRepository{
		addBalanceFunc: func(ctx context.Context, clientID modelclient.ID, addAmount uint32) error {
			if addAmount != 100 {
				t.Errorf("expected addAmount 100, got %d", addAmount)
			}
			return nil
		},
	}

	useCase := New(mockRepo, nil)
	clientID := uuid.New()

	err := useCase.AddBalance(context.Background(), clientID, 100)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSubtractBalance_Success(t *testing.T) {
	mockRepo := &mockRepository{
		showBalanceFunc: func(ctx context.Context, clientID modelclient.ID) (uint32, error) {
			return 500, nil
		},
		subtractBalanceFunc: func(ctx context.Context, clientID modelclient.ID, subAmount uint32) error {
			if subAmount != 100 {
				t.Errorf("expected subAmount 100, got %d", subAmount)
			}
			return nil
		},
	}

	useCase := New(mockRepo, nil)
	clientID := uuid.New()

	err := useCase.SubtractBalance(context.Background(), clientID, 100)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSubtractBalance_InsufficientBalance(t *testing.T) {
	mockRepo := &mockRepository{
		showBalanceFunc: func(ctx context.Context, clientID modelclient.ID) (uint32, error) {
			return 50, nil
		},
	}

	useCase := New(mockRepo, nil)
	clientID := uuid.New()

	err := useCase.SubtractBalance(context.Background(), clientID, 100)
	if err == nil {
		t.Error("expected error for insufficient balance, got nil")
	}

	if err.Error() != "Subtract more balance" {
		t.Errorf("expected 'Subtract more balance' error, got %v", err)
	}
}

func TestSubtractBalance_ShowBalanceError(t *testing.T) {
	expectedErr := errors.New("database error")
	mockRepo := &mockRepository{
		showBalanceFunc: func(ctx context.Context, clientID modelclient.ID) (uint32, error) {
			return 0, expectedErr
		},
	}

	useCase := New(mockRepo, nil)
	clientID := uuid.New()

	err := useCase.SubtractBalance(context.Background(), clientID, 100)
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestShow_Success(t *testing.T) {
	expectedUser := modelclient.User{
		ID:       uuid.New(),
		UserName: "testuser",
		Email:    "test@example.com",
	}

	mockRepo := &mockRepository{
		showFunc: func(ctx context.Context, clientID modelclient.ID) (modelclient.User, error) {
			return expectedUser, nil
		},
	}

	useCase := New(mockRepo, nil)

	user, err := useCase.Show(context.Background(), expectedUser.ID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if user.ID != expectedUser.ID {
		t.Errorf("expected user ID %v, got %v", expectedUser.ID, user.ID)
	}

	if user.UserName != expectedUser.UserName {
		t.Errorf("expected username %s, got %s", expectedUser.UserName, user.UserName)
	}
}

func TestUpdate_Success(t *testing.T) {
	mockRepo := &mockRepository{
		updateFunc: func(ctx context.Context, client modelclient.User) error {
			if client.UserName != "updateduser" {
				t.Errorf("expected username 'updateduser', got %s", client.UserName)
			}
			return nil
		},
	}

	useCase := New(mockRepo, nil)

	user := modelclient.User{
		ID:       uuid.New(),
		UserName: "updateduser",
		Email:    "updated@example.com",
	}

	err := useCase.Update(context.Background(), user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDelete_Success(t *testing.T) {
	mockRepo := &mockRepository{
		deleteFunc: func(ctx context.Context, clientID modelclient.ID) error {
			return nil
		},
	}

	useCase := New(mockRepo, nil)
	clientID := uuid.New()

	err := useCase.Delete(context.Background(), clientID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDelete_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")
	mockRepo := &mockRepository{
		deleteFunc: func(ctx context.Context, clientID modelclient.ID) error {
			return expectedErr
		},
	}

	useCase := New(mockRepo, nil)
	clientID := uuid.New()

	err := useCase.Delete(context.Background(), clientID)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error wrapping %v, got %v", expectedErr, err)
	}
}
