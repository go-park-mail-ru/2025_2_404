// Package ad implements the gRPC server for advertisement-related operations.
package ad

import (
	"2025_2_404/internal/delivery/grpc/interceptor"
	modelad "2025_2_404/internal/service/ad/domain/ad"
	modelfullad "2025_2_404/internal/service/ad/domain/ad_full_info"
	modeluser "2025_2_404/internal/service/ad/domain/user"
	"2025_2_404/pkg/utils"
	adv1 "2025_2_404/protos/gen/go/ad"
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type adUsecaseI interface {
	FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelfullad.AdFullInfo, error)
	Create(ctx context.Context, ad modelad.Ads) error
	Update(ctx context.Context, ad modelad.Ads) error
	Delete(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error
	GetOneAd(ctx context.Context, adID modelad.ID, clientID modeluser.ID) (modelfullad.AdFullInfo, int, error)
	GetAdDetailForSlot(ctx context.Context, id modelad.ID, eventType string) (modelfullad.DetailID, error)
	GetAdSlot(ctx context.Context, minCost uint32) (modelad.Ads, error)
	GetAdCount(ctx context.Context, clientID modeluser.ID) (int64, error)
}

type budgetI interface {
	UpdateBudget(ctx context.Context, adID modelad.ID, clientID modeluser.ID, budget uint32) error
}

type adService struct {
	adUsecase     adUsecaseI
	budgetUsecase budgetI
	adv1.UnimplementedAdServServer
	logger *zap.Logger
}

func New(adUsecase adUsecaseI, budgetUsecase budgetI, logger *zap.Logger) *adService {
	return &adService{
		adUsecase:     adUsecase,
		budgetUsecase: budgetUsecase,
		logger:        logger,
	}
}

func (s *adService) Create(ctx context.Context, req *adv1.CreateRequest) (*adv1.CreateResponse, error) {
	clientID, err := interceptor.GetUserID(ctx)
	if err != nil {
		s.logger.Warn("unauthorized create ad attempt")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	protoAd := req.GetAd()
	s.logger.Debug("received CreateAd request",
		zap.String("client_id", clientID.String()),
		zap.String("title", protoAd.GetTitle()),
		zap.Uint32("budget", protoAd.GetBudget()),
	)

	ad := modelad.Ads{
		Title:     protoAd.Title,
		ClientID:  clientID,
		Content:   protoAd.Content,
		Budget:    protoAd.Budget,
		ImagePath: protoAd.ImgPath,
		TargetURL: protoAd.Targeturl,
	}

	if err := s.adUsecase.Create(ctx, ad); err != nil {
		s.logger.Error("CreateAd usecase failed", zap.Error(err))
		return nil, utils.ToGRPCError(err)
	}

	s.logger.Info("ad created successfully",
		zap.String("client_id", clientID.String()),
		zap.String("title", ad.Title),
	)

	return &adv1.CreateResponse{}, nil
}

func (s *adService) GetAllAds(ctx context.Context, req *adv1.GetAllAdsRequest) (*adv1.GetAllAdsResponse, error) {
	clientID, err := interceptor.GetUserID(ctx)
	if err != nil {
		s.logger.Warn("unauthorized create ad attempt")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	s.logger.Debug("received GetAllAds request",
		zap.String("client_id", clientID.String()),
	)

	adsFull, err := s.adUsecase.FindByUserID(ctx, clientID)
	if err != nil {
		s.logger.Error("GetAllAds usecase failed",
			zap.String("client_id", clientID.String()),
			zap.Error(err),
		)
		return nil, utils.ToGRPCError(err)
	}

	var grpcAds []*adv1.Ad
	for _, a := range adsFull {
		grpcAds = append(grpcAds, &adv1.Ad{
			Id:        uuid.UUID(a.ID).String(),
			Title:     a.Title,
			Status:    a.Status,
			CreatedAt: a.CreatedAt.Format(time.RFC3339),
		})
	}

	s.logger.Info("successfully retrieved ads for user",
		zap.String("client_id", clientID.String()),
		zap.Int("ads_count", len(adsFull)),
	)

	return &adv1.GetAllAdsResponse{Ads: grpcAds}, nil
}

func (s *adService) Update(ctx context.Context, req *adv1.UpdateRequest) (*adv1.UpdateResponse, error) {
	clientID, err := interceptor.GetUserID(ctx)
	if err != nil {
		s.logger.Warn("unauthorized UpdateAd attempt")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	protoAd := req.GetAd()
	if protoAd == nil {
		s.logger.Warn("received UpdateAd request with nil ad", zap.String("client_id", clientID.String()))
		return nil, status.Error(codes.InvalidArgument, "ad is required")
	}
	id, err := uuid.Parse(protoAd.Id)
	if err != nil {
		s.logger.Warn("invalid ad ID in UpdateAd request",
			zap.String("client_id", clientID.String()),
			zap.String("ad_id", protoAd.Id),
		)
		return nil, status.Error(codes.InvalidArgument, "invalid ad ID")
	}

	ad := modelad.Ads{
		ID:        modelad.ID(id),
		ClientID:  clientID,
		Title:     protoAd.Title,
		Content:   protoAd.Content,
		ImagePath: protoAd.ImgPath,
		TargetURL: protoAd.Targeturl,
		Status:    protoAd.Status,
	}

	s.logger.Debug("received UpdateAd request",
		zap.String("client_id", clientID.String()),
		zap.String("ad_id", protoAd.Id),
		zap.String("title", ad.Title),
		zap.String("status", ad.Status),
	)

	if err := s.adUsecase.Update(ctx, ad); err != nil {
		s.logger.Error("UpdateAd usecase failed",
			zap.String("client_id", clientID.String()),
			zap.String("ad_id", protoAd.Id),
			zap.Error(err),
		)
		return nil, utils.ToGRPCError(err)
	}

	s.logger.Info("ad updated successfully",
		zap.String("client_id", clientID.String()),
		zap.String("ad_id", protoAd.Id),
	)

	return &adv1.UpdateResponse{}, nil
}

func (s *adService) Delete(ctx context.Context, req *adv1.DeleteRequest) (*adv1.DeleteResponse, error) {
	clientID, err := interceptor.GetUserID(ctx)
	if err != nil {
		s.logger.Warn("unauthorized DeleteAd attempt")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		s.logger.Warn("received DeleteAd request with empty ad ID",
			zap.String("client_id", clientID.String()),
		)
		return nil, status.Error(codes.InvalidArgument, "invalid ad ID")
	}
	adID := modelad.ID(id)

	if err := s.adUsecase.Delete(ctx, adID, modeluser.ID(clientID)); err != nil {
		s.logger.Error("DeleteAd usecase failed",
			zap.String("client_id", clientID.String()),
			zap.String("ad_id", adID.String()),
			zap.Error(err),
		)
		return nil, utils.ToGRPCError(err)
	}

	s.logger.Info("ad deleted successfully",
		zap.String("client_id", clientID.String()),
		zap.String("ad_id", adID.String()),
	)

	return &adv1.DeleteResponse{}, nil
}

func (s *adService) GetAd(ctx context.Context, req *adv1.GetAdRequest) (*adv1.GetAdResponse, error) {
	clientID, err := interceptor.GetUserID(ctx)
	if err != nil {
		s.logger.Warn("unauthorized GetAd attempt")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		s.logger.Warn("invalid ad ID format in GetAd request",
			zap.String("client_id", clientID.String()),
			zap.String("ad_id", id.String()),
		)
		return nil, status.Error(codes.InvalidArgument, "invalid ad ID")
	}
	adID := modelad.ID(id)

	s.logger.Debug("received GetAd request",
		zap.String("client_id", clientID.String()),
		zap.String("ad_id", id.String()),
	)

	adFull, _, err := s.adUsecase.GetOneAd(ctx, adID, modeluser.ID(clientID))
	if err != nil {
		s.logger.Error("GetAd usecase failed",
			zap.String("client_id", clientID.String()),
			zap.String("ad_id", id.String()),
			zap.Error(err),
		)
		return nil, utils.ToGRPCError(err)
	}

	s.logger.Debug("ad retrieved successfully",
		zap.String("client_id", clientID.String()),
		zap.String("ad_id", id.String()),
	)

	ad := &adv1.Ad{
		Id:          uuid.UUID(adFull.ID).String(),
		Title:       adFull.Title,
		Content:     adFull.Content,
		Targeturl:   adFull.TargetURL,
		ImgPath:     adFull.ImgPath,
		Budget:      adFull.Budget,
		Status:      adFull.Status,
		StartAt:     adFull.StartAt.Format(time.RFC3339),
		EndAt:       adFull.EndAt.Format(time.RFC3339),
		Clicks:      int64(adFull.Clicks),
		Impressions: int64(adFull.Impressions),
	}

	return &adv1.GetAdResponse{Ad: ad}, nil
}

func (s *adService) GetAdDetailForSlot(ctx context.Context, req *adv1.GetAdDetailIDRequest) (*adv1.GetAdDetailIDResponse, error) {
	id, err := uuid.Parse(req.GetAdId())
	if err != nil {
		s.logger.Warn("invalid ad ID format in GetAdDetailForSlot",
			zap.String("ad_id", req.AdId),
			zap.String("event_type", req.EventType),
		)
		return nil, status.Error(codes.InvalidArgument, "invalid ad ID")
	}

	s.logger.Debug("received GetAdDetailForSlot request",
		zap.String("ad_id", req.AdId),
		zap.String("event_type", req.EventType),
	)

	detailID, err := s.adUsecase.GetAdDetailForSlot(ctx, modelad.ID(id), req.GetEventType())
	if err != nil {
		s.logger.Error("GetAdDetailForSlot usecase failed",
			zap.String("ad_id", req.AdId),
			zap.String("event_type", req.EventType),
			zap.Error(err),
		)
		return nil, utils.ToGRPCError(err)
	}

	s.logger.Info("ad detail event processed successfully",
		zap.String("ad_id", id.String()),
		zap.String("event_type", req.EventType),
		zap.String("detail_id", detailID.String()),
	)

	return &adv1.GetAdDetailIDResponse{
		AdDetailId: detailID.String(),
	}, nil
}

func (s *adService) GetAdSlot(ctx context.Context, req *adv1.GetAdSlotRequest) (*adv1.GetAdSlotResponse, error) {
	s.logger.Debug("received GetAdSlot request",
		zap.Uint32("min_cost", req.MinCost),
	)

	adSlot, err := s.adUsecase.GetAdSlot(ctx, req.GetMinCost())
	if err != nil {
		s.logger.Debug("no ad slot found for min cost",
			zap.Uint32("min_cost", req.MinCost),
			zap.Error(err),
		)
		return &adv1.GetAdSlotResponse{
			Ad: &adv1.AdSlot{},
		}, nil
	}

	s.logger.Debug("ad slot served successfully",
		zap.String("ad_id", adSlot.ID.String()),
		zap.Uint32("min_cost", req.MinCost),
	)

	adRes := &adv1.AdSlot{
		Id:          adSlot.ID.String(),
		Title:       adSlot.Title,
		Description: adSlot.Content,
		ImageSrc:    adSlot.ImagePath,
		Link:        adSlot.TargetURL,
	}

	return &adv1.GetAdSlotResponse{Ad: adRes}, nil
}

func (s *adService) UpdateAdBudget(ctx context.Context, req *adv1.UpdateBudgetRequest) (*adv1.UpdateBudgetResponse, error) {
	clientID, err := interceptor.GetUserID(ctx)
	if err != nil {
		s.logger.Warn("unauthorized UpdateAdBudget attempt")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	adIDStr := req.GetId()
	if adIDStr == "" {
		s.logger.Warn("received UpdateAdBudget with empty ad ID",
			zap.String("client_id", clientID.String()),
		)
		return nil, status.Error(codes.InvalidArgument, "ad id is required")
	}

	id, err := uuid.Parse(adIDStr)
	if err != nil {
		s.logger.Warn("invalid ad ID format in UpdateAdBudget",
			zap.String("client_id", clientID.String()),
			zap.String("ad_id", adIDStr),
		)
		return nil, status.Error(codes.InvalidArgument, "invalid ad ID format")
	}

	newBudget := req.GetBudget()
	s.logger.Debug("received UpdateAdBudget request",
		zap.String("client_id", clientID.String()),
		zap.String("ad_id", adIDStr),
		zap.Uint32("new_budget", newBudget),
	)
	if err := s.budgetUsecase.UpdateBudget(ctx, modelad.ID(id), clientID, newBudget); err != nil {
		s.logger.Error("UpdateAdBudget usecase failed",
			zap.String("client_id", clientID.String()),
			zap.String("ad_id", adIDStr),
			zap.Uint32("new_budget", newBudget),
			zap.Error(err),
		)
		return nil, utils.ToGRPCError(err)
	}

	s.logger.Info("ad budget updated successfully",
		zap.String("client_id", clientID.String()),
		zap.String("ad_id", adIDStr),
		zap.Uint32("new_budget", newBudget),
	)

	return &adv1.UpdateBudgetResponse{
		Budget: newBudget,
	}, nil
}

func (s *adService) GetAdCount(ctx context.Context, req *adv1.GetAdCountRequest) (*adv1.GetAdCountResponse, error) {
	clientID, err := interceptor.GetUserID(ctx)
	if err != nil {
		s.logger.Warn("unauthorized GetAdCount attempt")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	s.logger.Debug("received GetAdCount request",
		zap.String("client_id", clientID.String()),
	)

	count, err := s.adUsecase.GetAdCount(ctx, clientID)
	if err != nil {
		s.logger.Error("GetAdCount usecase failed",
			zap.String("client_id", clientID.String()),
			zap.Error(err),
		)
		return nil, utils.ToGRPCError(err)
	}

	s.logger.Debug("ad count retrieved successfully",
		zap.String("client_id", clientID.String()),
		zap.Int64("count", count),
	)

	return &adv1.GetAdCountResponse{
		Count: count,
	}, nil
}
