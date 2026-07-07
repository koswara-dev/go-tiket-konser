package handler

import (
	"bytes"
	"encoding/json"
	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"go-tiket-konser/utils/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func init() {
	logger.Log = logrus.New()
	logger.Log.SetOutput(io.Discard)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("future_date", func(fl validator.FieldLevel) bool {
			dateStr, ok := fl.Field().Interface().(string)
			if !ok {
				return false
			}
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return false
			}
			return date.After(time.Now().AddDate(0, 0, 7))
		})
	}
}

func setTestUserContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", uuid.MustParse("99999999-9999-9999-9999-999999999999"))
		c.Set("customer_id", uuid.MustParse("11111111-1111-1111-1111-111111111111"))
		c.Set("role", "admin")
		c.Next()
	}
}

func TestCreateConcert_Handler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(setTestUserContext())

	mockServ := new(MockConcertService)
	mockSP := new(MockStorageProvider)
	h := NewConcertHandler(mockServ, mockSP)

	r.POST("/concerts", h.CreateConcert)

	// date is 10 days in the future to pass validation
	futureDateStr := time.Now().AddDate(0, 0, 10).Format("2006-01-02")
	reqBody := dto.ConcertRequest{
		Title:       "Konser Musik Hebat",
		Description: "Konser asik sekali",
		Date:        futureDateStr,
		Venue:       "Stadium Utama Jakarta",
		Status:      "active",
	}

	concertID := uuid.New()
	mockServ.On("CreateConcert", mock.AnythingOfType("*models.Concert")).Return(nil).Run(func(args mock.Arguments) {
		c := args.Get(0).(*models.Concert)
		c.ID = concertID
	})

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/concerts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp["status"].(bool))
	assert.Equal(t, "Konser berhasil ditambahkan", resp["message"])

	mockServ.AssertExpectations(t)
}

func TestGetConcerts_Handler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockServ := new(MockConcertService)
	mockSP := new(MockStorageProvider)
	h := NewConcertHandler(mockServ, mockSP)

	r.GET("/concerts", h.GetConcerts)

	query := dto.ConcertQueryRequest{
		Page:  1,
		Limit: 10,
	}

	id1 := uuid.New()
	concertsResp := []dto.ConcertResponse{
		{
			ID:    id1,
			Title: "Konser Musik Hebat",
		},
	}
	meta := dto.PaginationMeta{
		Page:      1,
		Limit:     10,
		TotalData: 1,
		TotalPage: 1,
	}

	mockServ.On("GetAllConcerts", query).Return(concertsResp, meta, nil)

	req, _ := http.NewRequest(http.MethodGet, "/concerts?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.WebResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)

	mockServ.AssertExpectations(t)
}

func TestGetConcertByID_Handler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockServ := new(MockConcertService)
	mockSP := new(MockStorageProvider)
	h := NewConcertHandler(mockServ, mockSP)

	r.GET("/concerts/:id", h.GetConcertByID)

	id1 := uuid.New()
	dummyConcert := models.Concert{
		BaseModel: models.BaseModel{
			ID: id1,
		},
		Title: "Konser Musik Hebat",
		Date:  time.Now(),
	}

	mockServ.On("GetConcertByID", id1).Return(dummyConcert, nil)

	req, _ := http.NewRequest(http.MethodGet, "/concerts/"+id1.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp["success"].(bool))

	mockServ.AssertExpectations(t)
}

func TestGetConcertByID_Handler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockServ := new(MockConcertService)
	mockSP := new(MockStorageProvider)
	h := NewConcertHandler(mockServ, mockSP)

	r.GET("/concerts/:id", h.GetConcertByID)

	idErr := uuid.New()
	mockServ.On("GetConcertByID", idErr).Return(models.Concert{}, gorm.ErrRecordNotFound)

	req, _ := http.NewRequest(http.MethodGet, "/concerts/"+idErr.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp["success"].(bool))
	assert.Equal(t, "Concert not found", resp["message"])

	mockServ.AssertExpectations(t)
}

func TestUpdateConcert_Handler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(setTestUserContext())

	mockServ := new(MockConcertService)
	mockSP := new(MockStorageProvider)
	h := NewConcertHandler(mockServ, mockSP)

	r.PUT("/concerts/:id", h.UpdateConcert)

	id1 := uuid.New()
	dummyConcert := models.Concert{
		BaseModel: models.BaseModel{
			ID: id1,
		},
		Title: "Old Title",
	}

	mockServ.On("GetConcertByID", id1).Return(dummyConcert, nil)
	mockServ.On("UpdateConcert", mock.AnythingOfType("*models.Concert")).Return(nil)

	futureDate := time.Now().AddDate(0, 0, 10)
	inputConcert := models.Concert{
		Title:        "New Title",
		Description:  "New Description",
		Date:         futureDate,
		Venue:        "New Venue",
		Status:       "active",
		PosterURL:    "http://example.com/poster.jpg",
		ThumbnailURL: "http://example.com/thumbnail.jpg",
		RulesPDFURL:  "http://example.com/rules.pdf",
	}

	body, _ := json.Marshal(inputConcert)
	req, _ := http.NewRequest(http.MethodPut, "/concerts/"+id1.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockServ.AssertExpectations(t)
}

func TestDeleteConcert_Handler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(setTestUserContext())

	mockServ := new(MockConcertService)
	mockSP := new(MockStorageProvider)
	h := NewConcertHandler(mockServ, mockSP)

	r.DELETE("/concerts/:id", h.DeleteConcert)

	id1 := uuid.New()
	adminUserID := uuid.MustParse("99999999-9999-9999-9999-999999999999")
	mockServ.On("DeleteConcert", id1, adminUserID).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/concerts/"+id1.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Concert deleted successfully", resp["message"])

	mockServ.AssertExpectations(t)
}
