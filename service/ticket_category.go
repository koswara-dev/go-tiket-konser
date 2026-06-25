package service

import (
	"go-tiket-konser/models"
	"go-tiket-konser/repository"
)

type TicketCategoryService interface {
	CreateTicketCategory(category *models.TicketCategory) error
	GetAllTicketCategories() ([]models.TicketCategory, error)
	GetTicketCategoryByID(id int) (models.TicketCategory, error)
	UpdateTicketCategory(category *models.TicketCategory) error
	DeleteTicketCategory(id int) error
}

type ticketCategoryService struct {
	categoryRepo repository.TicketCategoryRepository
	concertRepo  repository.ConcertRepository
}

func NewTicketCategoryService(
	categoryRepo repository.TicketCategoryRepository,
	concertRepo repository.ConcertRepository,
) TicketCategoryService {
	return &ticketCategoryService{
		categoryRepo: categoryRepo,
		concertRepo:  concertRepo,
	}
}

func (s *ticketCategoryService) CreateTicketCategory(category *models.TicketCategory) error {
	// Verify Concert exists
	_, err := s.concertRepo.FindByID(category.ConcertID)
	if err != nil {
		return err
	}
	return s.categoryRepo.Create(category)
}

func (s *ticketCategoryService) GetAllTicketCategories() ([]models.TicketCategory, error) {
	return s.categoryRepo.FindAll()
}

func (s *ticketCategoryService) GetTicketCategoryByID(id int) (models.TicketCategory, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *ticketCategoryService) UpdateTicketCategory(category *models.TicketCategory) error {
	// Verify Concert exists
	_, err := s.concertRepo.FindByID(category.ConcertID)
	if err != nil {
		return err
	}
	return s.categoryRepo.Update(category)
}

func (s *ticketCategoryService) DeleteTicketCategory(id int) error {
	return s.categoryRepo.Delete(id)
}
