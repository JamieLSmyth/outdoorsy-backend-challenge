package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"outdoorsy.com/backend/route"
	"outdoorsy.com/backend/repository"
)

func main() {

	//configure database(s)
	// Initialize Database
	dsn := "host=postgres user=root password=root dbname=testingwithrentals port=5432 sslmode=disable"
	rentalsDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//initial Repository
	rentalRepository := &repository.GORMRentalRepository{Database: rentalsDB}


	//This is just in here for debuging
    rentalsDB.Config.Logger = logger.Default.LogMode(logger.Info)

	//configure router
	router := gin.Default()
	route.Init(router, rentalRepository)
	// get application host and port from env so they can be configured for each environment in necessary
	host := os.Getenv("APP_HOST")// Defaults to all network interfaces
	port := os.Getenv("APP_PORT")
	if len(port) == 0 {
		port = "8080" //Default to 8080
	}
	router.Run(fmt.Sprintf("%s:%s", host, port)) // TODO: make this configurable

}