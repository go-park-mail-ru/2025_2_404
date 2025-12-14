package http

import (
	"2025_2_404/internal/service/profile/domain"
	"2025_2_404/pkg"
	pkgfile "2025_2_404/pkg/readerFile"
	"2025_2_404/pkg/utils"
	pbProfile "2025_2_404/protos/gen/go/profile"
	pbStorage "2025_2_404/protos/gen/go/storage"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ProfileHandler struct {
	client pbProfile.ProfileClient
	storageClient pbStorage.StorageClient
}

func NewProfileHandler(client pbProfile.ProfileClient, storageClient pbStorage.StorageClient) *ProfileHandler {
	return &ProfileHandler{client: client, storageClient: storageClient}
}

func (h *ProfileHandler) Show(w http.ResponseWriter, r *http.Request) {
    reqID := r.Header.Get("X-Request-ID")
    slog.Info("📥 Show profile request", "req_id", reqID)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if auth := r.Header.Get("Authorization"); auth != "" {
        ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
    }

    resp, err := h.client.Show(ctx, &pbProfile.ShowRequest{})
    if err != nil {
        st, _ := status.FromError(err)
        slog.Error("❌ Failed to show profile", "req_id", reqID, "error", st.Message(), "code", st.Code())
        http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
        return
    }

    var resIMG *pbStorage.GetResponse
    imgPath := resp.AvatarPath
    if imgPath != "" {
        resIMG, err = h.storageClient.Get(ctx, &pbStorage.GetRequest{ImagePath: imgPath})
        if err != nil {
            slog.Warn("⚠️ Failed to get avatar", "req_id", reqID, "image_path", imgPath, "error", err)
        }
    }

    slog.Info("✅ Profile shown successfully", "req_id", reqID, "user_name", resp.UserName)
    pkg.JSONResponse(w, http.StatusOK, "Profile retrieved successfully", map[string]interface{}{
        "user_name":    resp.UserName,
        "email":        resp.Email,
        "first_name":   resp.FirstName,
        "last_name":    resp.LastName,
        "company":      resp.Company,
        "phone":        resp.Phone,
        "profile_type": resp.ProfileType,
        "imageData":  resIMG,
    })
}

func parseMultipartForm(r *http.Request) error {
	const maxMemory = 32 << 20 // 32 MB
	return r.ParseMultipartForm(maxMemory)
}

func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
    reqID := r.Header.Get("X-Request-ID")
    slog.Info("📥 Update profile request", "req_id", reqID)

    if err := parseMultipartForm(r); err != nil {
        slog.Error("❌ Failed to parse multipart form", "req_id", reqID, "error", err)
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
        slog.Error("❌ Failed to extract avatar", "req_id", reqID, "error", err)
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
            slog.Error("❌ Failed to upload avatar", "req_id", reqID, "image_path", imageFilename, "error", err)
        } else {
            slog.Info("✅ Avatar uploaded", "req_id", reqID, "image_path", imageFilename)
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
        slog.Error("❌ Failed to update profile", "req_id", reqID, "error", st.Message(), "code", st.Code())
        http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
        return
    }

    slog.Info("✅ Profile updated successfully", "req_id", reqID, "user_name", userName)
    pkg.JSONResponse(w, http.StatusOK, "Profile updated successfully", resp)
}

func (h *ProfileHandler) Delete(w http.ResponseWriter, r *http.Request) {
    reqID := r.Header.Get("X-Request-ID")
    slog.Info("📥 Delete profile request", "req_id", reqID)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if auth := r.Header.Get("Authorization"); auth != "" {
        ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
    }

    _, err := h.client.Delete(ctx, &pbProfile.DeleteRequest{})
    if err != nil {
        st, _ := status.FromError(err)
        slog.Error("❌ Failed to delete profile", "req_id", reqID, "error", st.Message(), "code", st.Code())
        http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
        return
    }

    slog.Info("✅ Profile deleted successfully", "req_id", reqID)
    w.WriteHeader(http.StatusNoContent)
}
func (h *ProfileHandler) ShowBalance(w http.ResponseWriter, r *http.Request) {
    reqID := r.Header.Get("X-Request-ID")
    slog.Info("📥 Show balance request", "req_id", reqID)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if auth := r.Header.Get("Authorization"); auth != "" {
        ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
    }

    resp, err := h.client.ShowBalance(ctx, &pbProfile.ShowBalanceRequest{})
    if err != nil {
        st, _ := status.FromError(err)
        slog.Error("❌ Failed to show balance", "req_id", reqID, "error", st.Message(), "code", st.Code())
        http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
        return
    }

    slog.Info("✅ Balance shown successfully", "req_id", reqID, "balance", resp.Balance)
    pkg.JSONResponse(w, http.StatusOK, "Balance retrieved successfully", resp)
}

func (h *ProfileHandler) AddBalance(w http.ResponseWriter, r *http.Request) {
    reqID := r.Header.Get("X-Request-ID")
    slog.Info("📥 Add balance request", "req_id", reqID)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if auth := r.Header.Get("Authorization"); auth != "" {
        ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
    }

    var jsonReq user.BalanceOp
    if err := json.NewDecoder(r.Body).Decode(&jsonReq); err != nil {
        slog.Error("❌ Invalid JSON in add balance", "req_id", reqID, "error", err)
        http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
        return
    }

    req := &pbProfile.AddBalanceRequest{
        AddAmount: jsonReq.AddAmount,
    }

    resp, err := h.client.AddBalance(ctx, req)
    if err != nil {
        st, _ := status.FromError(err)
        slog.Error("❌ Failed to add balance", "req_id", reqID, "amount", jsonReq.AddAmount, "error", st.Message(), "code", st.Code())
        http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
        return
    }

    slog.Info("✅ Balance added successfully", "req_id", reqID, "amount", jsonReq.AddAmount, "new_balance", req.AddAmount)
    pkg.JSONResponse(w, http.StatusOK, "Balance added successfully", resp)
}

func (h *ProfileHandler) SubtractBalance(w http.ResponseWriter, r *http.Request) {
    reqID := r.Header.Get("X-Request-ID")
    slog.Info("📥 Subtract balance request", "req_id", reqID)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if auth := r.Header.Get("Authorization"); auth != "" {
        ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
    }

    var jsonReq user.BalanceOp
    if err := json.NewDecoder(r.Body).Decode(&jsonReq); err != nil {
        slog.Error("❌ Invalid JSON in subtract balance", "req_id", reqID, "error", err)
        http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
        return
    }

    req := &pbProfile.SubtractBalanceRequest{
        SubAmount: jsonReq.SubtractAmount,
    }

    resp, err := h.client.SubtractBalance(ctx, req)
    if err != nil {
        st, _ := status.FromError(err)
        slog.Error("❌ Failed to subtract balance", "req_id", reqID, "amount", jsonReq.SubtractAmount, "error", st.Message(), "code", st.Code())
        http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
        return
    }

    slog.Info("✅ Balance subtracted successfully", "req_id", reqID, "amount", jsonReq.SubtractAmount, "new_balance-", req.SubAmount)
    pkg.JSONResponse(w, http.StatusOK, "Balance subtracted successfully", resp)
}

func (h *ProfileHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
    reqID := r.Header.Get("X-Request-ID")
    slog.Info("📥 Create payment request", "req_id", reqID)

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

    reqPayment := &pbProfile.PaymentCreateRequest{
        Amount: jsonReq.AmountRub,
        PaymentMethod: jsonReq.PaymentMethod,
    }

    resp, err := h.client.CreatePayment(ctx, reqPayment)
    if err != nil{
        st, _ := status.FromError(err)
        slog.Error("❌ Failed to create payment", "req_id", reqID, "error", st.Message(), "code", st.Code())
        http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
        return
    }

    reqAddBalance := &pbProfile.AddBalanceRequest{
        AddAmount: jsonReq.AmountRub,
    }

    _, err = h.client.AddBalance(ctx, reqAddBalance)
    if err != nil {
        st, _ := status.FromError(err)
        slog.Error("❌ Failed to add balance after payment", "req_id", reqID, "amount", jsonReq.AmountRub, "error", st.Message(), "code", st.Code())
        http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
        return
    }

    slog.Info("✅ Payment created successfully", "req_id", reqID)
    pkg.JSONResponse(w, http.StatusOK, "Payment created successfully", resp)
}