package http

import (
	"2025_2_404/internal/service/profile/config"
	user "2025_2_404/internal/service/profile/domain"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const YooKassaAPI = "https://api.yookassa.ru/v3/payments"

type YooKassaHttp struct{
	secretKey	string
	shopID		string
}

func New(cfg *config.Config) *YooKassaHttp{
	return &YooKassaHttp{
		secretKey: cfg.PaymentConfig.SecretKey,
		shopID: cfg.PaymentConfig.ShopID,
	}
}

func (y *YooKassaHttp) CreatePayment(ctx context.Context, payment user.Payment)(string, error){
	paymentReq := user.PaymentRequest{}
	paymentReq.Amount.Value = fmt.Sprintf("%d.00", payment.AmountRub)
	paymentReq.Amount.Currency = "RUB"
	paymentReq.Confirmation.Type = "redirect"
	paymentReq.Confirmation.ReturnURL = "https://adnet.website/balance"
	paymentReq.Description = "Оплата заказа"

	body, _ := json.Marshal(paymentReq)

	client := &http.Client{}
	httpReq, _ := http.NewRequest("POST", YooKassaAPI, bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Idempotence-Key", payment.YooPaymentID)
	httpReq.SetBasicAuth(y.shopID, y.secretKey)

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("Bad Request to YooKassa")
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("YooKassa error: %d, body: %s\n", resp.StatusCode, string(respBody))
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("YooKassa error: %d, body: %s\n", resp.StatusCode, string(respBody))
		return "", err
	}

	var paymentResp user.PaymentResponse
	json.Unmarshal(respBody, &paymentResp)

	return paymentResp.Confirmation.URL, nil
}