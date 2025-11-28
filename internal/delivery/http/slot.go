package http

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"2025_2_404/pkg"
	"2025_2_404/pkg/utils"
	adpb "2025_2_404/protos/gen/go/ad"
	slotpb "2025_2_404/protos/gen/go/slot"
	"2025_2_404/internal/service/slot/domain/slot"
)

type SlotHandler struct {
	client   slotpb.SlotServClient
	adClient adpb.AdServClient
	tmpl     *template.Template
}

func NewSlotHandler(client slotpb.SlotServClient, adClient adpb.AdServClient) *SlotHandler {
	tmpl := template.Must(template.ParseFiles("template/template.html"))
	return &SlotHandler{client: client, adClient: adClient, tmpl: tmpl}
}

func (h *SlotHandler) ServeSlot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slotID := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetSlot(ctx, &slotpb.GetSlotRequest{Id: slotID})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == 5 { // NotFound
			http.Error(w, "", http.StatusNotFound)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	data := slot.SlotRenderData{
		Title:       resp.AdSlot.Title,
		Description: resp.AdSlot.Description,
		ImageSrc:    resp.AdSlot.ImageSrc,
		Link:        resp.AdSlot.Link,
		Background:  resp.Slot.BackColor,
		Color:       resp.Slot.TextColor,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "frame-ancestors 'self' http://localhost:8000 http://89.208.230.119:8000;")

	if err := h.tmpl.Execute(w, data); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
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
	var dto slotDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
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
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	pkg.JSONResponse(w, http.StatusCreated, map[string]string{"id": resp.Id})
}

func (h *SlotHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	resp, err := h.client.ListSlots(ctx, &slotpb.ListSlotsRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	pkg.JSONResponse(w, http.StatusOK, resp.Slots)
}

func (h *SlotHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, `{"error":"invalid UUID"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	resp, err := h.client.GetSlot(ctx, &slotpb.GetSlotRequest{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	pkg.JSONResponse(w, http.StatusOK, resp.Slot)
}

func (h *SlotHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var dto slotDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
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
			Id:              id,
			SlotName:        dto.SlotName,
			MinCostAdv:      dto.MinCostAdv,
			FormatOfBanner:  dto.FormatOfBanner,
			Status:          dto.Status,
			BackColor:       dto.BackColor,
			TextColor:       dto.TextColor,
		},
	}

	_, err := h.client.UpdateSlot(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SlotHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if auth := r.Header.Get("Authorization"); auth != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}

	_, err := h.client.DeleteSlot(ctx, &slotpb.DeleteSlotRequest{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		http.Error(w, `{"error":"`+st.Message()+`"}`, utils.HTTPStatusFromCode(st.Code()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}