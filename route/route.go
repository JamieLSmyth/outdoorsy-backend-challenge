package route

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "outdoorsy.com/backend/repository"
)

var RentalRepository *repository.GORMRentalRepository = nil

func Init(router *gin.Engine) {
    // Initialize Database
	dsn := "host=postgres user=root password=root dbname=testingwithrentals port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }
    //Initialize Repositories
    RentalRepository = &repository.GORMRentalRepository{Database: db}


    router.GET("/rentals", GetRentals)
	router.GET("/rentals/:id", GetRental)
}


func GetRental(context *gin.Context) {
    id, err := strconv.Atoi(context.Param("id"))
    if err != nil {
        context.AbortWithStatus(http.StatusNotFound)
        return
    }
    rental, err := RentalRepository.FindById(id)
    if err != nil {
        context.AbortWithStatus(http.StatusNotFound)
        return
    }
    context.IndentedJSON(http.StatusOK, rental)
}

func GetRentals(context *gin.Context) {
    context.IndentedJSON(http.StatusOK, RentalRepository.FindAll())
}