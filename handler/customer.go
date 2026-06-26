package handler

import (
	"net/http"
	"strconv"

	"go-tiket-konser/dto"
	"go-tiket-konser/service"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service service.CustomerService
}

func NewCustomerHandler(service service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) GetAllCustomers(c *gin.Context) {
	customers, err := h.service.GetAllCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var res []dto.CustomerResponse
	for _, cust := range customers {
		res = append(res, dto.CustomerResponse{
			ID:        cust.ID,
			UserID:    cust.UserID,
			Name:      cust.Name,
			Email:     cust.Email,
			CreatedAt: cust.CreatedAt,
			UpdatedAt: cust.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, res)
}

func (h *CustomerHandler) GetCustomerByID(c *gin.Context) {
	idStr := c.Param("id")
	requestedID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	tokenCustomerID := c.MustGet("customer_id").(uint)
	tokenRole := c.MustGet("role").(string)

	if tokenRole != "admin" && tokenCustomerID != uint(requestedID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: tidak dapat mengakses data customer lain"})
		return
	}

	cust, err := h.service.GetCustomerByID(requestedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, dto.CustomerResponse{
		ID:        cust.ID,
		UserID:    cust.UserID,
		Name:      cust.Name,
		Email:     cust.Email,
		CreatedAt: cust.CreatedAt,
		UpdatedAt: cust.UpdatedAt,
	})
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	idStr := c.Param("id")
	requestedID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	tokenCustomerID := c.MustGet("customer_id").(uint)
	tokenRole := c.MustGet("role").(string)

	if tokenRole != "admin" && tokenCustomerID != uint(requestedID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: tidak dapat memperbarui data customer lain"})
		return
	}

	var req dto.CustomerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cust, err := h.service.UpdateCustomer(requestedID, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.CustomerResponse{
		ID:        cust.ID,
		UserID:    cust.UserID,
		Name:      cust.Name,
		Email:     cust.Email,
		CreatedAt: cust.CreatedAt,
		UpdatedAt: cust.UpdatedAt,
	})
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	idStr := c.Param("id")
	requestedID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	err = h.service.DeleteCustomer(requestedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer berhasil dihapus"})
}
