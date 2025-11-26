package http

import (
	"google.golang.org/grpc/codes"
	slotpb "2025_2_404/protos/gen/go/slot"
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

type SlotHandler struct {
	client slotpb.SlotServClient
	tmpl   *template.Template
}

func NewSlotHandler(client slotpb.SlotServClient) *SlotHandler {
	tmpl := template.Must(template.ParseFiles("templates/slot.html"))
	return &SlotHandler{client: client, tmpl: tmpl}
}

func (h *SlotHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/ads/serving/:id", h.ServeSlot)
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