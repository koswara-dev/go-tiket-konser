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

// GetAllCustomers godoc
// @Summary      Get all customers
// @Description  Get list of all customers with search, pagination, and sorting (Admin only)
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        query    query     dto.CustomerQueryRequest  false  "Query Parameters"
// @Success      200      {object}  dto.WebResponse{data=[]dto.CustomerResponse,meta=dto.PaginationMeta}
// @Failure      400      {object}  dto.WebResponse{data=string}
// @Failure      401      {object}  dto.WebResponse{data=string}
// @Failure      403      {object}  dto.WebResponse{data=string}
// @Failure      500      {object}  dto.WebResponse{data=string}
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /customers [get]
func (h *CustomerHandler) GetAllCustomers(c *gin.Context) {
	var req dto.CustomerQueryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Validasi parameter pencarian gagal",
			Data:    err.Error(),
		})
		return
	}

	customers, meta, err := h.service.GetAllCustomers(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebResponse{
			Success: false,
			Message: "Gagal mengambil data customer",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Data berhasil diambil",
		Data:    customers,
		Meta:    meta,
	})
}

// GetCustomerByID godoc
// @Summary      Get customer by ID
// @Description  Get a customer's details by their ID (Admin or own customer check)
// @Tags         customers
// @Produce      json
// @Param        id   path      int  true  "Customer ID"
// @Success      200  {object}  dto.CustomerResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /customers/{id} [get]
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

// UpdateCustomer godoc
// @Summary      Update customer
// @Description  Update a customer's profile details by ID (Admin or own customer check)
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        id      path      int                         true  "Customer ID"
// @Param        request body      dto.CustomerUpdateRequest  true  "Update Info"
// @Success      200     {object}  dto.CustomerResponse
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Failure      404     {object}  map[string]interface{}
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /customers/{id} [put]
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

// DeleteCustomer godoc
// @Summary      Delete customer
// @Description  Delete a customer record by ID (Admin only)
// @Tags         customers
// @Produce      json
// @Param        id   path      int  true  "Customer ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /customers/{id} [delete]
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
