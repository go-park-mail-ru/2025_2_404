package storage

import (
    "context"
    storagev1 "2025_2_404/protos/gen/go/storage"
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
    err := s.storageUsecase.Create(ctx, req.ImageData, req.ImagePath)
    if err != nil {
        return nil, err
    }

    return &storagev1.CreateResponse{}, nil
}

func (s *storageService) Delete(ctx context.Context, req *storagev1.DeleteRequest) (*storagev1.DeleteResponse, error) {
    if err := s.storageUsecase.Delete(ctx, req.ImagePath); err != nil {
        return nil, err
    }

    return &storagev1.DeleteResponse{
        Success: true,
    }, nil
}

func (s *storageService) Get(ctx context.Context, req *storagev1.GetRequest) (*storagev1.GetResponse, error) {
    imageData, contentType, err := s.storageUsecase.Get(ctx, req.ImagePath)
    if err != nil {
        return nil, err
    }

    return &storagev1.GetResponse{
        ImageData:   imageData,
        ContentType: contentType,
    }, nil
}