// Package handler provides gRPC delivery handlers for profile-related operations.
package handler

import (
	"2025_2_404/internal/delivery/grpc/interceptor"
	modeluser "2025_2_404/internal/service/profile/domain"
	"2025_2_404/protos/gen/go/profile"
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileUsecaseI interface {
	Update(ctx context.Context, client modeluser.User) error
	Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, error)
	Delete(ctx context.Context, clientID modeluser.ID) error
	ShowBalance(ctx context.Context, clientID modeluser.ID) (uint32, error)
	AddBalance(ctx context.Context, clientID modeluser.ID, addAmount uint32) error
	SubtractBalance(ctx context.Context, clientID modeluser.ID, subAmount uint32) error
	CreatePayment(ctx context.Context, payment modeluser.Payment) (string, error)
	UpdatePaymentStatus(ctx context.Context, yooPaymentID string, status modeluser.PaymentStatus) (modeluser.ID, error)
	GetPaymentsByClientID(ctx context.Context, clientID modeluser.ID) ([]modeluser.Payment, error)
}

type ProfileServer struct {
	profile.UnimplementedProfileServer
	profileUsecase ProfileUsecaseI
}

func NewProfileServer(profileUsecase ProfileUsecaseI) *ProfileServer {
	return &ProfileServer{
		profileUsecase: profileUsecase,
	}
}

func (h *ProfileServer) Update(ctx context.Context, req *profile.UpdateRequest) (*profile.UpdateResponse, error) {
	clientID, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}
	fmt.Printf("DEBUG HANDLER: Получен ID из контекста: %v\n", clientID)
	client := modeluser.User{
		ID:            clientID,
		Email:         req.GetEmail(),
		UserName:      req.GetUserName(),
		UserFirstName: req.GetFirstName(),
		UserLastName:  req.GetLastName(),
		Company:       req.GetCompany(),
		Phone:         req.GetPhone(),
		ImagePath:     req.GetAvatarPath(),
	}

	err = h.profileUsecase.Update(ctx, client)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update profile: %v", err)
	}

	return &profile.UpdateResponse{
		UserName:   client.UserName,
		Email:      client.Email,
		FirstName:  client.UserFirstName,
		LastName:   client.UserLastName,
		Company:    client.Company,
		Phone:      client.Phone,
		AvatarPath: client.ImagePath,
	}, nil
}

func (h *ProfileServer) Show(ctx context.Context, req *profile.ShowRequest) (*profile.ShowResponse, error) {
	clientID, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	user, err := h.profileUsecase.Show(ctx, clientID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to show profile: %v", err)
	}

	return &profile.ShowResponse{
		UserName:   user.UserName,
		Email:      user.Email,
		FirstName:  user.UserFirstName,
		LastName:   user.UserLastName,
		Company:    user.Company,
		Phone:      user.Phone,
		AvatarPath: user.ImagePath,
		CreatedAt:  user.CreatedAt,
	}, nil
}

func (h *ProfileServer) Delete(ctx context.Context, req *profile.DeleteRequest) (*profile.DeleteResponse, error) {
	clientID, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = h.profileUsecase.Delete(ctx, clientID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete profile: %v", err)
	}

	return &profile.DeleteResponse{}, nil
}

func (h *ProfileServer) ShowBalance(ctx context.Context, req *profile.ShowBalanceRequest) (*profile.ShowBalanceResponse, error) {
	clientID, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	balance, err := h.profileUsecase.ShowBalance(ctx, clientID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to show balance: %v", err)
	}

	return &profile.ShowBalanceResponse{
		Balance: balance,
	}, nil
}

func (h *ProfileServer) AddBalance(ctx context.Context, req *profile.AddBalanceRequest) (*profile.AddBalanceResponse, error) {

	clientID, err := uuid.Parse(req.GetClientId())
	if err != nil {
		return nil, fmt.Errorf("clientId not correct")
	}

	err = h.profileUsecase.AddBalance(ctx, clientID, req.GetAddAmount())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add balance: %v", err)
	}

	return &profile.AddBalanceResponse{}, nil
}

func (h *ProfileServer) SubtractBalance(ctx context.Context, req *profile.SubtractBalanceRequest) (*profile.SubtractBalanceResponse, error) {
	clientID, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = h.profileUsecase.SubtractBalance(ctx, clientID, req.GetSubAmount())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to subtract balance: %v", err)
	}

	return &profile.SubtractBalanceResponse{}, nil
}

func (h *ProfileServer) GetPaymentsByClientID(ctx context.Context, req *profile.PaymentsByClientIDRequest) (*profile.PaymentsByClientIDResponse, error) {
	clientID, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	payments, err := h.profileUsecase.GetPaymentsByClientID(ctx, clientID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get payments: %v", err)
	}

	var pbPayments []*profile.Payment
	for _, payment := range payments {
		pbPayment := &profile.Payment{
			Id:            payment.ID.String(),
			Amount:        uint32(payment.AmountRub),
			Status:        string(payment.Status),
			YooPaymentId:  payment.YooPaymentID,
			MethodPayment: payment.PaymentMethod,
			CreatedAt:     payment.CreatedTime,
		}
		pbPayments = append(pbPayments, pbPayment)
	}

	return &profile.PaymentsByClientIDResponse{
		Payments: pbPayments,
	}, nil
}

func (h *ProfileServer) CreatePayment(ctx context.Context, req *profile.PaymentCreateRequest) (*profile.PaymentCreateResponse, error) {
	clientID, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	payment := modeluser.Payment{
		ClientID:      clientID,
		PaymentMethod: req.GetPaymentMethod(),
		AmountRub:     req.GetAmount(),
	}

	yooKassaLink, err := h.profileUsecase.CreatePayment(ctx, payment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment: %v", err)
	}

	return &profile.PaymentCreateResponse{
		PaymentUrl: yooKassaLink,
	}, nil
}

func (h *ProfileServer) UpdatePaymentStatus(ctx context.Context, req *profile.PaymentStatusRequest) (*profile.PaymentStatusResponse, error) {
	slog.Info("🔄 UpdatePaymentStatus called",
		"yookassa_id", req.GetYookassaId(),
		"status", req.GetStatus(),
		"amount", req.GetAmount(),
	)

	clientID, err := h.profileUsecase.UpdatePaymentStatus(ctx, req.GetYookassaId(), modeluser.PaymentStatus(req.GetStatus()))
	if err != nil {
		slog.Error("Failed to update payment status in DB",
			"yookassa_id", req.GetYookassaId(),
			"status", req.GetStatus(),
			"error", err,
		)
		return nil, status.Errorf(codes.Internal, "failed to update payment status: %v", err)
	}

	slog.Info("Payment status updated, client identified",
		"yookassa_id", req.GetYookassaId(),
		"client_id", clientID.String(),
	)

	amount, err := strconv.ParseUint(req.GetAmount(), 10, 32)
	if err != nil {
		slog.Error("Invalid amount format",
			"yookassa_id", req.GetYookassaId(),
			"amount", req.GetAmount(),
			"error", err,
		)
		return nil, status.Errorf(codes.InvalidArgument, "failed to convert amount: %v", err)
	}

	slog.Info("Adding balance",
		"client_id", clientID.String(),
		"add_amount_rub", amount,
	)
	if req.Status == string(modeluser.PaymentSucceeded) {
		_, err = h.AddBalance(ctx, &profile.AddBalanceRequest{
			ClientId:  clientID.String(),
			AddAmount: uint32(amount),
		})
		if err != nil {
			slog.Error("💸 Failed to add balance",
				"client_id", clientID.String(),
				"amount", amount,
				"error", err,
			)
		}
	}

	slog.Info("UpdatePaymentStatus completed successfully",
		"yookassa_id", req.GetYookassaId(),
		"client_id", clientID.String(),
		"amount_rub", amount,
	)

	return &profile.PaymentStatusResponse{}, nil
}
