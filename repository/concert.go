package repository

import (
	"go-tiket-konser/models"

	"gorm.io/gorm"
)

type ConcertRepository interface {
	Create(concert *models.Concert) error
	FindAll() ([]models.Concert, error)
	FindByID(id int) (models.Concert, error)
	Update(concert *models.Concert) error
	Delete(id int) error
}

type concertRepository struct {
	db *gorm.DB
}

func NewConcertRepository(db *gorm.DB) ConcertRepository {
	return &concertRepository{db: db}
}

func (r *concertRepository) Create(concert *models.Concert) error {
	return r.db.Create(concert).Error
}

func (r *concertRepository) FindAll() ([]models.Concert, error) {
	var concerts []models.Concert
	err := r.db.Find(&concerts).Error
	return concerts, err
}

func (r *concertRepository) FindByID(id int) (models.Concert, error) {
	var concert models.Concert
	err := r.db.First(&concert, id).Error
	return concert, err
}

func (r *concertRepository) Update(concert *models.Concert) error {
	return r.db.Save(concert).Error
}

func (r *concertRepository) Delete(id int) error {
	var concert models.Concert
	if err := r.db.First(&concert, id).Error; err != nil {
		return err
	}
	return r.db.Delete(&concert).Error
}
