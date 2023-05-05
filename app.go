package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"outdoorsy.com/backend/route"
	"outdoorsy.com/backend/repository"
)

func main() {
	router := setupRouter()
	// get application host and port from env so they can be configured for each environment in necessary
	host := os.Getenv("APP_HOST")// Defaults to all network interfaces
	port := os.Getenv("APP_PORT")
	if len(port) == 0 {
		port = "8080" //Default to 8080
	}
	router.Run(fmt.Sprintf("%s:%s", host, port)) // TODO: make this configurable

}

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
