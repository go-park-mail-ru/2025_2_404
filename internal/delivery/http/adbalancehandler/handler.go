package adbalancehandler

import()

type adBalanceUsecaseI interface{

}

type adBalanceHandler struct{
	adBalanceUsecase adBalanceUsecaseI
}

func New(adBalanceUsecase adBalanceUsecaseI) *adBalanceHandler{
	return &adBalanceHandler{
		adBalanceUsecase:	adBalanceUsecase,
	}
}