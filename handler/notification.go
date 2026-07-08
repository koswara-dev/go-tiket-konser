package handler

import (
	"io"
	"net/http"

	"go-tiket-konser/dto"
	"go-tiket-konser/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	broker *service.NotificationBroker
}

func NewNotificationHandler(b *service.NotificationBroker) *NotificationHandler {
	return &NotificationHandler{broker: b}
}

func getUserIDStr(c *gin.Context) string {
	val, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	if u, ok := val.(uuid.UUID); ok {
		return u.String()
	}
	return ""
}

func (h *NotificationHandler) Stream(c *gin.Context) {
	userIDStr := getUserIDStr(c)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	ch := make(chan string, 10)
	h.broker.Register(userIDStr, ch)
	defer h.broker.Unregister(userIDStr, ch)

	// Send initial event
	c.SSEvent("info", "SSE connection established")

	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-ch:
			if ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		case <-c.Request.Context().Done():
			return false
		}
	})
}

func (h *NotificationHandler) GetHistory(c *gin.Context) {
	userIDStr := getUserIDStr(c)

	history, err := h.broker.GetHistory(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebResponse{
			Success: false,
			Message: "Gagal mengambil riwayat notifikasi",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Riwayat notifikasi berhasil diambil",
		Data:    history,
	})
}
