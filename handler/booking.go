package handler

import (
	"net/http"
	"strconv"

	"go-tiket-konser/dto"
	"go-tiket-konser/service"

	"github.com/gin-gonic/gin"
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

	booking, err := h.service.CreateBooking(&req)
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
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID booking tidak valid"})
		return
	}

	booking, err := h.service.GetBookingByID(id)
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
