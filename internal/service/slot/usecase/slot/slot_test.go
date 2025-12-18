package slot

import (
	"2025_2_404/internal/service/slot/domain/slot"
	"context"
	"errors"
	"testing"
)

type mockSlotRepository struct {
	createFunc       func(ctx context.Context, s slot.Slot) (slot.ID, error)
	getByIDFunc      func(ctx context.Context, id slot.ID) (slot.Slot, error)
	listByUserIDFunc func(ctx context.Context, userID slot.UserID) ([]slot.Slot, error)
	updateFunc       func(ctx context.Context, s slot.Slot) error
	deleteFunc       func(ctx context.Context, id slot.ID, userID slot.UserID) error
}

func (m *mockSlotRepository) Create(ctx context.Context, s slot.Slot) (slot.ID, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, s)
	}
	return "", nil
}

func (m *mockSlotRepository) GetByID(ctx context.Context, id slot.ID) (slot.Slot, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return slot.Slot{}, nil
}

func (m *mockSlotRepository) ListByUserID(ctx context.Context, userID slot.UserID) ([]slot.Slot, error) {
	if m.listByUserIDFunc != nil {
		return m.listByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockSlotRepository) Update(ctx context.Context, s slot.Slot) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, s)
	}
	return nil
}

func (m *mockSlotRepository) Delete(ctx context.Context, id slot.ID, userID slot.UserID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id, userID)
	}
	return nil
}

func TestCreate_Success(t *testing.T) {
	expectedID := slot.ID("slot-123")
	mockRepo := &mockSlotRepository{
		createFunc: func(ctx context.Context, s slot.Slot) (slot.ID, error) {
			return expectedID, nil
		},
	}

	useCase := New(mockRepo)

	newSlot := slot.Slot{
		UserID:         "user-1",
		SlotName:       "Test Slot",
		MinCostAdv:     100,
		FormatOfBanner: "banner",
		Status:         "active",
	}

	id, err := useCase.Create(context.Background(), newSlot)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if id != expectedID {
		t.Errorf("expected ID %v, got %v", expectedID, id)
	}
}

func TestCreate_Error(t *testing.T) {
	expectedErr := errors.New("database error")
	mockRepo := &mockSlotRepository{
		createFunc: func(ctx context.Context, s slot.Slot) (slot.ID, error) {
			return "", expectedErr
		},
	}

	useCase := New(mockRepo)

	newSlot := slot.Slot{
		UserID:         "user-1",
		SlotName:       "Test Slot",
		MinCostAdv:     100,
		FormatOfBanner: "banner",
	}

	_, err := useCase.Create(context.Background(), newSlot)
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestGetByID_Success(t *testing.T) {
	expectedSlot := slot.Slot{
		ID:             "slot-123",
		UserID:         "user-1",
		SlotName:       "Test Slot",
		MinCostAdv:     100,
		FormatOfBanner: "banner",
		Status:         "active",
	}

	mockRepo := &mockSlotRepository{
		getByIDFunc: func(ctx context.Context, id slot.ID) (slot.Slot, error) {
			if id == expectedSlot.ID {
				return expectedSlot, nil
			}
			return slot.Slot{}, errors.New("not found")
		},
	}

	useCase := New(mockRepo)

	result, err := useCase.GetByID(context.Background(), expectedSlot.ID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if result.ID != expectedSlot.ID {
		t.Errorf("expected slot ID %v, got %v", expectedSlot.ID, result.ID)
	}

	if result.SlotName != expectedSlot.SlotName {
		t.Errorf("expected slot name %s, got %s", expectedSlot.SlotName, result.SlotName)
	}
}

func TestListByUserID_Success(t *testing.T) {
	userID := slot.UserID("user-1")
	expectedSlots := []slot.Slot{
		{ID: "slot-1", UserID: userID, SlotName: "Slot 1"},
		{ID: "slot-2", UserID: userID, SlotName: "Slot 2"},
		{ID: "slot-3", UserID: userID, SlotName: "Slot 3"},
	}

	mockRepo := &mockSlotRepository{
		listByUserIDFunc: func(ctx context.Context, uid slot.UserID) ([]slot.Slot, error) {
			if uid == userID {
				return expectedSlots, nil
			}
			return nil, nil
		},
	}

	useCase := New(mockRepo)

	result, err := useCase.ListByUserID(context.Background(), userID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(result) != len(expectedSlots) {
		t.Errorf("expected %d slots, got %d", len(expectedSlots), len(result))
	}

	for i := range result {
		if result[i].ID != expectedSlots[i].ID {
			t.Errorf("slot %d: expected ID %v, got %v", i, expectedSlots[i].ID, result[i].ID)
		}
	}
}

func TestListByUserID_Empty(t *testing.T) {
	mockRepo := &mockSlotRepository{
		listByUserIDFunc: func(ctx context.Context, uid slot.UserID) ([]slot.Slot, error) {
			return []slot.Slot{}, nil
		},
	}

	useCase := New(mockRepo)

	result, err := useCase.ListByUserID(context.Background(), "user-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty list, got %d items", len(result))
	}
}

func TestUpdate_Success(t *testing.T) {
	mockRepo := &mockSlotRepository{
		updateFunc: func(ctx context.Context, s slot.Slot) error {
			if s.SlotName == "Updated Slot" {
				return nil
			}
			return errors.New("invalid slot name")
		},
	}

	useCase := New(mockRepo)

	updatedSlot := slot.Slot{
		ID:             "slot-123",
		UserID:         "user-1",
		SlotName:       "Updated Slot",
		MinCostAdv:     200,
		FormatOfBanner: "banner",
	}

	err := useCase.Update(context.Background(), updatedSlot)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUpdate_Error(t *testing.T) {
	expectedErr := errors.New("update error")
	mockRepo := &mockSlotRepository{
		updateFunc: func(ctx context.Context, s slot.Slot) error {
			return expectedErr
		},
	}

	useCase := New(mockRepo)

	updatedSlot := slot.Slot{
		ID:             "slot-123",
		UserID:         "user-1",
		SlotName:       "Updated Slot",
		MinCostAdv:     200,
		FormatOfBanner: "banner",
	}

	err := useCase.Update(context.Background(), updatedSlot)
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestDelete_Success(t *testing.T) {
	slotID := slot.ID("slot-123")
	userID := slot.UserID("user-1")

	mockRepo := &mockSlotRepository{
		deleteFunc: func(ctx context.Context, id slot.ID, uid slot.UserID) error {
			if id == slotID && uid == userID {
				return nil
			}
			return errors.New("not found")
		},
	}

	useCase := New(mockRepo)

	err := useCase.Delete(context.Background(), slotID, userID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDelete_Error(t *testing.T) {
	expectedErr := errors.New("delete error")
	mockRepo := &mockSlotRepository{
		deleteFunc: func(ctx context.Context, id slot.ID, uid slot.UserID) error {
			return expectedErr
		},
	}

	useCase := New(mockRepo)

	err := useCase.Delete(context.Background(), "slot-123", "user-1")
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestNew(t *testing.T) {
	mockRepo := &mockSlotRepository{}
	useCase := New(mockRepo)

	if useCase == nil {
		t.Error("expected non-nil useCase")
	}

	if useCase.repo == nil {
		t.Error("expected non-nil repo")
	}
}
