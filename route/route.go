package route

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "outdoorsy.com/backend/model"
)

func Init(router *gin.Engine) {
    router.GET("/rentals", GetRentals)
	router.GET("/rentals/:id", GetRental)
}


func GetRental(context *gin.Context) {
    id, err := strconv.Atoi(context.Param("id"))
    if err != nil {
        context.AbortWithStatus(http.StatusNotFound)
        return
    }
    context.IndentedJSON(http.StatusOK, model.Rental{Id: id})
}

func GetRentals(context *gin.Context) {
    context.IndentedJSON(http.StatusOK, make([]model.Rental,10))
}