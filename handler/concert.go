package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"go-tiket-konser/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConcertHandler struct {
	service service.ConcertService
}

func NewConcertHandler(service service.ConcertService) *ConcertHandler {
	return &ConcertHandler{service: service}
}

// CreateConcert handles POST /api/v1/concerts
func (h *ConcertHandler) CreateConcert(c *gin.Context) {
	var req dto.ConcertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Validasi input gagal",
			"error":   err.Error(),
		})
		return
	}

	parseDate, _ := time.Parse("2006-01-02", req.Date)

	concert := models.Concert{
		Title:       req.Title,
		Description: req.Description,
		Date:        parseDate,
		Venue:       req.Venue,
		Status:      req.Status,
	}

	err := h.service.CreateConcert(&concert)
	if err != nil {
		statusCode, errMsg := mapError(err)
		c.JSON(statusCode, gin.H{
			"status":  false,
			"message": "Gagal menambahkan konser",
			"error":   errMsg,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Konser berhasil ditambahkan",
		"data": dto.ConcertResponse{
			ID:          concert.ID,
			Title:       concert.Title,
			Description: concert.Description,
			Date:        concert.Date.Format("2006-01-02"),
			Venue:       concert.Venue,
			Status:      concert.Status,
			CreatedAt:   concert.CreatedAt,
			UpdatedAt:   concert.UpdatedAt,
		},
	})
}

// GetConcerts handles GET /api/v1/concerts
func (h *ConcertHandler) GetConcerts(c *gin.Context) {
	var req dto.ConcertQueryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Validasi parameter pencarian gagal",
			Data:    err.Error(),
		})
		return
	}

	concerts, meta, err := h.service.GetAllConcerts(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebResponse{
			Success: false,
			Message: "Gagal mengambil data konser",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Data berhasil diambil",
		Data:    concerts,
		Meta:    meta,
	})
}

// GetConcertByID handles GET /api/v1/concerts/:id
func (h *ConcertHandler) GetConcertByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	concert, err := h.service.GetConcertByID(id)
	if err != nil {
		statusCode, errMsg := mapError(err)
		c.JSON(statusCode, gin.H{
			"success": false,
			"message": errMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": dto.ConcertResponse{
			ID:          concert.ID,
			Title:       concert.Title,
			Description: concert.Description,
			Date:        concert.Date.Format("2006-01-02"),
			Venue:       concert.Venue,
			Status:      concert.Status,
			CreatedAt:   concert.CreatedAt,
			UpdatedAt:   concert.UpdatedAt,
		},
	})
}

// UpdateConcert handles PUT /api/v1/concerts/:id
func (h *ConcertHandler) UpdateConcert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	concert, err := h.service.GetConcertByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Concert not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input models.Concert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	concert.Title = input.Title
	concert.Description = input.Description
	concert.Date = input.Date
	concert.Venue = input.Venue
	concert.Status = input.Status

	if err := h.service.UpdateConcert(&concert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update concert: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, concert)
}

// DeleteConcert handles DELETE /api/v1/concerts/:id
func (h *ConcertHandler) DeleteConcert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteConcert(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Concert not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Concert deleted successfully"})
}

func mapError(err error) (int, string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, "Concert not found"
	}
	if errors.Is(err, models.ErrConcertAlreadyExists) {
		return http.StatusConflict, "Concert already exists"
	}
	return http.StatusInternalServerError, err.Error()
}
