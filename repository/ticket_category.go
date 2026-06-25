package repository

import (
	"go-tiket-konser/models"

	"gorm.io/gorm"
)

type TicketCategoryRepository interface {
	Create(category *models.TicketCategory) error
	FindAll() ([]models.TicketCategory, error)
	FindByID(id int) (models.TicketCategory, error)
	Update(category *models.TicketCategory) error
	Delete(id int) error
}

type ticketCategoryRepository struct {
	db *gorm.DB
}

func NewTicketCategoryRepository(db *gorm.DB) TicketCategoryRepository {
	return &ticketCategoryRepository{db: db}
}

func (r *ticketCategoryRepository) Create(category *models.TicketCategory) error {
	return r.db.Create(category).Error
}

func (r *ticketCategoryRepository) FindAll() ([]models.TicketCategory, error) {
	var categories []models.TicketCategory
	err := r.db.Preload("Concert").Find(&categories).Error
	return categories, err
}

func (r *ticketCategoryRepository) FindByID(id int) (models.TicketCategory, error) {
	var category models.TicketCategory
	err := r.db.Preload("Concert").First(&category, id).Error
	return category, err
}

func (r *ticketCategoryRepository) Update(category *models.TicketCategory) error {
	return r.db.Save(category).Error
}

func (r *ticketCategoryRepository) Delete(id int) error {
	var category models.TicketCategory
	if err := r.db.First(&category, id).Error; err != nil {
		return err
	}
	return r.db.Delete(&category).Error
}
