package metric

import (
	"2025_2_404/internal/service/slot/domain/metric"
	"context"
)

type metricRepositiryI interface{
	CreateMetric(ctx context.Context, metric metric.Metric) error
}

type MetricUsecase struct{
	repo 	metricRepositiryI
}

func New(repo metricRepositiryI) *MetricUsecase{
	return &MetricUsecase{
		repo: repo,
	}
}

func (u *MetricUsecase) CreateMetric(ctx context.Context, metric metric.Metric) error{
	return u.repo.CreateMetric(ctx, metric)
}