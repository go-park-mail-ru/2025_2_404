package http

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"2025_2_404/internal/service/slot/domain/metric"
	"2025_2_404/internal/service/slot/domain/slot"
	"2025_2_404/pkg"
	convertimage "2025_2_404/pkg/convertImage"
	"2025_2_404/pkg/utils"
	adpb "2025_2_404/protos/gen/go/ad"
	slotpb "2025_2_404/protos/gen/go/slot"
	storagepb "2025_2_404/protos/gen/go/storage"
)

type SlotHandler struct {
	client        slotpb.SlotServClient
	adClient      adpb.AdServClient
	storageClient storagepb.StorageClient
	tmpl          *template.Template
}

func NewSlotHandler(client slotpb.SlotServClient, adClient adpb.AdServClient, storageClient storagepb.StorageClient) *SlotHandler {
	tmpl := template.Must(template.ParseFiles("template/template.html"))
	return &SlotHandler{
		client:        client,
		adClient:      adClient,
		storageClient: storageClient,
		tmpl:          tmpl,
	}
}

func (h *SlotHandler) ServeSlot(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на отображение слота: %s", r.URL.Path)

	vars := mux.Vars(r)
	slotID := vars["id"]

	if slotID == "" {
		log.Printf("Ошибка: ID слота отсутствует в URL")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Запрос данных слота с ID: %s", slotID)
	resp, err := h.client.GetSlot(ctx, &slotpb.GetSlotRequest{Id: slotID})
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("Ошибка при получении слота (ID=%s): код=%d, сообщение=%q", slotID, st.Code(), st.Message())
		if st.Code() == 5 { // NotFound
			http.Error(w, "Not found slot", http.StatusNotFound)
		} else {
			http.Error(w, "Problem in logic back", http.StatusExpectationFailed)
		}
		return
	}

	log.Printf("Слот найден. Запрос изображения по пути: %s", resp.AdSlot.ImageSrc)
	imgData, err := h.storageClient.Get(ctx, &storagepb.GetRequest{ImagePath: resp.AdSlot.ImageSrc})
	if err != nil {
		log.Printf("Ошибка при получении изображения (путь=%s): %v", resp.AdSlot.ImageSrc, err)
		http.Error(w, "Problem with load image", http.StatusNotFound)
		return
	}

	if len(imgData.ImageData) == 0 {
		log.Printf("Предупреждение: получены пустые данные изображения для пути %s", resp.AdSlot.ImageSrc)
	}

	var imageSrc string
	if len(imgData.ImageData) != 0 {

		imageSrc = convertimage.ConvertImageToBase64(imgData.ImageData, imgData.ContentType)
		log.Printf("Изображение успешно конвертировано в Base64: %s", imageSrc[:30]+"...")
		if imageSrc == "" {
			log.Printf("Ошибка: ConvertImageToBase64 вернула пустую строку")
			http.Error(w, "Problem with convert image", http.StatusBadRequest)
			return
		}
	}

	data := slot.SlotRenderData{
		Title:       resp.AdSlot.Title,
		Description: resp.AdSlot.Description,
		ImageData:   template.URL(imageSrc),
		Link:        resp.AdSlot.Link,
		Background:  resp.Slot.BackColor,
		Color:       resp.Slot.TextColor,
		Banner:      resp.AdSlot.Id,
		Slot:        slotID,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "frame-ancestors 'self' http://localhost:8000 http://localhost:8080 http://89.208.230.119:8000 http://terabithia.online https://flintmail.ru;")

	log.Printf("Рендеринг HTML для слота ID=%s", slotID)
	if err := h.tmpl.Execute(w, data); err != nil {
		log.Printf("Ошибка при рендеринге шаблона: %v", err)
		http.Error(w, "Render html bad", http.StatusNotFound)
	}
}

type slotDTO struct {
	SlotName       string `json:"slot_name"`
	MinCostAdv     int32  `json:"min_cost_adv"`
	FormatOfBanner string `json:"format_of_banner"`
	Status         string `json:"status"`
	BackColor      string `json:"back_color"`
	TextColor      string `json:"text_color"`
}

func (h *SlotHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на создание слота")
	var dto slotDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		log.Printf("Ошибка парсинга JSON при создании слота: %v", err)
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	req := &slotpb.CreateSlotRequest{
		Slot: &slotpb.Slot{
			SlotName:       dto.SlotName,
			MinCostAdv:     dto.MinCostAdv,
			FormatOfBanner: dto.FormatOfBanner,
			Status:         dto.Status,
			BackColor:      dto.BackColor,
			TextColor:      dto.TextColor,
		},
	}

	resp, err := h.client.CreateSlot(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("Ошибка gRPC при создании слота: код=%d, сообщение=%q", st.Code(), st.Message())
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	log.Printf("Слот успешно создан с ID=%s", resp.Id)
	pkg.JSONResponse(w, http.StatusCreated, "Slot created successfully", map[string]string{"id": resp.Id})
}

func (h *SlotHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на получение всех слотов")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	resp, err := h.client.ListSlots(ctx, &slotpb.ListSlotsRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("Ошибка gRPC при получении списка слотов: код=%d, сообщение=%q", st.Code(), st.Message())
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	log.Printf("Успешно получено %d слотов", len(resp.Slots))
	pkg.JSONResponse(w, http.StatusOK, "Slots retrieved successfully", resp.Slots)
}

func (h *SlotHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := uuid.Parse(id); err != nil {
		log.Printf("Неверный UUID в запросе: %s", id)
		http.Error(w, `{"error":"invalid UUID"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Получен запрос на получение слота с ID=%s", id)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	resp, err := h.client.GetSlot(ctx, &slotpb.GetSlotRequest{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("Ошибка gRPC при получении слота (ID=%s): код=%d, сообщение=%q", id, st.Code(), st.Message())
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	log.Printf("Слот с ID=%s успешно получен", id)
	pkg.JSONResponse(w, http.StatusOK, "Slot retrieved successfully", resp.Slot)
}

func (h *SlotHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("Получен запрос на обновление слота с ID=%s", id)
	var dto slotDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		log.Printf("Ошибка парсинга JSON при обновлении слота (ID=%s): %v", id, err)
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	req := &slotpb.UpdateSlotRequest{
		Slot: &slotpb.Slot{
			Id:             id,
			SlotName:       dto.SlotName,
			MinCostAdv:     dto.MinCostAdv,
			FormatOfBanner: dto.FormatOfBanner,
			Status:         dto.Status,
			BackColor:      dto.BackColor,
			TextColor:      dto.TextColor,
		},
	}

	_, err := h.client.UpdateSlot(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("Ошибка gRPC при обновлении слота (ID=%s): код=%d, сообщение=%q", id, st.Code(), st.Message())
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	log.Printf("Слот с ID=%s успешно обновлён", id)
	w.WriteHeader(http.StatusOK)
}

func (h *SlotHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("Получен запрос на удаление слота с ID=%s", id)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	_, err := h.client.DeleteSlot(ctx, &slotpb.DeleteSlotRequest{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("Ошибка gRPC при удалении слота (ID=%s): код=%d, сообщение=%q", id, st.Code(), st.Message())
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	log.Printf("Слот с ID=%s успешно удалён", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *SlotHandler) CreateMetric(w http.ResponseWriter, r *http.Request) {
	bannerID := r.URL.Query().Get("banner")
	slotID := r.URL.Query().Get("slot")
	action := r.URL.Query().Get("action")

	log.Printf("Получен запрос на запись метрики: banner=%q, slot=%q, action=%q", bannerID, slotID, action)

	if bannerID == "" || slotID == "" || action == "" {
		http.Error(w, "missing required params: banner, slot, action", http.StatusBadRequest)
		return
	}

	if action == "shown" {
		action = "impression"
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := h.client.CreateMetric(ctx, &slotpb.CreateMetricRequest{
		AdId:      bannerID,
		SlotId:    slotID,
		EventType: action,
	})
	if err != nil {
		log.Printf("Failed to record metric: %v", err)
		http.Error(w, "Metrics not created", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SlotHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("Получен запрос на получение статистики слота с ID=%s", id)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	res, err := h.client.GetMetrics(ctx, &slotpb.GetMetricsRequest{SlotId: id})
	if err != nil {
		log.Printf("Failed to get metric: %v", err)
		http.Error(w, "Not found metrics", http.StatusNotFound)
		return
	}

	var metricsForDay []metric.MetricsForDay

	for _, m := range res.GetMetrics() {
		metricForDay := metric.MetricsForDay{
			Clicks:      m.Clicks,
			Impressions: m.Impressions,
			EventDate:   m.EventData,
		}
		metricsForDay = append(metricsForDay, metricForDay)
	}

	metricRes := metric.GetMetricsResponse{
		SlotID:           res.GetSlotId(),
		TotalImpressions: res.GetTotalImpressions(),
		TotalClicks:      res.GetTotalClicks(),
		Metrics:          metricsForDay,
	}

	log.Printf("Статистика слота с ID=%s успешно получен", id)
	pkg.JSONResponse(w, http.StatusOK, "Slot retrieved successfully", metricRes)
}
