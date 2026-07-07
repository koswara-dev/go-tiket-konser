package handler

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"go-tiket-konser/service"
	"go-tiket-konser/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ConcertHandler struct {
	service         service.ConcertService
	storageProvider service.StorageProvider
}

func NewConcertHandler(s service.ConcertService, sp service.StorageProvider) *ConcertHandler {
	return &ConcertHandler{service: s, storageProvider: sp}
}

func (h *ConcertHandler) UploadTumbnail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	file, err := c.FormFile("thumbnail")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Failed to get thumbnail file",
		})
		return
	}

	fileUrl, err := h.storageProvider.UploadFile(file, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Failed to upload thumbnail",
		})
		return
	}

	concert, err := h.service.GetConcertByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Failed to get concert",
		})
		return
	}

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	concert.ThumbnailURL = fileUrl
	concert.UpdatedBy = &userID
	if err := h.service.UpdateConcert(&concert); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Failed to update thumbnail",
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Thumbnail uploaded successfully",
		Data:    fileUrl,
	})
}

func (h *ConcertHandler) UploadRulesPDF(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{Success: false, Message: "ID Konser tidak valid"})
		return
	}

	file, err := c.FormFile("rules_pdf")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{Success: false, Message: "Berkas PDF tata tertib wajib dikirim"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" {
		c.JSON(http.StatusBadRequest, dto.WebResponse{Success: false, Message: "Format berkas tidak didukung: hanya diperbolehkan PDF (.pdf)"})
		return
	}

	fileURL, err := h.storageProvider.UploadFile(file, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{Success: false, Message: err.Error()})
		return
	}

	concert, err := h.service.GetConcertByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.WebResponse{Success: false, Message: "Konser tidak ditemukan"})
		return
	}

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	concert.RulesPDFURL = fileURL
	concert.UpdatedBy = &userID
	err = h.service.UpdateConcert(&concert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebResponse{Success: false, Message: "Gagal menyimpan berkas PDF tata tertib"})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "PDF Tata tertib konser berhasil diunggah dan disimpan",
		Data:    fileURL,
	})
}

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

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	concert := models.Concert{
		BaseModel: models.BaseModel{
			CreatedBy: &userID,
		},
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

	logger.Log.WithFields(logrus.Fields{
		"concert_id": concert.ID,
		"client_ip":  c.ClientIP(),
		"action":     "concert_created",
	}).Info("Konser baru berhasil ditambahkan")
}

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

	logger.Log.WithFields(logrus.Fields{
		"query":     req,
		"client_ip": c.ClientIP(),
		"action":    "concert_list_retrieved",
	}).Info("Daftar konser berhasil diambil")
}

func (h *ConcertHandler) GetConcertByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
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

func (h *ConcertHandler) UpdateConcert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	concert, err := h.service.GetConcertByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "concert not found" {
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

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	concert.Title = input.Title
	concert.Description = input.Description
	concert.Date = input.Date
	concert.Venue = input.Venue
	concert.Status = input.Status
	concert.UpdatedBy = &userID

	if err := h.service.UpdateConcert(&concert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update concert: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, concert)
}

func (h *ConcertHandler) DeleteConcert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	if err := h.service.DeleteConcert(id, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "concert not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Concert not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Concert deleted successfully"})
}

func mapError(err error) (int, string) {
	if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "concert not found" {
		return http.StatusNotFound, "Concert not found"
	}
	if errors.Is(err, models.ErrConcertAlreadyExists) {
		return http.StatusConflict, "Concert already exists"
	}
	return http.StatusInternalServerError, err.Error()
}
