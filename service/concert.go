package service

import (
	"fmt"
	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"go-tiket-konser/repository"
	"math"

	"github.com/google/uuid"
)

type ConcertService interface {
	CreateConcert(concert *models.Concert) error
	GetAllConcerts(req dto.ConcertQueryRequest) ([]dto.ConcertResponse, dto.PaginationMeta, error)
	GetConcertByID(id uuid.UUID) (models.Concert, error)
	UpdateConcert(concert *models.Concert) error
	DeleteConcert(id uuid.UUID, deleterID uuid.UUID) error
}

type concertService struct {
	repo   repository.ConcertRepository
	broker *NotificationBroker
}

func NewConcertService(repo repository.ConcertRepository, broker *NotificationBroker) ConcertService {
	return &concertService{repo: repo, broker: broker}
}

func (s *concertService) CreateConcert(concert *models.Concert) error {
	// validasi bisnis: cek duplikat judul konser
	allConcerts, err := s.repo.FindAll()
	if err == nil {
		for _, c := range allConcerts {
			if c.Title == concert.Title {
				return models.ErrConcertAlreadyExists
			}
		}
	}

	err = s.repo.Create(concert)
	if err == nil && s.broker != nil {
		msg := fmt.Sprintf("Konser baru '%s' telah ditambahkan di %s pada %s!", concert.Title, concert.Venue, concert.Date.Format("2006-01-02"))
		_ = s.broker.SendNotification("", "Konser Baru!", msg)
	}
	return err
}

func (s *concertService) GetAllConcerts(req dto.ConcertQueryRequest) ([]dto.ConcertResponse, dto.PaginationMeta, error) {
	// 1. Validasi & Set Nilai Default untuk Limit
	limit := req.Limit
	if limit <= 0 {
		limit = 10 // Default limit per halaman
	}
	if limit > 100 {
		limit = 100 // Batasan keamanan maksimum limit demi memitigasi serangan DOS
	}

	// 2. Validasi & Set Nilai Default untuk Page
	page := req.Page
	if page <= 0 {
		page = 1 // Default halaman pertama
	}

	// 3. Hitung Kalkulasi Offset
	offset := (page - 1) * limit

	// 4. Panggil Kueri ke Repositori
	concerts, totalData, err := s.repo.FindAllPaginated(req.Search, limit, offset, req.Sort)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	// 5. Transformasi Entity Model Database ke dalam bentuk DTO Response
	var responses []dto.ConcertResponse = make([]dto.ConcertResponse, 0)
	for _, c := range concerts {
		responses = append(responses, dto.ConcertResponse{
			ID:           c.ID,
			Title:        c.Title,
			Description:  c.Description,
			Date:         c.Date.Format("2006-01-02"),
			Venue:        c.Venue,
			Status:       c.Status,
			PosterURL:    c.PosterURL,
			ThumbnailURL: c.ThumbnailURL,
			RulesPDFURL:  c.RulesPDFURL,
			CreatedAt:    c.CreatedAt,
			UpdatedAt:    c.UpdatedAt,
		})
	}

	// 6. Hitung Total Halaman (Membulatkan ke Atas)
	totalPage := 0
	if totalData > 0 {
		totalPage = int(math.Ceil(float64(totalData) / float64(limit)))
	}

	// 7. Bentuk metadata pagination
	meta := dto.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: totalData,
		TotalPage: totalPage,
	}

	return responses, meta, nil
}

func (s *concertService) GetConcertByID(id uuid.UUID) (models.Concert, error) {
	concert, err := s.repo.FindByID(id)
	if err != nil {
		return models.Concert{}, models.ErrConcertNotFound
	}
	return concert, nil
}

func (s *concertService) UpdateConcert(concert *models.Concert) error {
	_, err := s.repo.FindByID(concert.ID)
	if err != nil {
		return models.ErrConcertNotFound
	}
	return s.repo.Update(concert)
}

func (s *concertService) DeleteConcert(id uuid.UUID, deleterID uuid.UUID) error {
	concert, err := s.repo.FindByID(id)
	if err != nil {
		return models.ErrConcertNotFound
	}
	concert.DeletedBy = &deleterID
	_ = s.repo.Update(&concert)
	return s.repo.Delete(id)
}
