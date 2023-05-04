package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.Run("localhost:8080") // TODO: make this configurable externally

}