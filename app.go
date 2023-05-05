package main

import (
	"github.com/gin-gonic/gin"
	"outdoorsy.com/backend/route"
)

func main() {

	router := gin.Default()
	route.Init(router)

	router.Run(":8080") // TODO: make this configurable

}