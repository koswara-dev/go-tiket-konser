package handler

import (
	"net/http"

	"go-tiket-konser/dto"
	"go-tiket-konser/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingHandler struct {
	service service.BookingService
}

func NewBookingHandler(service service.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req dto.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi input gagal",
			"error":   err.Error(),
		})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Akses tidak sah"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	booking, err := h.service.CreateBooking(&req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Reservasi tiket gagal diproses",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pemesanan tiket berhasil dikonfirmasi!",
		"data":    booking,
	})
}

func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID booking tidak valid"})
		return
	}

	// 1. Ekstrak userId dan role hasil suntikan JWTAuthMiddleware
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses tidak sah"})
		return
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tipe data user ID tidak valid"})
		return
	}
	role := c.MustGet("role").(string)

	// 2. kirim parameter audit keamanan layer service
	booking, err := h.service.GetBookingByID(id, userID, role)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Faktur pemesanan tiket tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    booking,
	})
}
