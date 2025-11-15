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
	GetOneSup(ctx context.Context, supID modelsup.ID, clientID modeluser.ID) (modelsup.Support, error)
	GetAllSups(ctx context.Context, superID int64) ([]modelsup.Support, error)
	Update(ctx context.Context, sup modelsup.Support) error 
	Delete(ctx context.Context, supID modelsup.ID, clientID modeluser.ID) error
	FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelsup.Support, error) 
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
func (u *UseCase) Delete(ctx context.Context, supID modelsup.ID, clientID modeluser.ID) error {
	return u.supRepo.Delete(ctx, supID, clientID)
}

func (u *UseCase) GetOneUved(ctx context.Context, supID modelsup.ID, clientID modeluser.ID) (modelsup.Support, []byte, error) {
	supInfo, err := u.supRepo.GetOneSup(ctx, supID, clientID)
	if err != nil {
		return modelsup.Support{}, nil, fmt.Errorf("Failed to get ad with id error %w", err)
	}

	imgBytes, err := u.fileStorage.ReadImageAsBytes(supInfo.ImagePath)
	if err != nil {
		return  modelsup.Support{}, nil, fmt.Errorf("problem in convert")
	}

	return supInfo, imgBytes, nil 
}

func (u *UseCase) FindByUserID(ctx context.Context, clientID modeluser.ID) ([]modelsup.Support, error){
	return u.supRepo.FindByUserID(ctx, clientID)
}
