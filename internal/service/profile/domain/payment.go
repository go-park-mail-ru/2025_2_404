package user

import (
	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentSucceeded PaymentStatus = "succeeded"
	PaymentCanceled  PaymentStatus = "canceled"
)

type Payment struct {
	ID            uuid.UUID
	ClientID      uuid.UUID
	AmountRub     uint32        `json:"amount"`
	PaymentMethod string        `json:"payment_method"`
	Status        PaymentStatus `json:"status"`
	YooPaymentID  string        `json:"yoo_payment_id"`
	CreatedTime   string        `json:"created_at"`
}

type PaymentRequest struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Confirmation struct {
		Type      string `json:"type"`
		ReturnURL string `json:"return_url"`
	} `json:"confirmation"`
	Description string                 `json:"description"`
	Capture     bool                   `json:"capture"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type PaymentResponse struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Paid         bool   `json:"paid"`
	Confirmation struct {
		Type string `json:"type"`
		URL  string `json:"confirmation_url,omitempty"`
	} `json:"confirmation"`
}

type BalanceResponse struct {
	Balance  int64     `json:"balance"`
	Payments []Payment `json:"payments"`
}

type YooKassaWebhook struct {
	Event  string `json:"event"`
	Object struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Amount struct {
			Value string `json:"value"`
		} `json:"amount"`
	} `json:"object"`
}
