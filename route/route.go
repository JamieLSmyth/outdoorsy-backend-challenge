package route

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

    db.Config.Logger = logger.Default.LogMode(logger.Info)

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
	filter := repository.RentalFilter{}
	if value, err := strconv.ParseFloat(context.Query("price_min"), 64); err == nil {
		filter.PriceMin = &value
	}
	if value, err := strconv.ParseFloat(context.Query("price_max"), 64); err == nil {
		filter.PriceMax = &value
	}
	ids_string := context.Query("ids")
	if len(ids_string) > 0 { // TODO Shouldn't need this once doing the convert below
		ids := strings.Split(ids_string, ",")
		//TODO convert these to integers and error ignore anything that is not an int
		filter.IDs = &ids
	}
    near := strings.Split(context.Query("near"), ",")
    if len(near) == 2 {
        latLong := repository.LatLong{}
        if value, err := strconv.ParseFloat(near[0],64); err == nil {
            latLong.Latitude = value
        }
        if value, err := strconv.ParseFloat(near[1],64); err == nil {
            latLong.Longitude = value
        }
        filter.Near = &latLong
    }

	context.IndentedJSON(http.StatusOK, RentalRepository.FindAllByFilter(filter))
}
