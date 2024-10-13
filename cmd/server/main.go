package main

import (
	"github.com/gin-gonic/gin"
	"practicum-middle/pkg/handler"
)

func main() {
	router := gin.Default()

	router.POST("/", handler.HandleShortenURL)
	router.GET("/:shortID", handler.HandleRedirect)
	if err := router.Run(":8085"); err != nil {
		panic(err)
	}
}
