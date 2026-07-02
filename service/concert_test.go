package service

import (
	"errors"
	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateConcert_Success(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo)

	now := time.Now().AddDate(0, 0, 10)
	inputConcert := &models.Concert{
		Title:       "Konser A",
		Description: "Deskripsi",
		Date:        now,
		Venue:       "Stadion Utama",
		Status:      "active",
	}

	// mock FindAll to return nil (no duplicate check fails)
	mockRepo.On("FindAll").Return([]models.Concert{}, nil)
	mockRepo.On("Create", inputConcert).Return(nil)

	err := serv.CreateConcert(inputConcert)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateConcert_DuplicateTitle(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo)

	now := time.Now().AddDate(0, 0, 10)
	inputConcert := &models.Concert{
		Title: "Konser A",
	}

	existingConcerts := []models.Concert{
		{ID: 1, Title: "Konser A", Date: now},
	}

	mockRepo.On("FindAll").Return(existingConcerts, nil)

	err := serv.CreateConcert(inputConcert)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrConcertAlreadyExists))
	mockRepo.AssertExpectations(t)
}

func TestGetAllConcerts_Success(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo)

	req := dto.ConcertQueryRequest{
		Page:   1,
		Limit:  10,
		Search: "Konser",
		Sort:   "date_asc",
	}

	now := time.Now()
	concerts := []models.Concert{
		{
			ID:          1,
			Title:       "Konser A",
			Description: "Deskripsi A",
			Date:        now,
			Venue:       "Venue A",
			Status:      "active",
		},
	}

	mockRepo.On("FindAllPaginated", "Konser", 10, 0, "date_asc").Return(concerts, int64(1), nil)

	responses, meta, err := serv.GetAllConcerts(req)

	assert.NoError(t, err)
	assert.Len(t, responses, 1)
	assert.Equal(t, "Konser A", responses[0].Title)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 10, meta.Limit)
	assert.Equal(t, int64(1), meta.TotalData)
	assert.Equal(t, 1, meta.TotalPage)
	mockRepo.AssertExpectations(t)
}

func TestGetConcertByID_Success(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo)

	dummyConcert := models.Concert{
		ID:    1,
		Title: "Konser A",
	}

	mockRepo.On("FindByID", 1).Return(dummyConcert, nil)

	concert, err := serv.GetConcertByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "Konser A", concert.Title)
	mockRepo.AssertExpectations(t)
}

func TestGetConcertByID_NotFound(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo)

	mockRepo.On("FindByID", 99).Return(models.Concert{}, errors.New("not found"))

	_, err := serv.GetConcertByID(99)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrConcertNotFound))
	mockRepo.AssertExpectations(t)
}

func TestUpdateConcert_Success(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo)

	dummyConcert := models.Concert{
		ID:    1,
		Title: "Konser A",
	}

	mockRepo.On("FindByID", 1).Return(dummyConcert, nil)
	mockRepo.On("Update", &dummyConcert).Return(nil)

	err := serv.UpdateConcert(&dummyConcert)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateConcert_NotFound(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo)

	inputConcert := &models.Concert{
		ID: 99,
	}

	mockRepo.On("FindByID", 99).Return(models.Concert{}, errors.New("not found"))

	err := serv.UpdateConcert(inputConcert)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrConcertNotFound))
	mockRepo.AssertExpectations(t)
}

func TestDeleteConcert_Success(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo)

	dummyConcert := models.Concert{
		ID: 1,
	}

	mockRepo.On("FindByID", 1).Return(dummyConcert, nil)
	mockRepo.On("Delete", 1).Return(nil)

	err := serv.DeleteConcert(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteConcert_NotFound(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo)

	mockRepo.On("FindByID", 99).Return(models.Concert{}, errors.New("not found"))

	err := serv.DeleteConcert(99)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrConcertNotFound))
	mockRepo.AssertExpectations(t)
}
