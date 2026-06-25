package service

import (
	"go-tiket-konser/models"
	"go-tiket-konser/repository"
)

type ConcertService interface {
	CreateConcert(concert *models.Concert) error
	GetAllConcerts() ([]models.Concert, error)
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

func (s *concertService) GetAllConcerts() ([]models.Concert, error) {
	return s.repo.FindAll()
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
