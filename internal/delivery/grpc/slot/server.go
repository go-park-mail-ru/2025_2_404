// Package slot provides gRPC delivery handlers for slot-related operations.
package slot

import (
	"2025_2_404/internal/delivery/grpc/interceptor"
	user "2025_2_404/internal/service/profile/domain"
	"2025_2_404/internal/service/slot/domain/metric"
	"2025_2_404/internal/service/slot/domain/slot"
	"2025_2_404/pkg/utils"
	adpb "2025_2_404/protos/gen/go/ad"
	profilepb "2025_2_404/protos/gen/go/profile"
	slotpb "2025_2_404/protos/gen/go/slot"
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type slotUsecaseI interface {
	Create(ctx context.Context, s slot.Slot) (slot.ID, error)
	GetByID(ctx context.Context, id slot.ID) (slot.Slot, error)
	ListByUserID(ctx context.Context, userID slot.UserID) ([]slot.Slot, error)
	Update(ctx context.Context, s slot.Slot) error
	Delete(ctx context.Context, id slot.ID, userID slot.UserID) error
}

type metricUsecaseI interface {
	CreateMetric(ctx context.Context, metric metric.Metric) (user.ID, error) //лучшее что я придумал добавить вывод из метрик юзера без транзакции
	GetMetricForSlot(ctx context.Context, slotID metric.SlotID) (int, int, []metric.GetMetric, error)
}

type slotService struct {
	slotUsecase   slotUsecaseI
	metricUsecase metricUsecaseI
	slotpb.UnimplementedSlotServServer
	clientAD      adpb.AdServClient
	clientProfile profilepb.ProfileClient
	logger        *slog.Logger
}

func New(u slotUsecaseI, metricUsecase metricUsecaseI, clientAD adpb.AdServClient, clientProfile profilepb.ProfileClient) *slotService {
	return &slotService{
		slotUsecase:   u,
		metricUsecase: metricUsecase,
		clientAD:      clientAD,
		clientProfile: clientProfile,
		logger:        slog.Default(),
	}
}

func (s *slotService) CreateSlot(ctx context.Context, req *slotpb.CreateSlotRequest) (*slotpb.CreateSlotResponse, error) {
	s.logger.Debug("CreateSlot called", "req", req)

	userID, err := interceptor.GetUserID(ctx)
	if err != nil {
		s.logger.Warn("Unauthorized CreateSlot", "error", err)
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	pbSlot := req.GetSlot()
	domainSlot := slot.Slot{
		UserID:         slot.UserID(userID.String()),
		SlotName:       pbSlot.GetSlotName(),
		MinCostAdv:     pbSlot.GetMinCostAdv(),
		FormatOfBanner: pbSlot.GetFormatOfBanner(),
		Status:         pbSlot.GetStatus(),
		BackColor:      pbSlot.GetBackColor(),
		TextColor:      pbSlot.GetTextColor(),
	}

	s.logger.Debug("Creating slot in usecase", "slot", domainSlot)

	slotID, err := s.slotUsecase.Create(ctx, domainSlot)
	if err != nil {
		s.logger.Error("Failed to create slot", "error", err)
		return nil, utils.ToGRPCError(err)
	}

	s.logger.Info("Slot created successfully", "slot_id", slotID)
	return &slotpb.CreateSlotResponse{Id: string(slotID)}, nil
}

func (s *slotService) GetSlot(ctx context.Context, req *slotpb.GetSlotRequest) (*slotpb.GetSlotResponse, error) {
	id := req.GetId()
	s.logger.Debug("GetSlot called", "slot_id", id)

	slotDomain, err := s.slotUsecase.GetByID(ctx, slot.ID(id))
	if err != nil {
		s.logger.Warn("Slot not found", "slot_id", id, "error", err)
		return nil, utils.ToGRPCError(err)
	}

	adSlot, err := s.clientAD.GetAdSlot(ctx, &adpb.GetAdSlotRequest{
		ClientId: string(slotDomain.UserID),
		MinCost:  uint32(slotDomain.MinCostAdv),
	})
	if err != nil {
		s.logger.Warn("Bad request in adservice", "error", err)
		return nil, utils.ToGRPCError(err)
	}

	s.logger.Debug("Fetched slot data", "slot", slotDomain, "adSlot", adSlot)

	resp := &slotpb.GetSlotResponse{
		Slot: &slotpb.Slot{
			Id:             string(slotDomain.ID),
			UserId:         string(slotDomain.UserID),
			SlotName:       slotDomain.SlotName,
			MinCostAdv:     slotDomain.MinCostAdv,
			FormatOfBanner: slotDomain.FormatOfBanner,
			Status:         slotDomain.Status,
			BackColor:      slotDomain.BackColor,
			TextColor:      slotDomain.TextColor,
		},
		AdSlot: &slotpb.AdSlot{
			Id:          adSlot.Ad.GetId(),
			Title:       adSlot.Ad.GetTitle(),
			Description: adSlot.Ad.GetDescription(),
			ImageSrc:    adSlot.Ad.GetImageSrc(),
			Link:        adSlot.Ad.GetLink(),
		},
	}

	s.logger.Debug("GetSlot response prepared", "response", resp)
	return resp, nil
}

func (s *slotService) ListSlots(ctx context.Context, req *slotpb.ListSlotsRequest) (*slotpb.ListSlotsResponse, error) {
	userID, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	slots, err := s.slotUsecase.ListByUserID(ctx, slot.UserID(userID.String()))
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	var grpcSlots []*slotpb.Slot
	for _, s := range slots {
		grpcSlots = append(grpcSlots, &slotpb.Slot{
			Id:             string(s.ID),
			UserId:         string(s.UserID),
			SlotName:       s.SlotName,
			MinCostAdv:     s.MinCostAdv,
			FormatOfBanner: s.FormatOfBanner,
			Status:         s.Status,
			BackColor:      s.BackColor,
			TextColor:      s.TextColor,
		})
	}

	return &slotpb.ListSlotsResponse{Slots: grpcSlots}, nil
}

func (s *slotService) UpdateSlot(ctx context.Context, req *slotpb.UpdateSlotRequest) (*slotpb.UpdateSlotResponse, error) {
	userID, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	pbSlot := req.GetSlot()
	domainSlot := slot.Slot{
		ID:             slot.ID(pbSlot.Id),
		UserID:         slot.UserID(userID.String()),
		SlotName:       pbSlot.SlotName,
		MinCostAdv:     pbSlot.MinCostAdv,
		FormatOfBanner: pbSlot.FormatOfBanner,
		Status:         pbSlot.Status,
		BackColor:      pbSlot.BackColor,
		TextColor:      pbSlot.TextColor,
	}

	if err := s.slotUsecase.Update(ctx, domainSlot); err != nil {
		return nil, utils.ToGRPCError(err)
	}

	return &slotpb.UpdateSlotResponse{}, nil
}

func (s *slotService) DeleteSlot(ctx context.Context, req *slotpb.DeleteSlotRequest) (*slotpb.DeleteSlotResponse, error) {
	userID, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	id := req.GetId()

	if err := s.slotUsecase.Delete(ctx, slot.ID(id), slot.UserID(userID.String())); err != nil {
		return nil, utils.ToGRPCError(err)
	}

	return &slotpb.DeleteSlotResponse{}, nil
}

func (s *slotService) CreateMetric(ctx context.Context, req *slotpb.CreateMetricRequest) (*slotpb.CreateMetricResponse, error) {
	s.logger.Debug("CreateMetric called", "req", req)

	slotID, err := uuid.Parse(req.SlotId)
	if err != nil {
		s.logger.Warn("Invalid slot_id in CreateMetric", "slot_id", req.SlotId, "error", err)
		return nil, utils.ToGRPCError(err)
	}

	adID, err := uuid.Parse(req.AdId)
	if err != nil {
		s.logger.Warn("Invalid ad_id in CreateMetric", "ad_id", req.AdId, "error", err)
		return nil, utils.ToGRPCError(err)
	}

	// Запрос в ad-service
	reqAd := &adpb.GetAdDetailIDRequest{
		AdId: adID.String(),
	}
	s.logger.Debug("Calling ad-service.GetAdDetailForSlot", "ad_id", adID.String())

	resAd, err := s.clientAD.GetAdDetailForSlot(ctx, reqAd)
	if err != nil {
		s.logger.Error("ad-service.GetAdDetailForSlot failed", "error", err)
		return nil, utils.ToGRPCError(err)
	}

	adDetailID, err := uuid.Parse(resAd.GetAdDetailId())
	if err != nil {
		s.logger.Error("Invalid ad_detail_id from ad-service", "ad_detail_id", resAd.GetAdDetailId(), "error", err)
		return nil, status.Error(codes.DataLoss, "adDetailId not correctly formatted")
	}

	metric := metric.Metric{
		SlotID:     metric.SlotID(slotID),
		AdDetailID: metric.AdDetailID(adDetailID),
		EventType:  req.GetEventType(),
	}

	s.logger.Debug("Creating metric", "metric", metric)
	clientID, err := s.metricUsecase.CreateMetric(ctx, metric)
	if err != nil {
		s.logger.Error("Failed to store metric", "error", err)
		return nil, utils.ToGRPCError(err)
	}

	_, err = s.clientProfile.AddBalance(ctx, &profilepb.AddBalanceRequest{
		ClientId:  clientID.String(),
		AddAmount: 2,
	})

	if err != nil {
		return &slotpb.CreateMetricResponse{}, utils.ToGRPCError(err)
	}

	s.logger.Info("Metric recorded successfully", "slot_id", slotID, "ad_detail_id", adDetailID, "event", req.GetEventType())
	return &slotpb.CreateMetricResponse{}, nil
}

func (s *slotService) GetMetrics(ctx context.Context, req *slotpb.GetMetricsRequest) (*slotpb.GetMetricsResponse, error) {
	s.logger.Debug("GetMetrics called", "req", req)

	slotID, err := uuid.Parse(req.GetSlotId())
	if err != nil {
		s.logger.Warn("Invalid slot_id in GetMetrics", "slot_id", req.SlotId, "error", err)
		return nil, status.Error(codes.InvalidArgument, "bad slot id")
	}

	totalClicks, totalImpressions, metrics, err := s.metricUsecase.GetMetricForSlot(ctx, metric.SlotID(slotID))
	if err != nil {
		s.logger.Error("Failed to get metric", "error", err)
		return nil, utils.ToGRPCError(err)
	}

	var grpcMetrics []*slotpb.MetricsForDay

	for _, metric := range metrics {
		grpcMetrics = append(grpcMetrics, &slotpb.MetricsForDay{
			SlotId:      slotID.String(),
			Impressions: int32(metric.Impressions),
			Clicks:      int32(metric.Clicks),
			EventData:   metric.EventDate.GoString(),
		})
	}

	return &slotpb.GetMetricsResponse{
		SlotId:           slotID.String(),
		TotalClicks:      int32(totalClicks),
		TotalImpressions: int32(totalImpressions),
		Metrics:          grpcMetrics,
	}, nil
}
