package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-tiket-konser/dto"
	"go-tiket-konser/utils/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Initialize a dummy logger to avoid nil pointer dereference during testing
	logger.Log = logrus.New()
	logger.Log.SetOutput(io.Discard)
}

func TestHandlerLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockAuthService := new(MockAuthService)
	h := NewAuthHandler(mockAuthService)
	r.POST("/login", h.Login)

	loginReq := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockAuthService.On("Login", loginReq.Email, loginReq.Password).Return("mocked-jwt-token", nil)

	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.WebResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Login berhasil", response.Message)

	dataMap := response.Data.(map[string]interface{})
	assert.Equal(t, "mocked-jwt-token", dataMap["token"])

	mockAuthService.AssertExpectations(t)
}

func TestHandlerLogin_ValidationFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockAuthService := new(MockAuthService)
	h := NewAuthHandler(mockAuthService)
	r.POST("/login", h.Login)

	// Password too short (min=6) or empty email
	loginReq := dto.LoginRequest{
		Email:    "invalid-email",
		Password: "123",
	}

	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.WebResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Validasi input gagal", response.Message)

	mockAuthService.AssertNotCalled(t, "Login")
}

func TestHandlerLogin_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockAuthService := new(MockAuthService)
	h := NewAuthHandler(mockAuthService)
	r.POST("/login", h.Login)

	loginReq := dto.LoginRequest{
		Email:    "wrong@example.com",
		Password: "password123",
	}

	mockAuthService.On("Login", loginReq.Email, loginReq.Password).Return("", errors.New("invalid email"))

	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.WebResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Email atau password salah", response.Data)

	mockAuthService.AssertExpectations(t)
}
