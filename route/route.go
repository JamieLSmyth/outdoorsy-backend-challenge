package route

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
    router.GET("/rentals", GetRentals)
	router.GET("/rentals/:id", GetRental)
}


func GetRentals(context *gin.Context) {
    context.IndentedJSON(http.StatusOK, struct {}{})
}

func GetRental(context *gin.Context) {
    context.IndentedJSON(http.StatusOK, struct {}{})
}