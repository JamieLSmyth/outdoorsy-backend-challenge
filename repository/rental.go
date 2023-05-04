package repository

import (
	"gorm.io/gorm"
	"outdoorsy.com/backend/model"
)

const MILE_IN_METERS = 1609.34

type LatLong struct {
	Latitude  float64
	Longitude float64
}

type RentalFilter struct {
	PriceMax *float64
	PriceMin *float64
	IDs      *[]string //TODO this should be a list of ints
	Near     *LatLong
}

type RentalRepository interface {
	FindById(id int) (model.Rental, error)
	FindAllByFilter(RentalFilter) []model.Rental
}

type GORMRentalRepository struct {
	Database *gorm.DB
}

func (repository *GORMRentalRepository) FindById(id int) (model.Rental, error) {
	var rental model.Rental
	err := repository.Database.Preload("User").First(&rental, "id = ?", id).Error
	return rental, err
}

func (repository *GORMRentalRepository) FindAllByFilter(filter RentalFilter, offset int, limit int, sort []string) []model.Rental {
	var rentals []model.Rental
	query := repository.Database.Preload("User").Offset(offset)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if filter.PriceMin != nil {
		query = query.Where("price_per_day > ?", *filter.PriceMin)
	}
	if filter.PriceMax != nil {
		query = query.Where("price_per_day < ?", *filter.PriceMax)
	}
	if filter.IDs != nil && len(*filter.IDs) > 0 {
		query = query.Where("id IN (?)", *filter.IDs)
	}
	if filter.Near !=nil {
		query = query.Where(`
		ST_Distance_Sphere(
			ST_MakePoint(lng, lat),
			ST_MakePoint(?, ?)
		) <= ?
		`, filter.Near.Longitude, filter.Near.Latitude, MILE_IN_METERS * 100)
	}
	query.Find(&rentals)
	return rentals
}
