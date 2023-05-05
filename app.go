package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"outdoorsy.com/backend/route"
)

func main() {

	router := gin.Default()
	route.Init(router)
	// get application host and port from env so they can be configured for each environment in necessary
	host := os.Getenv("APP_HOST")// Defaults to all network interfaces
	port := os.Getenv("APP_PORT")
	if len(port) == 0 {
		port = "8080" //Default to 8080
	}
	router.Run(fmt.Sprintf("%s:%s", host, port)) // TODO: make this configurable

}