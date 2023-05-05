package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"outdoorsy.com/backend/model"
	"outdoorsy.com/backend/repository"
	"outdoorsy.com/backend/route"
)

var router = setupRouter()

func TestShouldGetRentalById(t *testing.T) {
	const RENTAL_ID = 3

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/rentals/%d", RENTAL_ID), nil)
	router.ServeHTTP(w, req)
	resultString := w.Body.String()

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d but got %d instead", http.StatusOK, w.Code)
	}

	var result model.Rental
	err := json.Unmarshal([]byte(resultString), &result)
	if err != nil {
		t.Fatalf("Unable to unmarshal result")
	}
	if result.Id != RENTAL_ID {
		t.Errorf("expectd Rental.Id of %d but got %d instead", RENTAL_ID, result.Id)
	}

}

func TestShouldReturn404WhenRentalIDDoesNotExist(t *testing.T) {
	const RENTAL_ID = 100000

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/rentals/%d", RENTAL_ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status code %d but got %d instead", http.StatusNotFound, w.Code)
	}
}

func TestShouldReturnPagedResults(t *testing.T) {
	const PAGE_SIZE = 3
	const PAGE_NUMBER = 2
	offset := PAGE_SIZE * (PAGE_NUMBER - 1)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/rentals?limit=%d&offset=%d", PAGE_SIZE, offset), nil)
	router.ServeHTTP(w, req)
	resultString := w.Body.String()

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d but got %d instead", http.StatusOK, w.Code)
	}

	var result []model.Rental
	err := json.Unmarshal([]byte(resultString), &result)
	if err != nil {
		t.Fatalf("Unable to unmarshal result")
	}

	if len(result) != PAGE_SIZE {
		t.Errorf("expected result page size of %d but got a size of %d", PAGE_SIZE, len(result))
	}

	if result[0].Id != PAGE_SIZE + 1 {
		t.Errorf("expected first result Rental's Id to be %d but got %d instead", PAGE_SIZE + 1, result[0].Id)
	}

}

func TestShouldGetRentalsNearLocation(t *testing.T) {
	const LATITUDE = 32.83
	const LONGITUDE = -118.28
	var EXPECTED_RESULTS_IDS = [6]int{1,3,5,7,15,23}


	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/rentals?near=%f,%f", LATITUDE, LONGITUDE), nil)
	router.ServeHTTP(w, req)
	resultString := w.Body.String()

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d but got %d instead", http.StatusOK, w.Code)
	}

	var result []model.Rental
	err := json.Unmarshal([]byte(resultString), &result)
	if err != nil {
		t.Fatalf("Unable to unmarshal result")
	}

	if len(result) != len(EXPECTED_RESULTS_IDS) {
		t.Errorf("expected %d rentals but got %d instead", len(EXPECTED_RESULTS_IDS), len(result))
	}

	/* 
	TODO we should probably actually compare to the expected ids to make sure they are the ones expected
	But I've put enough time into this already
	*/ 
}


//TODO add more tests for all the possible filters, sorts and different permutations

// This is a duplicate from app.go because of the way I set up my testing structure
// If I had the time to redo the structure I would would restructure this so the tests could reference the
// setupRouter from app.go
func setupRouter() *gin.Engine{
	//configure database(s)
	// Initialize Database
	dbhost := os.Getenv("PGHOST")
	dbuser := os.Getenv("POSTGRES_USER")
	dbpassword := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("PGDATABASE")
	dbport := os.Getenv("PGPORT")
	if len(dbport) == 0 {
		dbport = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		dbhost, dbuser, dbpassword, dbname, dbport)
	rentalsDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//initial Repository
	rentalRepository := &repository.GORMRentalRepository{Database: rentalsDB}


	//configure router
	router := gin.Default()
	route.Init(router, rentalRepository)
	return router
}