package http

import (
	"google.golang.org/grpc/codes"
	slotpb "2025_2_404/protos/gen/go/slot"
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"2025_2_404/pkg/utils"
	"google.golang.org/grpc/status"
)

type SlotHandler struct {
	client slotpb.SlotServClient
	tmpl   *template.Template
}

func NewSlotHandler(client slotpb.SlotServClient) *SlotHandler {
	tmpl := template.Must(template.ParseFiles("template/template.html"))
	return &SlotHandler{client: client, tmpl: tmpl}
}

func (h *SlotHandler) RegisterRoutes(r *gin.Engine) {
	r.Group("/slots")
	{
		r.GET("/serving/:id", h.ServeSlot)
		r.POST("", h.Create)
		r.GET("", h.GetAll)
		r.GET("/:id", h.GetOne)
		r.PUT("/:id", h.Update)
		r.DELETE("/:id", h.Delete)
	}
}

type slotRenderData struct {
	Title      string
	Background string
	Color      string
}

func (h *SlotHandler) ServeSlot(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &slotpb.GetSlotRequest{Id: id}
	resp, err := h.client.GetSlot(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	data := slotRenderData{
		Title:      resp.Slot.SlotName,
		Background: resp.Slot.BackColor,
		Color:      resp.Slot.TextColor,
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("X-Content-Type-Options", "nosniff")
	if err := h.tmpl.Execute(c.Writer, data); err != nil {
		c.Status(http.StatusInternalServerError)
	}
}

type slotDTO struct {
	SlotName       string `json:"slot_name" binding:"required"`
	MinCostAdv     int32  `json:"min_cost_adv" binding:"required"`
	FormatOfBanner string `json:"format_of_banner" binding:"required"`
	Status         string `json:"status"`
	BackColor      string `json:"back_color" binding:"required"`
	TextColor      string `json:"text_color" binding:"required"`
}

func (h *SlotHandler) Create(c *gin.Context) {
	var dto slotDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
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
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": resp.Id})
}

func (h *SlotHandler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	}

	req := &slotpb.ListSlotsRequest{}
	resp, err := h.client.ListSlots(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusOK, resp.Slots)
}

func (h *SlotHandler) GetOne(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	}

	req := &slotpb.GetSlotRequest{Id: id}
	resp, err := h.client.GetSlot(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusOK, resp.Slot)
}

func (h *SlotHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var dto slotDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
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
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *SlotHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	}

	req := &slotpb.DeleteSlotRequest{Id: id}
	_, err := h.client.DeleteSlot(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.Status(http.StatusNoContent)
}