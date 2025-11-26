package http

import (
	adpb "2025_2_404/protos/gen/go/ad"
	slotpb "2025_2_404/protos/gen/go/slot"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"2025_2_404/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type SlotHandler struct {
	client slotpb.SlotServClient
	adClient adpb.AdServClient
	tmpl   *template.Template
}

func NewSlotHandler(client slotpb.SlotServClient, adClient adpb.AdServClient) *SlotHandler {
	tmpl := template.Must(template.ParseFiles("template/template.html"))
	return &SlotHandler{client: client,adClient: adClient, tmpl: tmpl}
}

func (h *SlotHandler) RegisterRoutes(r *gin.Engine) {
	slots := r.Group("/slots") 
	{
		slots.GET("/serving/:id", h.ServeSlot) // не работает 
		slots.POST("", h.Create)
		slots.GET("", h.GetAll)
		slots.GET("/:id", h.GetOne)
		slots.PUT("/:id", h.Update) // не работает тут мб такая же проблема с типом id slot должен быть uuid а передается string 
		slots.DELETE("/:id", h.Delete)
	}
}

type slotRenderData struct {
	Title      string 
	ImageSrc   string 
	Link       string 
	Background string 
	Color      string 
}

func (h *SlotHandler) ServeSlot(c *gin.Context) {
	slotID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slotReq := &slotpb.GetSlotRequest{Id: slotID}
	slotResp, err := h.client.GetSlot(ctx, slotReq)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == 5 { 
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	adReq := &adpb.GetAllAdsRequest{}
	adResp, err := h.adClient.GetAllAds(ctx, adReq)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(utils.HTTPStatusFromCode(st.Code()), gin.H{"error": "failed to fetch ads"})
		return
	}

	if len(adResp.Ads) == 0 {
		c.Status(http.StatusNotFound) 
		return
	}

	ad := adResp.Ads[0]

	data := slotRenderData{
		Title:      ad.Title,
		ImageSrc:   "", 
		Link:       ad.Targeturl,
		Background: slotResp.Slot.BackColor,
		Color:      slotResp.Slot.TextColor,
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
	slotID, err := uuid.Parse(id)
	if err != nil {
		fmt.Println("can`t parse into uuid err %w", err)
		return 
	} 

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	}

	req := &slotpb.GetSlotRequest{Id: slotID.String()}
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