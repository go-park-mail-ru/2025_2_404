package ad

import (
	modelad "2025_2_404/internal/domain/models/ad"
	modelfullad "2025_2_404/internal/domain/models/ad_full_info"
	modeluser "2025_2_404/internal/domain/models/user"
	adv1 "2025_2_404/protos/gen/go/ad"
	"context"
)

type adUsecaseI interface{
	FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelad.Ads, error)
	Create(ctx context.Context, ad modelad.Ads) (error)
	Update(ctx context.Context, ad modelad.Ads) error
	Delete(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error
	GetOneAd(ctx context.Context, adID int64) (modelfullad.AdFullInfo, int, error)
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
	protoAd := req.GetAd()
	
	ad := modelad.Ads{
		ID:	modelad.ID(protoAd.Id),
		ClientID: 	modeluser.ID(protoAd.ClientID),
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
	clientID := req.GetClientID()

	ads, err := s.adUsecase.FindByUserID(ctx, modeluser.ID(clientID))
	if err != nil {
		return nil, err
	}

	var grpcAds []*adv1.Ad
	for _, a := range ads {
		grpcAds = append(grpcAds, &adv1.Ad{
			Id:        int64(a.ID),
			ClientID:  int64(a.ClientID),
			Title:     a.Title,
			Content:   a.Content,
			Targeturl: a.TargetUrl,
		})
	}

	return &adv1.GetAllAdsResponse{Ads: grpcAds}, nil
}

func (s *adService) Update(ctx context.Context, req *adv1.UpdateRequest) (*adv1.UpdateResponse, error){
	protoAd := req.GetAd()

	ad := modelad.Ads {
		ID: modelad.ID(protoAd.Id),
		ClientID: modeluser.ID(protoAd.ClientID),
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
	adID := req.GetId()
	clientID := req.GetClientID()

	if err := s.adUsecase.Delete(ctx, modelad.ID(adID), modeluser.ID(clientID)); err != nil{
		return nil, err
	}

	return &adv1.DeleteResponse{}, nil
}

func (s *adService) GetAd(ctx context.Context, req *adv1.GetAdRequest) (*adv1.GetAdResponse, error){
	adID := req.GetId()
	clientID := req.GetClientID()

	adFull, _, err := s.adUsecase.GetOneAd(ctx, adID)
	if err != nil{
		return &adv1.GetAdResponse{}, err
	}

	ad := &adv1.Ad{
		Id: int64(adFull.ID),
		ClientID: int64(clientID),
		Title: adFull.Title,
		Content: adFull.Content,
		Targeturl: adFull.TargetUrl,
	}

	return &adv1.GetAdResponse{Ad: ad}, nil
}


