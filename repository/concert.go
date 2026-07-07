package repository

import (
	"go-tiket-konser/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConcertRepository interface {
	Create(concert *models.Concert) error
	FindAll() ([]models.Concert, error)
	FindAllPaginated(search string, limit int, offset int, sort string) ([]models.Concert, int64, error)
	FindByID(id uuid.UUID) (models.Concert, error)
	Update(concert *models.Concert) error
	Delete(id uuid.UUID) error
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

func (r *concertRepository) FindByID(id uuid.UUID) (models.Concert, error) {
	var concert models.Concert
	err := r.db.First(&concert, id).Error
	return concert, err
}

func (r *concertRepository) Update(concert *models.Concert) error {
	return r.db.Save(concert).Error
}

func (r *concertRepository) Delete(id uuid.UUID) error {
	var concert models.Concert
	if err := r.db.First(&concert, id).Error; err != nil {
		return err
	}
	return r.db.Delete(&concert).Error
}

func (r *concertRepository) FindAllPaginated(search string, limit int, offset int, sort string) ([]models.Concert, int64, error) {
	var concerts []models.Concert
	var total int64

	query := r.db.Model(&models.Concert{})

	if search != "" {
		query = query.Where("title ILIKE ? OR venue ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch sort {
	case "date_asc":
		query = query.Order("date ASC")
	case "date_desc":
		query = query.Order("date DESC")
	case "title_asc":
		query = query.Order("title ASC")
	case "title_desc":
		query = query.Order("title DESC")
	default:
		query = query.Order("id DESC")
	}

	err := query.Limit(limit).Offset(offset).Find(&concerts).Error
	return concerts, total, err
}
