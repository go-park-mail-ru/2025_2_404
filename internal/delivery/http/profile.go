package http

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	user "2025_2_404/internal/service/profile/domain"
	"2025_2_404/pkg"
	pkgfile "2025_2_404/pkg/readerFile"
	"2025_2_404/pkg/utils"

	// pkgyookassa "2025_2_404/pkg/ReadYooKassaIP"
	pbAd "2025_2_404/protos/gen/go/ad"
	pbProfile "2025_2_404/protos/gen/go/profile"
	pbStorage "2025_2_404/protos/gen/go/storage"
)

type ProfileHandler struct {
	client        pbProfile.ProfileClient
	storageClient pbStorage.StorageClient
	adClient      pbAd.AdServClient
}

func NewProfileHandler(client pbProfile.ProfileClient, storageClient pbStorage.StorageClient, adClient pbAd.AdServClient) *ProfileHandler {
	return &ProfileHandler{client: client, storageClient: storageClient, adClient: adClient}
}

func (h *ProfileHandler) Show(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get("X-Request-ID")
	slog.Info("Show profile request", "req_id", reqID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	resp, err := h.client.Show(ctx, &pbProfile.ShowRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		slog.Error("Failed to show profile", "req_id", reqID, "error", st.Message(), "code", st.Code())
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	var resIMG *pbStorage.GetResponse
	imgPath := resp.AvatarPath
	if imgPath != "" {
		resIMG, err = h.storageClient.Get(ctx, &pbStorage.GetRequest{ImagePath: imgPath})
		if err != nil {
			slog.Warn("Failed to get avatar", "req_id", reqID, "image_path", imgPath, "error", err)
		}
	}

	var adsCount int64 = 0
	adResp, err := h.adClient.GetAdCount(ctx, &pbAd.GetAdCountRequest{})

	if err != nil {
		log.Printf("Failed to get ad count, req_id: %v, error: %v", reqID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		adsCount = adResp.Count
	}

	slog.Info("Profile shown successfully", "req_id", reqID, "user_name", resp.UserName)
	pkg.JSONResponse(w, http.StatusOK, "Profile retrieved successfully", map[string]interface{}{
		"user_name":    resp.UserName,
		"email":        resp.Email,
		"first_name":   resp.FirstName,
		"last_name":    resp.LastName,
		"company":      resp.Company,
		"phone":        resp.Phone,
		"profile_type": resp.ProfileType,
		"imageData":    resIMG,
		"ads_count":    adsCount,
		"created_at":   resp.CreatedAt,
	})
}

func parseMultipartForm(r *http.Request) error {
	const maxMemory = 32 << 20 // 32 MB
	return r.ParseMultipartForm(maxMemory)
}

func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get("X-Request-ID")
	slog.Info("Update profile request", "req_id", reqID)

	if err := parseMultipartForm(r); err != nil {
		slog.Error("Failed to parse multipart form", "req_id", reqID, "error", err)
		http.Error(w, `{"error":"failed to parse form"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}
	userName := r.FormValue("user_name")
	email := r.FormValue("email")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	company := r.FormValue("company")
	phone := r.FormValue("phone")
	profileType := r.FormValue("profile_type")

	var avatarPath string

	fileBytes, imageFilename, err := pkgfile.ExtractImage(r, "avatar/", "avatar")
	if err != nil {
		slog.Error("Failed to extract avatar", "req_id", reqID, "error", err)
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	if len(fileBytes) > 0 {
		// Отдельный контекст для storage
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := h.storageClient.Create(ctx, &pbStorage.CreateRequest{
			ImagePath: imageFilename,
			ImageData: fileBytes,
		})
		if err != nil {
			slog.Error("Failed to upload avatar", "req_id", reqID, "image_path", imageFilename, "error", err)
		} else {
			slog.Info("Avatar uploaded", "req_id", reqID, "image_path", imageFilename)
			avatarPath = imageFilename // ← важно: сохраняем путь!
		}
	}

	req := &pbProfile.UpdateRequest{
		UserName:    userName,
		Email:       email,
		Password:    r.FormValue("password"), // пароль не логируем!
		FirstName:   firstName,
		LastName:    lastName,
		Company:     company,
		Phone:       phone,
		ProfileType: profileType,
		AvatarPath:  avatarPath,
	}

	resp, err := h.client.Update(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		slog.Error("Failed to update profile", "req_id", reqID, "error", st.Message(), "code", st.Code())
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	slog.Info("Profile updated successfully", "req_id", reqID, "user_name", userName)
	pkg.JSONResponse(w, http.StatusOK, "Profile updated successfully", resp)
}

func (h *ProfileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get("X-Request-ID")
	slog.Info("Delete profile request", "req_id", reqID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	_, err := h.client.Delete(ctx, &pbProfile.DeleteRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		slog.Error("Failed to delete profile", "req_id", reqID, "error", st.Message(), "code", st.Code())
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	slog.Info("Profile deleted successfully", "req_id", reqID)
	w.WriteHeader(http.StatusNoContent)
}
func (h *ProfileHandler) ShowBalance(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get("X-Request-ID")
	slog.Info("Show balance request", "req_id", reqID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	resp, err := h.client.ShowBalance(ctx, &pbProfile.ShowBalanceRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		slog.Error("Failed to show balance", "req_id", reqID, "error", st.Message(), "code", st.Code())
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	payments, err := h.client.GetPaymentsByClientID(ctx, &pbProfile.PaymentsByClientIDRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		slog.Error("Failed to show history payment", "req_id", reqID, "error", st.Message(), "code", st.Code())
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	var paymentsResp []user.Payment

	for _, payment := range payments.GetPayments() {
		payID, err := uuid.Parse(payment.GetId())
		if err != nil {
			log.Printf("Ошибка: Взят неправильный uuid в истории платежей")
			http.Error(w, "Invalid payment ID format", http.StatusBadRequest)
			return
		}
		historyPayment := user.Payment{
			ID:            payID,
			AmountRub:     payment.GetAmount(),
			PaymentMethod: payment.GetMethodPayment(),
			Status:        user.PaymentStatus(payment.Status),
			YooPaymentID:  payment.GetYooPaymentId(),
			CreatedTime:   payment.CreatedAt,
		}
		paymentsResp = append(paymentsResp, historyPayment)
	}

	slog.Info("Balance shown successfully", "req_id", reqID, "balance", resp.Balance)
	pkg.JSONResponse(w, http.StatusOK, "Balance retrieved successfully", &user.BalanceResponse{
		Balance:  int64(resp.GetBalance()),
		Payments: paymentsResp,
	})
}

func (h *ProfileHandler) AddBalance(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get("X-Request-ID")
	slog.Info("Add balance request", "req_id", reqID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	var jsonReq user.BalanceOp
	if err := json.NewDecoder(r.Body).Decode(&jsonReq); err != nil {
		slog.Error("Invalid JSON in add balance", "req_id", reqID, "error", err)
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	req := &pbProfile.AddBalanceRequest{
		AddAmount: jsonReq.AddAmount,
	}

	resp, err := h.client.AddBalance(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		slog.Error("Failed to add balance", "req_id", reqID, "amount", jsonReq.AddAmount, "error", st.Message(), "code", st.Code())
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	slog.Info("Balance added successfully", "req_id", reqID, "amount", jsonReq.AddAmount, "new_balance", req.AddAmount)
	pkg.JSONResponse(w, http.StatusOK, "Balance added successfully", resp)
}

func (h *ProfileHandler) SubtractBalance(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get("X-Request-ID")
	slog.Info("Subtract balance request", "req_id", reqID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	var jsonReq user.BalanceOp
	if err := json.NewDecoder(r.Body).Decode(&jsonReq); err != nil {
		slog.Error("Invalid JSON in subtract balance", "req_id", reqID, "error", err)
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	req := &pbProfile.SubtractBalanceRequest{
		SubAmount: jsonReq.SubtractAmount,
	}

	resp, err := h.client.SubtractBalance(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		slog.Error("Failed to subtract balance", "req_id", reqID, "amount", jsonReq.SubtractAmount, "error", st.Message(), "code", st.Code())
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	slog.Info("Balance subtracted successfully", "req_id", reqID, "amount", jsonReq.SubtractAmount, "new_balance-", req.SubAmount)
	pkg.JSONResponse(w, http.StatusOK, "Balance subtracted successfully", resp)
}

func (h *ProfileHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get("X-Request-ID")
	slog.Info("Create payment request", "req_id", reqID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	var jsonReq user.Payment
	if err := json.NewDecoder(r.Body).Decode(&jsonReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if jsonReq.AmountRub > 100000 {
		http.Error(w, "Payment more 100000 RUB", http.StatusBadRequest)
		return
	}

	reqPayment := &pbProfile.PaymentCreateRequest{
		Amount:        jsonReq.AmountRub,
		PaymentMethod: jsonReq.PaymentMethod,
	}

	resp, err := h.client.CreatePayment(ctx, reqPayment)
	if err != nil {
		st, _ := status.FromError(err)
		slog.Error("Failed to create payment", "req_id", reqID, "error", st.Message(), "code", st.Code())
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	slog.Info("Payment created successfully", "req_id", reqID)
	pkg.JSONResponse(w, http.StatusOK, "Payment created successfully", resp)
}

// func (h *ProfileHandler) HandleYooKassaWebhook(w http.ResponseWriter, r *http.Request) {
//     ip := r.Header.Get("X-Forwarded-For")
// 	if ip == "" {
// 		ip = r.RemoteAddr
// 	}

//     // if !pkgyookassa.IsYooKassaIP(ip) {
// 	// 	slog.Warn("Rejected webhook from unauthorized IP", "ip", ip)
// 	// 	http.Error(w, "Forbidden", http.StatusForbidden)
// 	// 	return
// 	// }

//     var notification user.YooKassaNotification

//     if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
// 		http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 		return
// 	}

//     if notification.Object.Status == "waiting_for_capture"{
//         return
//     }

//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     _, err := h.client.CheckPayment(ctx, &pbProfile.CheckPaymentRequest{
//         YookassaId: notification.Object.ID,
//         Status: notification.Object.Status,
//     })

//     if err != nil {
//         http.Error(w, "status not update", http.StatusBadGateway)
//         return
//     }

//     w.WriteHeader(http.StatusOK)
// }

func (h *ProfileHandler) HandleYooKassaWebhook(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	// Логируем входящий webhook
	slog.Info("Received YooKassa webhook", "ip", ip)

	var notification user.YooKassaWebhook
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		slog.Error("Failed to decode YooKassa webhook JSON", "ip", ip, "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	yookassaID := notification.Object.ID
	status := notification.Object.Status
	rublesStr := strings.Split(notification.Object.Amount.Value, ".")[0]

	slog.Info("Parsed YooKassa webhook",
		"ip", ip,
		"yookassa_id", yookassaID,
		"status", status,
		"amount", rublesStr,
	)

	if status == "waiting_for_capture" {
		slog.Info("Skipping 'waiting_for_capture' status", "yookassa_id", yookassaID)
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.client.UpdatePaymentStatus(ctx, &pbProfile.PaymentStatusRequest{
		YookassaId: yookassaID,
		Status:     status,
		Amount:     rublesStr,
	})
	if err != nil {
		slog.Error("Failed to update payment status via gRPC",
			"yookassa_id", yookassaID,
			"status", status,
			"error", err,
		)
		http.Error(w, "status not update", http.StatusBadRequest)
		return
	}

	slog.Info("Successfully processed YooKassa webhook", "yookassa_id", yookassaID)
	w.WriteHeader(http.StatusOK)
}

// func (h *ProfileHandler) HandleYooKassaWebhook(w http.ResponseWriter, r *http.Request) {
//     ip := r.Header.Get("X-Forwarded-For")
// 	if ip == "" {
// 		ip = r.RemoteAddr
// 	}

//     // if !pkgyookassa.IsYooKassaIP(ip) {
// 	// 	slog.Warn("Rejected webhook from unauthorized IP", "ip", ip)
// 	// 	http.Error(w, "Forbidden", http.StatusForbidden)
// 	// 	return
// 	// }

//     var notification user.YooKassaNotification

//     if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
// 		http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 		return
// 	}

//     if notification.Object.Status == "waiting_for_capture"{
//         return
//     }

//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     _, err := h.client.CheckPayment(ctx, &pbProfile.CheckPaymentRequest{
//         YookassaId: notification.Object.ID,
//         Status: notification.Object.Status,
//     })

//     if err != nil {
//         http.Error(w, "status not update", http.StatusBadGateway)
//         return
//     }

//     w.WriteHeader(http.StatusOK)
// }

// func (h *ProfileHandler) HandleYooKassaWebhook(w http.ResponseWriter, r *http.Request) {
//     ip := r.Header.Get("X-Forwarded-For")
//     if ip == "" {
//         ip = r.RemoteAddr
//     }

//     // Логируем входящий webhook
//     slog.Info("Received YooKassa webhook", "ip", ip, "yookassa_id", "unknown", "status", "unknown")

//     var notification user.YooKassaNotification
//     if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
//         slog.Error("Failed to decode YooKassa webhook JSON", "ip", ip, "error", err)
//         http.Error(w, "Invalid JSON", http.StatusBadRequest)
//         return
//     }

//     yookassaID := notification.ID
//     status := notification.Status
//     rublesStr := strings.Split(notification.Amount.Value, ".")[0]

//     slog.Info("Parsed YooKassa webhook",
//         "ip", ip,
//         "yookassa_id", yookassaID,
//         "status", status,
//         "amount", notification.Amount.Value,
//     )

//     if status == "waiting_for_capture" {
//         slog.Info("Skipping 'waiting_for_capture' status", "yookassa_id", yookassaID)
//         w.WriteHeader(http.StatusOK)
//         return
//     }

//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     _, err := h.client.UpdatePaymentStatus(ctx, &pbProfile.PaymentStatusRequest{
//         YookassaId: yookassaID,
//         Status:     status,
//         Amount:     rublesStr,
//     })
//     if err != nil {
//         slog.Error("Failed to update payment status via gRPC",
//             "yookassa_id", yookassaID,
//             "status", status,
//             "error", err,
//         )
//         http.Error(w, "status not update", http.StatusBadRequest)
//         return
//     }

//     slog.Info("Successfully processed YooKassa webhook", "yookassa_id", yookassaID)
//     w.WriteHeader(http.StatusOK)
// }
