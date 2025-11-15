package support

import (
	modelsup "2025_2_404/internal/domain/models/support"
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/google/uuid"
)

type repositoryI interface {
	Create(ctx context.Context, sup modelsup.Support) (error)
	GetOneSup(ctx context.Context, supID modelsup.ID, clientId modeluser.ID) (modelsup.Support, error)
	Update(ctx context.Context, sup modelsup.Support) error
	GetAllSups(ctx context.Context, superID int64) ([]modelsup.Support, error)
}

type fileStorageI interface {
	Save(uploadPath string, data io.Reader) error
	ReadImageAsBytes(path string) ([]byte, error)
	Delete(path string) error
}

type UseCase struct {
	supRepo repositoryI
	fileStorage fileStorageI
}

func New(repo repositoryI, fileStorage fileStorageI) *UseCase {
	return &UseCase{
		supRepo: repo,
		fileStorage: fileStorage,
	}
}

func (u *UseCase) Create(ctx context.Context, sup modelsup.Support, file io.Reader, ext string) (error) {
	filename := uuid.New().String() + ext

	uploadPath := filepath.Join("support/", filename)
	if file != nil {
		if err := u.fileStorage.Save(uploadPath, file); err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}
		sup.ImagePath = uploadPath
	}
	return u.supRepo.Create(ctx, sup)
}

func (u *UseCase) Update(ctx context.Context, sup modelsup.Support, file io.Reader, ext string) (error) {
	filename := uuid.New().String() + ext

	uploadPath := filepath.Join("support/", filename)
	if file != nil {
		if err := u.fileStorage.Save(uploadPath, file); err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}
		sup.ImagePath = uploadPath
	}
	return u.supRepo.Update(ctx, sup)
}

func (u *UseCase) GetAllSups(ctx context.Context, superID int64) ([]modelsup.Support, error){
	return u.supRepo.GetAllSups(ctx, superID)
}
