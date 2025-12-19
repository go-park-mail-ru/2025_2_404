// Package storage provides gRPC delivery handlers for storage-related operations.
package storage

import (
	"2025_2_404/pkg/utils"
	storagev1 "2025_2_404/protos/gen/go/storage"
	"context"
	"log/slog"
)

type storageUsecaseI interface {
	Create(ctx context.Context, imageData []byte, imagePath string) error
	Delete(ctx context.Context, imagePath string) error
	Get(ctx context.Context, imagePath string) ([]byte, string, error)
}

type storageService struct {
	storageUsecase storageUsecaseI
	storagev1.UnimplementedStorageServer
}

func New(storageUsecase storageUsecaseI) *storageService {
	return &storageService{
		storageUsecase: storageUsecase,
	}
}

func (s *storageService) Create(ctx context.Context, req *storagev1.CreateRequest) (*storagev1.CreateResponse, error) {
	slog.Debug("📂 Create: received request", "image_path", req.ImagePath, "size", len(req.ImageData))
	err := s.storageUsecase.Create(ctx, req.ImageData, req.ImagePath)
	if err != nil {
		slog.Error("❌ Create failed", "image_path", req.ImagePath, "error", err)
		return nil, utils.ToGRPCError(err)
	}

	slog.Info("✅ Create succeeded", "image_path", req.ImagePath)
	return &storagev1.CreateResponse{}, nil
}

func (s *storageService) Delete(ctx context.Context, req *storagev1.DeleteRequest) (*storagev1.DeleteResponse, error) {
	slog.Debug("🗑️ Delete: received request", "image_path", req.ImagePath)
	err := s.storageUsecase.Delete(ctx, req.ImagePath)
	if err != nil {
		slog.Error("❌ Delete failed", "image_path", req.ImagePath, "error", err)
		return nil, utils.ToGRPCError(err)
	}

	slog.Info("✅ Delete succeeded", "image_path", req.ImagePath)
	return &storagev1.DeleteResponse{
		Success: true,
	}, nil
}

func (s *storageService) Get(ctx context.Context, req *storagev1.GetRequest) (*storagev1.GetResponse, error) {
	slog.Debug("📥 Get: received request", "image_path", req.ImagePath)
	imageData, contentType, err := s.storageUsecase.Get(ctx, req.ImagePath)
	if err != nil {
		slog.Error("❌ Get failed", "image_path", req.ImagePath, "error", err)
		return nil, utils.ToGRPCError(err)
	}

	slog.Info("✅ Get succeeded", "image_path", req.ImagePath, "size", len(imageData), "content_type", contentType)
	return &storagev1.GetResponse{
		ImageData:   imageData,
		ContentType: contentType,
	}, nil
}
