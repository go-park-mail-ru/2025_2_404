package slot

import (
	"2025_2_404/internal/delivery/grpc/interceptor"
	"2025_2_404/internal/service/slot/domain/slot"
	slotpb "2025_2_404/protos/gen/go/slot"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type slotUsecase interface {
	Create(ctx context.Context, s slot.Slot) (slot.ID, error) 
	GetByID(ctx context.Context, id slot.ID) (slot.Slot, error)
	ListByUserID(ctx context.Context, userID slot.UserID) ([]slot.Slot, error)
	Update(ctx context.Context, s slot.Slot) error
	Delete(ctx context.Context, id slot.ID, userID slot.UserID) error
}

type slotService struct {
	slotUsecase slotUsecase
	slotpb.UnimplementedSlotServServer
}

func New(u slotUsecase) *slotService {
	return &slotService{slotUsecase: u}
}

func (s *slotService) CreateSlot(ctx context.Context, req *slotpb.CreateSlotRequest) (*slotpb.CreateSlotResponse, error) {
	userID, err := interceptor.GetUserID(ctx)
	if err != nil {
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

	slotID, err := s.slotUsecase.Create(ctx, domainSlot)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create slot: %v", err)
	}

	return &slotpb.CreateSlotResponse{Id: string(slotID)}, nil
}

func (s *slotService) GetSlot(ctx context.Context, req *slotpb.GetSlotRequest) (*slotpb.GetSlotResponse, error) {
	id := req.GetId()

	slot, err := s.slotUsecase.GetByID(ctx, slot.ID(id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "slot not found: %v", err)
	}

	return &slotpb.GetSlotResponse{
		Slot: &slotpb.Slot{
			Id:              string(slot.ID),
			UserId:          string(slot.UserID),
			SlotName:        slot.SlotName,
			MinCostAdv:      slot.MinCostAdv,
			FormatOfBanner:  slot.FormatOfBanner,
			Status:          slot.Status,
			BackColor:       slot.BackColor,
			TextColor:       slot.TextColor,
		},
	}, nil
}

func (s *slotService) ListSlots(ctx context.Context, req *slotpb.ListSlotsRequest) (*slotpb.ListSlotsResponse, error) {
	userID, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	slots, err := s.slotUsecase.ListByUserID(ctx, slot.UserID(userID.String()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list slots: %v", err)
	}

	var grpcSlots []*slotpb.Slot
	for _, s := range slots {
		grpcSlots = append(grpcSlots, &slotpb.Slot{
			Id:              string(s.ID),
			UserId:          string(s.UserID),
			SlotName:        s.SlotName,
			MinCostAdv:      s.MinCostAdv,
			FormatOfBanner:  s.FormatOfBanner,
			Status:          s.Status,
			BackColor:       s.BackColor,
			TextColor:       s.TextColor,
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
		return nil, status.Errorf(codes.Internal, "update slot: %v", err)
	}

	return &slotpb.UpdateSlotResponse{}, nil
}

func (s *slotService) DeleteSlot(ctx context.Context, req *slotpb.DeleteSlotRequest) (*slotpb.DeleteSlotResponse, error) {
	userID, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	id := req.GetId()

	if err := s.slotUsecase.Delete(ctx, slot.ID(id), slot.UserID(userID.String())); err != nil {
		return nil, status.Errorf(codes.Internal, "delete slot: %v", err)
	}

	return &slotpb.DeleteSlotResponse{}, nil
}