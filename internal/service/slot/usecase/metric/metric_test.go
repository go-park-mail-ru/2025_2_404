package metric

import (
	user "2025_2_404/internal/service/profile/domain"
	"2025_2_404/internal/service/slot/domain/metric"
	"context"
	"errors"
	"testing"
)

type mockMetricRepository struct {
	createMetricFunc   func(ctx context.Context, m metric.Metric) (user.ID, error)
	getMetricForDayFunc func(ctx context.Context, slotID metric.SlotID) ([]metric.GetMetric, error)
}

func (m *mockMetricRepository) CreateMetric(ctx context.Context, met metric.Metric) (user.ID, error) {
	if m.createMetricFunc != nil {
		return m.createMetricFunc(ctx, met)
	}
	return user.ID{}, nil
}

func (m *mockMetricRepository) GetMetricForDay(ctx context.Context, slotID metric.SlotID) ([]metric.GetMetric, error) {
	if m.getMetricForDayFunc != nil {
		return m.getMetricForDayFunc(ctx, slotID)
	}
	return nil, nil
}

func TestNew(t *testing.T) {
	mockRepo := &mockMetricRepository{}
	useCase := New(mockRepo)

	if useCase == nil {
		t.Fatal("expected useCase to be created, got nil")
	}

	if useCase.repo == nil {
		t.Error("expected repo to be set")
	}
}

func TestCreateMetric_Success(t *testing.T) {
	called := false
	expectedID := user.ID{}

	mockRepo := &mockMetricRepository{
		createMetricFunc: func(ctx context.Context, m metric.Metric) (user.ID, error) {
			called = true
			return expectedID, nil
		},
	}

	useCase := New(mockRepo)
	met := metric.Metric{}

	id, err := useCase.CreateMetric(context.Background(), met)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if id != expectedID {
		t.Errorf("expected ID %v, got %v", expectedID, id)
	}

	if !called {
		t.Error("expected repository CreateMetric to be called")
	}
}

func TestCreateMetric_RepositoryError(t *testing.T) {
	expectedError := errors.New("database error")

	mockRepo := &mockMetricRepository{
		createMetricFunc: func(ctx context.Context, m metric.Metric) (user.ID, error) {
			return user.ID{}, expectedError
		},
	}

	useCase := New(mockRepo)
	met := metric.Metric{}

	_, err := useCase.CreateMetric(context.Background(), met)

	if err != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}

func TestGetMetricForSlot_Success(t *testing.T) {
	mockMetrics := []metric.GetMetric{
		{Clicks: 10, Impressions: 100},
		{Clicks: 20, Impressions: 200},
		{Clicks: 30, Impressions: 300},
	}

	mockRepo := &mockMetricRepository{
		getMetricForDayFunc: func(ctx context.Context, slotID metric.SlotID) ([]metric.GetMetric, error) {
			return mockMetrics, nil
		},
	}

	useCase := New(mockRepo)
	var slotID metric.SlotID

	clicks, impressions, metrics, err := useCase.GetMetricForSlot(context.Background(), slotID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expectedClicks := 60 // 10 + 20 + 30
	if clicks != expectedClicks {
		t.Errorf("expected totalClicks %d, got %d", expectedClicks, clicks)
	}

	expectedImpressions := 600 // 100 + 200 + 300
	if impressions != expectedImpressions {
		t.Errorf("expected totalImpressions %d, got %d", expectedImpressions, impressions)
	}

	if len(metrics) != 3 {
		t.Errorf("expected 3 metrics, got %d", len(metrics))
	}
}

func TestGetMetricForSlot_EmptyMetrics(t *testing.T) {
	mockRepo := &mockMetricRepository{
		getMetricForDayFunc: func(ctx context.Context, slotID metric.SlotID) ([]metric.GetMetric, error) {
			return []metric.GetMetric{}, nil
		},
	}

	useCase := New(mockRepo)
	var slotID metric.SlotID

	clicks, impressions, metrics, err := useCase.GetMetricForSlot(context.Background(), slotID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if clicks != 0 {
		t.Errorf("expected totalClicks 0, got %d", clicks)
	}

	if impressions != 0 {
		t.Errorf("expected totalImpressions 0, got %d", impressions)
	}

	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics, got %d", len(metrics))
	}
}

func TestGetMetricForSlot_RepositoryError(t *testing.T) {
	expectedError := errors.New("database connection failed")

	mockRepo := &mockMetricRepository{
		getMetricForDayFunc: func(ctx context.Context, slotID metric.SlotID) ([]metric.GetMetric, error) {
			return nil, expectedError
		},
	}

	useCase := New(mockRepo)
	var slotID metric.SlotID

	clicks, impressions, metrics, err := useCase.GetMetricForSlot(context.Background(), slotID)

	if err != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}

	if clicks != 0 {
		t.Errorf("expected totalClicks 0, got %d", clicks)
	}

	if impressions != 0 {
		t.Errorf("expected totalImpressions 0, got %d", impressions)
	}

	if metrics != nil {
		t.Errorf("expected nil metrics, got %v", metrics)
	}
}
