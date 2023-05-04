package repository
import (
	"gorm.io/gorm"
	"outdoorsy.com/backend/model"
)

type RentalRepository interface {
	FindById(id int) (model.Rental, error)
	FindAll() ([]model.Rental)
}

type GORMRentalRepository struct {
	Database *gorm.DB
}

func (repository *GORMRentalRepository) FindById(id int) (model.Rental, error) {
	var rental model.Rental
	err := repository.Database.Preload("User").First(&rental, "id = ?", id).Error
	return rental, err
}

func (repository *GORMRentalRepository) FindAll() ([]model.Rental) {
	var rentals []model.Rental
    repository.Database.Preload("User").Find(&rentals)
	return rentals
}