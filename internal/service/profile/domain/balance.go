// Package user defines data structures related to user balance operations.
package user

type BalanceOp struct {
	AddAmount      uint32 `json:"add_amount"`
	SubtractAmount uint32 `json:"subtract_amount"`
	Balance        uint32 `json:"balance"`
}
