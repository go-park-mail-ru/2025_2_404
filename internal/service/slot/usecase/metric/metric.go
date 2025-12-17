// Package metric provides use case implementations for managing slot metrics.
package metric

import (
	user "2025_2_404/internal/service/profile/domain"
	"2025_2_404/internal/service/slot/domain/metric"
	"context"
)

type metricRepositiryI interface {
	CreateMetric(ctx context.Context, metric metric.Metric) (user.ID, error)
	GetMetricForDay(ctx context.Context, slotID metric.SlotID) ([]metric.GetMetric, error)
}

type MetricUsecase struct {
	repo metricRepositiryI
}

func New(repo metricRepositiryI) *MetricUsecase {
	return &MetricUsecase{
		repo: repo,
	}
}

func (u *MetricUsecase) CreateMetric(ctx context.Context, metric metric.Metric) (user.ID, error) {
	return u.repo.CreateMetric(ctx, metric)
}

func (u *MetricUsecase) GetMetricForSlot(ctx context.Context, slotID metric.SlotID) (int, int, []metric.GetMetric, error) {
	metrics, err := u.repo.GetMetricForDay(ctx, slotID)
	if err != nil {
		return 0, 0, nil, err
	}

	totalClicks := 0
	totalImpressions := 0
	for _, m := range metrics {
		totalClicks += m.Clicks
		totalImpressions += m.Impressions
	}

	return totalClicks, totalImpressions, metrics, nil
}
