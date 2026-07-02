package handler

import (
	"go-tiket-konser/dto"
	"go-tiket-konser/service"
	"go-tiket-konser/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user and automatically create a customer record
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RegisterRequest  true  "Registration Info"
// @Success      201      {object}  dto.WebResponse{data=dto.UserResponse}
// @Failure      400      {object}  dto.WebResponse{data=string}
// @Failure      409      {object}  dto.WebResponse{data=string}
// @Failure      500      {object}  dto.WebResponse{data=string}
// @Security     ApiKeyAuth
// @Router       /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Validasi input gagal",
			Data:    err.Error(),
		})
		return
	}

	user, err := h.authService.Register(req.Email, req.Password, req.FullName)
	if err != nil {
		statusCode, errMsg := mapAuthError(err)
		c.JSON(statusCode, dto.WebResponse{
			Success: false,
			Message: "Registrasi gagal",
			Data:    errMsg,
		})
		return
	}

	var customerID uint
	if user.Customer != nil {
		customerID = uint(user.Customer.ID)
	}

	c.JSON(http.StatusCreated, dto.WebResponse{
		Success: true,
		Message: "User registered successfully",
		Data: dto.UserResponse{
			ID:         user.ID,
			FullName:   user.FullName,
			Email:      user.Email,
			Role:       user.Role,
			CustomerID: customerID,
		},
	})

	logger.Log.WithFields(logrus.Fields{
		"email":     user.Email,
		"client_ip": c.ClientIP(),
		"action":    "auth_register_success",
	}).Info("User baru berhasil terdaftar")
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user with email and password to retrieve a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "Login Credentials"
// @Success      200      {object}  dto.WebResponse{data=map[string]interface{}}
// @Failure      400      {object}  dto.WebResponse{data=string}
// @Failure      401      {object}  dto.WebResponse{data=string}
// @Failure      500      {object}  dto.WebResponse{data=string}
// @Security     ApiKeyAuth
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Validasi input gagal",
			Data:    err.Error(),
		})
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		statusCode, errMsg := mapAuthError(err)
		c.JSON(statusCode, dto.WebResponse{
			Success: false,
			Message: "Invalid Email or password",
			Data:    errMsg,
		})
		logger.Log.WithFields(logrus.Fields{
			"email":     req.Email,
			"client_ip": c.ClientIP(),
			"action":    "auth_login_failed",
		}).Warn("Percobaan login gagal: email tidak terdaftar")
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Login berhasil",
		Data: map[string]interface{}{
			"token": token,
		},
	})

	logger.Log.WithFields(logrus.Fields{
		"email":     req.Email,
		"client_ip": c.ClientIP(),
		"action":    "auth_login_success",
	}).Info("User berhasil login")

}

// Logout godoc
// @Summary      Logout user
// @Description  Invalidate the active JWT token by blacklisting it
// @Tags         auth
// @Produce      json
// @Success      200      {object}  dto.WebResponse{}
// @Failure      400      {object}  dto.WebResponse{data=string}
// @Failure      500      {object}  dto.WebResponse{data=string}
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /logout [post]
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
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Token tidak ditemukan",
			Data:    "missing or invalid authorization header",
		})
		return
	}

	err := h.authService.Logout(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebResponse{
			Success: false,
			Message: "Logout gagal",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "User logged out successfully",
	})
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Retrieve profile data of the currently logged-in user
// @Tags         auth
// @Produce      json
// @Success      200      {object}  dto.WebResponse{data=dto.UserResponse}
// @Failure      401      {object}  dto.WebResponse{data=string}
// @Failure      500      {object}  dto.WebResponse{data=string}
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.WebResponse{
			Success: false,
			Message: "Akses tidak sah",
			Data:    "unauthorized",
		})
		return
	}
	userID := userIDVal.(uint)

	user, err := h.authService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebResponse{
			Success: false,
			Message: "Gagal mengambil data profile",
			Data:    err.Error(),
		})
		return
	}

	var customerID uint
	if user.Customer != nil {
		customerID = uint(user.Customer.ID)
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "User profile retrieved successfully",
		Data: dto.UserResponse{
			ID:         user.ID,
			FullName:   user.FullName,
			Email:      user.Email,
			Role:       user.Role,
			CustomerID: customerID,
		},
	})
}

// VerifyOTP godoc
// @Summary      Verify registration OTP
// @Description  Verify the 6-digit OTP code sent to user's email during registration
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.VerifyOTPRequest  true  "Verify OTP Payload"
// @Success      200      {object}  dto.WebResponse{}
// @Failure      400      {object}  dto.WebResponse{data=string}
// @Security     ApiKeyAuth
// @Router       /verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Validasi input gagal",
			Data:    err.Error(),
		})
		return
	}

	err := h.authService.VerifyOTP(req.Email, req.OTP)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Verifikasi OTP gagal",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Verifikasi OTP berhasil, silakan login",
	})

	logger.Log.WithFields(logrus.Fields{
		"email":     req.Email,
		"client_ip": c.ClientIP(),
		"action":    "auth_otp_verified",
	}).Info("Kode OTP sukses diverifikasi, token JWT berhasil diterbitkan")
}
