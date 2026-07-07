package service

import (
	"errors"
	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateConcert_Success(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo, nil)

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
	serv := NewConcertService(mockRepo, nil)

	now := time.Now().AddDate(0, 0, 10)
	inputConcert := &models.Concert{
		Title: "Konser A",
	}

	id1 := uuid.New()
	existingConcerts := []models.Concert{
		{
			BaseModel: models.BaseModel{
				ID: id1,
			},
			Title: "Konser A",
			Date:  now,
		},
	}

	mockRepo.On("FindAll").Return(existingConcerts, nil)

	err := serv.CreateConcert(inputConcert)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrConcertAlreadyExists))
	mockRepo.AssertExpectations(t)
}

func TestGetAllConcerts_Success(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo, nil)

	req := dto.ConcertQueryRequest{
		Page:   1,
		Limit:  10,
		Search: "Konser",
		Sort:   "date_asc",
	}

	now := time.Now()
	id1 := uuid.New()
	concerts := []models.Concert{
		{
			BaseModel: models.BaseModel{
				ID: id1,
			},
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
	serv := NewConcertService(mockRepo, nil)

	id1 := uuid.New()
	dummyConcert := models.Concert{
		BaseModel: models.BaseModel{
			ID: id1,
		},
		Title: "Konser A",
	}

	mockRepo.On("FindByID", id1).Return(dummyConcert, nil)

	concert, err := serv.GetConcertByID(id1)

	assert.NoError(t, err)
	assert.Equal(t, "Konser A", concert.Title)
	mockRepo.AssertExpectations(t)
}

func TestGetConcertByID_NotFound(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo, nil)

	idErr := uuid.New()
	mockRepo.On("FindByID", idErr).Return(models.Concert{}, errors.New("not found"))

	_, err := serv.GetConcertByID(idErr)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrConcertNotFound))
	mockRepo.AssertExpectations(t)
}

func TestUpdateConcert_Success(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo, nil)

	id1 := uuid.New()
	dummyConcert := models.Concert{
		BaseModel: models.BaseModel{
			ID: id1,
		},
		Title: "Konser A",
	}

	mockRepo.On("FindByID", id1).Return(dummyConcert, nil)
	mockRepo.On("Update", &dummyConcert).Return(nil)

	err := serv.UpdateConcert(&dummyConcert)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateConcert_NotFound(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo, nil)

	idErr := uuid.New()
	inputConcert := &models.Concert{
		BaseModel: models.BaseModel{
			ID: idErr,
		},
	}

	mockRepo.On("FindByID", idErr).Return(models.Concert{}, errors.New("not found"))

	err := serv.UpdateConcert(inputConcert)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrConcertNotFound))
	mockRepo.AssertExpectations(t)
}

func TestDeleteConcert_Success(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo, nil)

	id1 := uuid.New()
	dummyConcert := models.Concert{
		BaseModel: models.BaseModel{
			ID: id1,
		},
	}

	mockRepo.On("FindByID", id1).Return(dummyConcert, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)
	mockRepo.On("Delete", id1).Return(nil)

	userID := uuid.New()
	err := serv.DeleteConcert(id1, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteConcert_NotFound(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	serv := NewConcertService(mockRepo, nil)

	idErr := uuid.New()
	mockRepo.On("FindByID", idErr).Return(models.Concert{}, errors.New("not found"))

	userID := uuid.New()
	err := serv.DeleteConcert(idErr, userID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrConcertNotFound))
	mockRepo.AssertExpectations(t)
}
