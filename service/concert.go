package service

import (
	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"go-tiket-konser/repository"
	"math"
)

type ConcertService interface {
	CreateConcert(concert *models.Concert) error
	GetAllConcerts(req dto.ConcertQueryRequest) ([]dto.ConcertResponse, dto.PaginationMeta, error)
	GetConcertByID(id int) (models.Concert, error)
	UpdateConcert(concert *models.Concert) error
	DeleteConcert(id int) error
}

type concertService struct {
	repo repository.ConcertRepository
}

func NewConcertService(repo repository.ConcertRepository) ConcertService {
	return &concertService{repo: repo}
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

	return s.repo.Create(concert)
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
	var responses []dto.ConcertResponse
	for _, c := range concerts {
		responses = append(responses, dto.ConcertResponse{
			ID:          c.ID,
			Title:       c.Title,
			Description: c.Description,
			Date:        c.Date.Format("2006-01-02"),
			Venue:       c.Venue,
			Status:      c.Status,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
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

func (s *concertService) GetConcertByID(id int) (models.Concert, error) {
	concert, err := s.repo.FindByID(id)
	if err != nil {
		return models.Concert{}, models.ErrConcertNotFound
	}
	return concert, nil
}

func (s *concertService) UpdateConcert(concert *models.Concert) error {
	_, err := s.repo.FindByID(int(concert.ID))
	if err != nil {
		return models.ErrConcertNotFound
	}
	return s.repo.Update(concert)
}

func (s *concertService) DeleteConcert(id int) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return models.ErrConcertNotFound
	}
	return s.repo.Delete(id)
}
