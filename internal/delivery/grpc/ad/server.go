package ad

import (
	"2025_2_404/internal/delivery/grpc/interceptor"
	modelad "2025_2_404/internal/service/ad/domain/ad"
	modelfullad "2025_2_404/internal/service/ad/domain/ad_full_info"
	modeluser "2025_2_404/internal/service/ad/domain/user"
	adv1 "2025_2_404/protos/gen/go/ad"
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type adUsecaseI interface{
	FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelad.Ads, error)
	Create(ctx context.Context, ad modelad.Ads) (error)
	Update(ctx context.Context, ad modelad.Ads) error
	Delete(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error
	GetOneAd(ctx context.Context, adID modelad.ID) (modelfullad.AdFullInfo, int, error)
}

type adService struct{
	adUsecase	adUsecaseI
	adv1.UnimplementedAdServServer
}

func New(adUsecase adUsecaseI) *adService{
	return &adService{
		adUsecase: adUsecase,
	}
}

func (s *adService) Create(ctx context.Context, req *adv1.CreateRequest) (*adv1.CreateResponse, error){
	clientID, err := interceptor.GetUserID(ctx)
	fmt.Printf("DEBUG INTERCEPTOR: Auth returned clientID string: '%s'\n", clientID)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}
	protoAd := req.GetAd()
	
	ad := modelad.Ads{
		ClientID: 	clientID,
		Title: protoAd.Title,
		Content: protoAd.Content,
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

	ads, err := s.adUsecase.FindByUserID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	var grpcAds []*adv1.Ad
	for _, a := range ads {
		grpcAds = append(grpcAds, &adv1.Ad{
			Id:        uuid.UUID(a.ID).String(),
			ClientID:  clientID.String(),
			Title:     a.Title,
			Content:   a.Content,
			Targeturl: a.TargetUrl,
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

	adFull, _, err := s.adUsecase.GetOneAd(ctx, adID)
	if err != nil{
		return &adv1.GetAdResponse{}, err
	}

	ad := &adv1.Ad{
		Id: uuid.UUID(adFull.ID).String(),
		ClientID: uuid.UUID(clientID).String(),
		Title: adFull.Title,
		Content: adFull.Content,
		Targeturl: adFull.TargetUrl,
	}

	return &adv1.GetAdResponse{Ad: ad}, nil
}


