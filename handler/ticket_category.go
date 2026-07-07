package handler

import (
	"errors"
	"net/http"

	"go-tiket-konser/models"
	"go-tiket-konser/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketCategoryHandler struct {
	service service.TicketCategoryService
}

func NewTicketCategoryHandler(service service.TicketCategoryService) *TicketCategoryHandler {
	return &TicketCategoryHandler{service: service}
}

// CreateTicketCategory handles POST /api/v1/ticket-categories
func (h *TicketCategoryHandler) CreateTicketCategory(c *gin.Context) {
	var category models.TicketCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)
	category.CreatedBy = &userID

	if err := h.service.CreateTicketCategory(&category); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Associated concert not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket category: " + err.Error()})
		return
	}

	// Fetch category with preloaded concert for response
	resCategory, err := h.service.GetTicketCategoryByID(category.ID)
	if err == nil {
		category = resCategory
	}

	c.JSON(http.StatusCreated, category)
}

// GetTicketCategories handles GET /api/v1/ticket-categories
func (h *TicketCategoryHandler) GetTicketCategories(c *gin.Context) {
	categories, err := h.service.GetAllTicketCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ticket categories: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetTicketCategoryByID handles GET /api/v1/ticket-categories/:id
func (h *TicketCategoryHandler) GetTicketCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	category, err := h.service.GetTicketCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateTicketCategory handles PUT /api/v1/ticket-categories/:id
func (h *TicketCategoryHandler) UpdateTicketCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	category, err := h.service.GetTicketCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input models.TicketCategory
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	category.ConcertID = input.ConcertID
	category.Name = input.Name
	category.Price = input.Price
	category.TotalQuota = input.TotalQuota
	category.AvailableQuota = input.AvailableQuota
	category.UpdatedBy = &userID

	if err := h.service.UpdateTicketCategory(&category); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Associated concert not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket category: " + err.Error()})
		return
	}

	// Fetch category with preloaded concert for response
	resCategory, err := h.service.GetTicketCategoryByID(category.ID)
	if err == nil {
		category = resCategory
	}

	c.JSON(http.StatusOK, category)
}

// DeleteTicketCategory handles DELETE /api/v1/ticket-categories/:id
func (h *TicketCategoryHandler) DeleteTicketCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	if err := h.service.DeleteTicketCategory(id, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket category deleted successfully"})
}
