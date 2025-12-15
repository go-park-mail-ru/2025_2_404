package metric

import (
	user "2025_2_404/internal/service/profile/domain"
	"2025_2_404/internal/service/slot/domain/metric"
	"context"
	"fmt"
)

type metricRepositiryI interface{
	CreateMetric(ctx context.Context, metric metric.Metric) (user.ID ,error)
	GetMetricForDay(ctx context.Context, slotID metric.SlotID)  ([]metric.GetMetric, error)
}

type MetricUsecase struct{
	repo 	metricRepositiryI
}

func New(repo metricRepositiryI) *MetricUsecase{
	return &MetricUsecase{
		repo: repo,
	}
}

func (u *MetricUsecase) CreateMetric(ctx context.Context, metric metric.Metric) (user.ID ,error){
	return u.repo.CreateMetric(ctx, metric)
}

func (u *MetricUsecase) GetMetricForSlot(ctx context.Context, slotID metric.SlotID) (int, int, []metric.GetMetric, error){

	metrics, err := u.repo.GetMetricForDay(ctx, slotID)
	if err != nil{
		return 0, 0, nil, fmt.Errorf("failed work with sql:%w", err)
	}

	total_clicks := 0
	total_impressions := 0
	for i := 0; i < len(metrics); i ++{
		total_clicks = total_clicks + metrics[i].Clicks
		total_impressions = total_impressions + metrics[i].Impressions
	}
	
	return  total_clicks, total_impressions, metrics, nil
}