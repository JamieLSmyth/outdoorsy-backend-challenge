package route

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"outdoorsy.com/backend/repository"
)

type RentalQueryParams struct {
	PriceMin *float64 `form:"price_min"`
	PriceMax *float64 `form:"price_max"`
	IDs string `form:"ids"`
	Near string `form:"near"`
	Limit int `form:"limit"`
	Offset int `form:"offset"`
	Sort string `form:"sort"`
}

/* TODO there is probably a better way to do this so the names match the json out put 
and the fields are automatically mapped through the ORM */ 
var SORT_FIELD_MAP = map[string]string{
	"price":"price_per_day",
	"name":"name",
	"type":"type",
	"city":"home_city",
	"state":"home_state",
	"country":"home_country",
}

var RentalRepository repository.RentalRepository = nil

func Init(router *gin.Engine, repository repository.RentalRepository) {

	RentalRepository = repository

	router.GET("/rentals", getRentals)
	router.GET("/rentals/:id", getRental)
}

func getRental(context *gin.Context) {
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

func getRentals(context *gin.Context) {
	var query = RentalQueryParams{}
	if err := context.ShouldBindQuery(&query); err != nil {
        context.AbortWithError(http.StatusBadRequest, err)
		return
    }

	filter := repository.RentalFilter{}
	filter.PriceMin = query.PriceMin
	filter.PriceMax = query.PriceMax

	//TODO this could get way fancier by handling trailing commas, empty entries and lots of other things
	if len(query.IDs) > 0 {
		idStrings := strings.Split(query.IDs, ",")
		ids := make([]int, len(idStrings))

		for i, str := range idStrings {
			id, err := strconv.Atoi(strings.TrimSpace(str))
			if err != nil {
				context.String(http.StatusBadRequest, `"id" parameter contains invalid value. Should contain a comma separated list of integers`)
				context.Abort()
				return
			}
			ids[i] = id
		}
		filter.IDs = ids
	}

	near := strings.Split(query.Near, ",")
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
	sort, ok := SORT_FIELD_MAP[query.Sort]
	if len(query.Sort) > 0 && !ok {
		//TODO probably best to return the list of valid parameters here
		context.String(http.StatusBadRequest, "Sort parameter is not valid")
		context.Abort()
		return
	}

	rentals, err := RentalRepository.FindAllByFilter(filter, query.Offset, query.Limit, sort)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
	}

	context.IndentedJSON(http.StatusOK, rentals)
}
