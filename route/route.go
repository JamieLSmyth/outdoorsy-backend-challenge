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

type GeoLocation struct {
	Latitude float64
	Longitude float64
}

type RentalQueryParams struct {
	PriceMin *float64 `form:"price_min"`
	PriceMax *float64 `form:"price_max"`
	IDs *string `form:"price_max"`
	Near *string `form:"near"`
	Limit *int `form:"limit"`
	Offset *int `form:"offset"`
	Sort *string `form:"sort"`
}

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
	var query = RentalQueryParams{}
	if err := context.ShouldBindQuery(&query); err != nil {
        context.AbortWithError(http.StatusBadRequest, err)
		return
    }

	filter := repository.RentalFilter{}
	filter.PriceMin = query.PriceMin
	filter.PriceMax = query.PriceMax
	if query.IDs != nil && len(*query.IDs) > 0 { // TODO Shouldn't need this once doing the convert below
		ids := strings.Split(*query.IDs, ",")
		//TODO convert these to integers and error ignore anything that is not an int
		filter.IDs = &ids
	}
	if query.Near != nil {
    	near := strings.Split(*query.Near, ",")
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
	}

	context.IndentedJSON(http.StatusOK, RentalRepository.FindAllByFilter(filter))
}
