package handler

import (
	"net/http"

	"go-tiket-konser/dto"
	"go-tiket-konser/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Get list of all users with search, pagination, and sorting (Admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        query    query     dto.UserQueryRequest  false  "Query Parameters"
// @Success      200      {object}  dto.WebResponse{data=[]dto.UserResponse,meta=dto.PaginationMeta}
// @Failure      400      {object}  dto.WebResponse{data=string}
// @Failure      401      {object}  dto.WebResponse{data=string}
// @Failure      403      {object}  dto.WebResponse{data=string}
// @Failure      500      {object}  dto.WebResponse{data=string}
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var req dto.UserQueryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Validasi parameter pencarian gagal",
			Data:    err.Error(),
		})
		return
	}

	users, meta, err := h.service.GetAllUsers(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebResponse{
			Success: false,
			Message: "Gagal mengambil data user",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Data berhasil diambil",
		Data:    users,
		Meta:    meta,
	})
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Get a user's details by their ID (Admin or own user check)
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	requestedID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	tokenUserID := c.MustGet("user_id").(uuid.UUID)
	tokenRole := c.MustGet("role").(string)

	if tokenRole != "admin" && tokenUserID != requestedID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: tidak dapat mengakses data user lain"})
		return
	}

	user, err := h.service.GetUserByID(requestedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	var customerID uuid.UUID
	if user.Customer != nil {
		customerID = user.Customer.ID
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:         user.ID,
		FullName:   user.FullName,
		Email:      user.Email,
		Role:       user.Role,
		CustomerID: customerID,
	})
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update a user's profile information by ID (Admin or own user check)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id      path      string                  true  "User ID"
// @Param        request body      dto.UserUpdateRequest  true  "Update Info"
// @Success      200     {object}  dto.UserResponse
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Failure      404     {object}  map[string]interface{}
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	requestedID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	tokenUserID := c.MustGet("user_id").(uuid.UUID)
	tokenRole := c.MustGet("role").(string)

	if tokenRole != "admin" && tokenUserID != requestedID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: tidak dapat memperbarui data user lain"})
		return
	}

	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Customer cannot change their own role
	if tokenRole != "admin" {
		req.Role = "customer"
	}

	user, err := h.service.UpdateUser(requestedID, &req, tokenUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var customerID uuid.UUID
	if user.Customer != nil {
		customerID = user.Customer.ID
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:         user.ID,
		FullName:   user.FullName,
		Email:      user.Email,
		Role:       user.Role,
		CustomerID: customerID,
	})
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete a user record by ID (Admin only)
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	requestedID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	tokenUserID := c.MustGet("user_id").(uuid.UUID)

	err = h.service.DeleteUser(requestedID, tokenUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dihapus"})
}
