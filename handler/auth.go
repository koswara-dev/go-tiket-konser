package handler

import (
	"go-tiket-konser/dto"
	"go-tiket-konser/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func mapAuthError(err error) (int, string) {
	if err.Error() == "email already registered" {
		return http.StatusConflict, "Email sudah terdaftar"
	}
	if err.Error() == "invalid email" || err.Error() == "invalid password" {
		return http.StatusUnauthorized, "Email atau password salah"
	}
	return http.StatusInternalServerError, err.Error()
}

// Register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi input gagal",
			"error":   err.Error(),
		})
		return
	}

	user, err := h.authService.Register(req.Email, req.Password, req.FullName)
	if err != nil {
		statusCode, errMsg := mapAuthError(err)
		c.JSON(statusCode, gin.H{
			"success": false,
			"message": "Registrasi gagal",
			"error":   errMsg,
		})
		return
	}

	var customerID uint
	if user.Customer != nil {
		customerID = uint(user.Customer.ID)
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully",
		"data": dto.UserResponse{
			ID:         user.ID,
			FullName:   user.FullName,
			Email:      user.Email,
			Role:       user.Role,
			CustomerID: customerID,
		},
	})
}

// Login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi input gagal",
			"error":   err.Error(),
		})
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		statusCode, errMsg := mapAuthError(err)
		c.JSON(statusCode, gin.H{
			"success": false,
			"message": "Login gagal",
			"error":   errMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data": gin.H{
			"token": token,
		},
	})
}

// Logout
func (h *AuthHandler) Logout(c *gin.Context) {
	tokenVal, exists := c.Get("token")
	var token string
	if exists {
		token = tokenVal.(string)
	} else {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 7 {
			token = authHeader[7:]
		}
	}

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Token tidak ditemukan",
			"error":   "missing or invalid authorization header",
		})
		return
	}

	err := h.authService.Logout(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Logout gagal",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User logged out successfully",
	})
}

// Get Profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Akses tidak sah",
			"error":   "unauthorized",
		})
		return
	}
	userID := userIDVal.(uint)

	user, err := h.authService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data profile",
			"error":   err.Error(),
		})
		return
	}

	var customerID uint
	if user.Customer != nil {
		customerID = uint(user.Customer.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User profile retrieved successfully",
		"data": dto.UserResponse{
			ID:         user.ID,
			FullName:   user.FullName,
			Email:      user.Email,
			Role:       user.Role,
			CustomerID: customerID,
		},
	})
}
