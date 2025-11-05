package ad

import (
	modelad "2025_2_404/internal/domain/models/ad"
	modelfullad "2025_2_404/internal/domain/models/ad_full_info"
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/google/uuid"
)

type adRepositoryI interface {
	FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelad.Ads, error)
	Create(ctx context.Context, ad modelad.Ads) (error)
	GetOneAd(ctx context.Context, adID int64) (modelfullad.AdFullInfo, error)
	Update(ctx context.Context, ad modelad.Ads) error
	Delete(ctx context.Context, adID int64) error
}

type fileStorageI interface{
	Save(uploadPath string, data io.Reader) error
	ReadImageAsBytes(path string) ([]byte, error)
}

type UseCase struct {
	adRepo adRepositoryI
	fileStorage	fileStorageI
}

func New(adRepo adRepositoryI, filefileStorage fileStorageI) *UseCase {
	return &UseCase{
		adRepo: adRepo,
		fileStorage: filefileStorage,
	}
}

func (u *UseCase) FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelad.Ads, error) {
	return u.adRepo.FindByUserID(ctx, userID)
}

func (u *UseCase) Create(ctx context.Context, ad modelad.Ads, file io.Reader, ext string) (error) {
	filename := uuid.New().String() + ext

	uploadPath := filepath.Join("ad/", filename)
	if file != nil {
		if err := u.fileStorage.Save(uploadPath, file); err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}
		ad.ImagePath = uploadPath
	}
	return u.adRepo.Create(ctx, ad)
}

func (u *UseCase) Update(ctx context.Context, ad modelad.Ads, file io.Reader, ext string) error {
	filename := uuid.New().String() + ext

	uploadPath := filepath.Join("ad/", filename)
	if file != nil {
		if err := u.fileStorage.Save(uploadPath, file); err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}
		ad.ImagePath = uploadPath
	}
	return u.adRepo.Update(ctx, ad)
}

func (u *UseCase) Delete(ctx context.Context, adID int64) error {
	return u.adRepo.Delete(ctx, adID)
}

func (u *UseCase) GetOneAd(ctx context.Context, adID int64) (modelfullad.AdFullInfo, int, []byte, error) {
	adInfo, err := u.adRepo.GetOneAd(ctx, adID)
	conversion := -1
	if err != nil {
		return modelfullad.AdFullInfo{}, conversion, nil, fmt.Errorf("Failed to get ad with id error %w", err)
	}
	if adInfo.Impressions != 0{
		conversion = adInfo.Clicks / adInfo.Impressions
	}
	imgBytes, err := u.fileStorage.ReadImageAsBytes(adInfo.ImgPath)
	if err != nil {
		return  modelfullad.AdFullInfo{}, conversion, nil, fmt.Errorf("problem in convert")
	}

	return adInfo, conversion, imgBytes, nil 
}
