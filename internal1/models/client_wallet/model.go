package clientwallet

import modeluser "2025_2_404/internal/domain/models/user"

type ID int
type Balance int

type Wallet struct{
	ID	ID	`json:"wallet_id"`
	ClientID	modeluser.ID
	Balance	Balance
}