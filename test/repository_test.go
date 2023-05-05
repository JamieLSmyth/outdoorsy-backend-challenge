package test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"outdoorsy.com/backend/repository"
)

func TestShouldFindExistingRentalByID(t *testing.T){
	const A_RENTAL_ID = 1

	rentalRepository, mock := createRepositoryAndSQLMock(t)

	expected := mock.NewRows([]string{"id"}).AddRow(A_RENTAL_ID)
	mock.ExpectQuery("SELECT \\* FROM \"rentals\" WHERE id = \\$1").
		WithArgs(A_RENTAL_ID).
		WillReturnRows(expected)


	result, _ := rentalRepository.FindById(A_RENTAL_ID)
	if result.Id != 1 {
		t.Errorf("expected Rental.Id to be %d but was %d", A_RENTAL_ID, result.Id)
	}
	
	if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("failed to meet mock expectations: %s", err)
    }
}

func TestShouldErrorWhenRentalDoesNotExistForID(t *testing.T){
	const A_RENTAL_ID = 1
	
	rentalRepository, mock := createRepositoryAndSQLMock(t)

	mock.ExpectQuery("SELECT \\* FROM \"rentals\" WHERE id = \\$1").
		WithArgs(A_RENTAL_ID).
		WillReturnError(fmt.Errorf("Entity Not Found"))


	_, err := rentalRepository.FindById(A_RENTAL_ID)
	if err == nil {
		t.Errorf("Expected error but none was returned")
	}
	
	if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("failed to meet mock expectations: %s", err)
    }
}

func TestShouldGetAllRentalsWhenNoFiltersSupplied(t *testing.T) {
	const RENTALS_COUNT = 3 //This is currently a magic number but should be used to generate number of result rows

	rentalRepository, mock := createRepositoryAndSQLMock(t)

	expected := mock.NewRows([]string{"id"}).
		AddRow(1).
		AddRow(2).
		AddRow(3)
	mock.ExpectQuery("SELECT \\* FROM \"rentals\"").
	WillReturnRows(expected)

	result, _ := rentalRepository.FindAllByFilter(repository.RentalFilter{}, 0, 0, "")
	if len(result) != RENTALS_COUNT {
		t.Errorf("expected %d rentals but got %d", RENTALS_COUNT, len(result))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("failed to meet mock expectations: %s", err)
    }

}

func TestShouldFilterByProximityWhenNearFilterIsSupplied(t *testing.T) {
	const RENTALS_COUNT = 3 //This is currently a magic number but should be used to generate number of result rows
	const FILTER_LATITUDE = 10.03 //This could be randomly generated
	const FILTER_LONGITUDE = -9.12 //Ths could be randomly generated
	const ONE_HUNDRED_MILES_IN_METERS = 160934.0
	
	rentalRepository, mock := createRepositoryAndSQLMock(t)

	expected := mock.NewRows([]string{"id"}).
		AddRow(1).
		AddRow(2).
		AddRow(3)
	mock.ExpectQuery("SELECT \\* FROM \"rentals\" WHERE ST_Distance_Sphere\\( ST_MakePoint\\(lng, lat\\), ST_MakePoint\\(\\$1, \\$2\\) \\) <= \\$3").
	WithArgs(FILTER_LONGITUDE, FILTER_LATITUDE, ONE_HUNDRED_MILES_IN_METERS).
	WillReturnRows(expected)

	result, _ := rentalRepository.FindAllByFilter(repository.RentalFilter{Near: &repository.LatLong{Latitude: FILTER_LATITUDE, Longitude: FILTER_LONGITUDE}}, 0, 0, "")
	if len(result) != RENTALS_COUNT {
		t.Errorf("expected %d rentals but got %d", RENTALS_COUNT, len(result))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("failed to meet mock expectations: %s", err)
    }
}

//TODO Implement the many other test for filters, limit, offset and different permutations

func createRepositoryAndSQLMock(t *testing.T) (repository.RentalRepository, sqlmock.Sqlmock)  {
	db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error occured creating mock database")
    }

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: db,
    }), &gorm.Config{})
	if err !=nil {
		t.Fatalf("error occured creating gorm database")
	}

	repository := &repository.GORMRentalRepository{Database: gormDB}

	return repository, mock
}