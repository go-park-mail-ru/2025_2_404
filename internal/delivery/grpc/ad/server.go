package ad

import (
	"2025_2_404/internal/delivery/grpc/interceptor"
	modelad "2025_2_404/internal/service/ad/domain/ad"
	modelfullad "2025_2_404/internal/service/ad/domain/ad_full_info"
	modeluser "2025_2_404/internal/service/ad/domain/user"
	adv1 "2025_2_404/protos/gen/go/ad"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type adUsecaseI interface{
	FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelfullad.AdFullInfo, error)
	Create(ctx context.Context, ad modelad.Ads) (error)
	Update(ctx context.Context, ad modelad.Ads) error
	Delete(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error
	GetOneAd(ctx context.Context, adID modelad.ID, clientID modeluser.ID) (modelfullad.AdFullInfo, int, error)
	GetAdDetailForSlot(ctx context.Context, id modelad.ID) (modelfullad.DetailID, error)
}

//TODO отдельныый интерфейс юкейз budgetI

type adService struct{
	adUsecase	adUsecaseI
	adv1.UnimplementedAdServServer
}

func New(adUsecase adUsecaseI) *adService{
	return &adService{
		adUsecase: adUsecase,
		//TODO budgetI
	}
}

func (s *adService) Create(ctx context.Context, req *adv1.CreateRequest) (*adv1.CreateResponse, error){
	clientID, err := interceptor.GetUserID(ctx)
	fmt.Printf("DEBUG INTERCEPTOR: Auth returned clientID string: '%s'\n", clientID)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}
	protoAd := req.GetAd()
	
	ad := modelad.Ads {
		Title: protoAd.Title,
		ClientID: clientID,
		Content: protoAd.Content,
		Budget: protoAd.Budget,
		ImagePath: protoAd.ImgPath,
		TargetUrl: protoAd.Targeturl,
	}
	if err := s.adUsecase.Create(ctx, ad); err != nil{
		return nil, err
	}

	return &adv1.CreateResponse{}, nil
}

func (s *adService) GetAllAds(ctx context.Context, req *adv1.GetAllAdsRequest) (*adv1.GetAllAdsResponse, error){
	clientID, err := interceptor.GetUserID(ctx)
	fmt.Printf("DEBUG INTERCEPTOR: Auth returned clientID string: '%s'\n", clientID)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	adsFull, err := s.adUsecase.FindByUserID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	var grpcAds []*adv1.Ad
	for _, a := range adsFull {
		grpcAds = append(grpcAds, &adv1.Ad{
			Id:        uuid.UUID(a.ID).String(),
			ClientID:  clientID.String(),
			Title:     a.Title,
			Content:   a.Content,
			Targeturl: a.TargetUrl,
			ImgPath:   a.ImgPath,
			Budget:    a.Budget,
			Status:    a.Status,
			StartAt:   a.StartAt.Format(time.RFC3339),
			EndAt:     a.EndAt.Format(time.RFC3339),
		})
	}

	return &adv1.GetAllAdsResponse{Ads: grpcAds}, nil
}

func (s *adService) Update(ctx context.Context, req *adv1.UpdateRequest) (*adv1.UpdateResponse, error){
	clientID, err := interceptor.GetUserID(ctx)
	fmt.Printf("DEBUG INTERCEPTOR: Auth returned clientID string: '%s'\n", clientID)
	
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}
	protoAd := req.GetAd()

	id, err := uuid.Parse(protoAd.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ad ID")
	}

	ad := modelad.Ads {
		ID: modelad.ID(id),
		ClientID: clientID,
		Title: protoAd.Title,
		Content: protoAd.Content,
		ImagePath: protoAd.ImgPath,
		TargetUrl: protoAd.Targeturl,
	}

	if err := s.adUsecase.Update(ctx, ad); err != nil{
		return nil, err
	}

	return &adv1.UpdateResponse{}, nil
}

func (s *adService) Delete(ctx context.Context, req *adv1.DeleteRequest) (*adv1.DeleteResponse, error){
	clientID, err := interceptor.GetUserID(ctx)
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ad ID")
	}
	adID := modelad.ID(id)

	if err := s.adUsecase.Delete(ctx, modelad.ID(adID), modeluser.ID(clientID)); err != nil{
		return nil, err
	}

	return &adv1.DeleteResponse{}, nil
}

func (s *adService) GetAd(ctx context.Context, req *adv1.GetAdRequest) (*adv1.GetAdResponse, error){
	clientID, err := interceptor.GetUserID(ctx)
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ad ID")
	}
	adID := modelad.ID(id)

	adFull, _, err := s.adUsecase.GetOneAd(ctx, adID, modeluser.ID(clientID))
	if err != nil{
		return &adv1.GetAdResponse{}, err
	}

	ad := &adv1.Ad{
		Id: uuid.UUID(adFull.ID).String(),
		Title: adFull.Title,
		Content: adFull.Content,
		Targeturl: adFull.TargetUrl,
		ImgPath: adFull.ImgPath,
		Budget: adFull.Budget,
		Status: adFull.Status,
		StartAt: adFull.StartAt.Format(time.RFC3339),
    	EndAt:   adFull.EndAt.Format(time.RFC3339),
	}

	return &adv1.GetAdResponse{Ad: ad}, nil
}

func (s *adService) GetAdDetailForSlot(ctx context.Context, req *adv1.GetAdDetailIDRequest) (* adv1.GetAdDetailIDResponse, error){
	id, err := uuid.Parse(req.GetAdId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ad ID")
	}

	detailId, err := s.adUsecase.GetAdDetailForSlot(ctx, modelfullad.ID(id))
	if err != nil{
		return nil, err
	}

	return &adv1.GetAdDetailIDResponse{
		AdDetailId: detailId.String(),
	}, nil
}

//TODO создать ручку пополнения бюджета в рекламе 