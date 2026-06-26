package handler

import (
	"net/http"
	"strconv"

	"go-tiket-konser/dto"
	"go-tiket-konser/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var res []dto.UserResponse
	for _, user := range users {
		var customerID uint
		if user.Customer != nil {
			customerID = uint(user.Customer.ID)
		}
		res = append(res, dto.UserResponse{
			ID:         user.ID,
			FullName:   user.FullName,
			Email:      user.Email,
			Role:       user.Role,
			CustomerID: customerID,
		})
	}
	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	requestedID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	tokenUserID := c.MustGet("user_id").(uint)
	tokenRole := c.MustGet("role").(string)

	if tokenRole != "admin" && tokenUserID != uint(requestedID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: tidak dapat mengakses data user lain"})
		return
	}

	user, err := h.service.GetUserByID(uint(requestedID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	var customerID uint
	if user.Customer != nil {
		customerID = uint(user.Customer.ID)
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:         user.ID,
		FullName:   user.FullName,
		Email:      user.Email,
		Role:       user.Role,
		CustomerID: customerID,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	requestedID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	tokenUserID := c.MustGet("user_id").(uint)
	tokenRole := c.MustGet("role").(string)

	if tokenRole != "admin" && tokenUserID != uint(requestedID) {
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

	user, err := h.service.UpdateUser(uint(requestedID), &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var customerID uint
	if user.Customer != nil {
		customerID = uint(user.Customer.ID)
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:         user.ID,
		FullName:   user.FullName,
		Email:      user.Email,
		Role:       user.Role,
		CustomerID: customerID,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	requestedID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	err = h.service.DeleteUser(uint(requestedID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dihapus"})
}
