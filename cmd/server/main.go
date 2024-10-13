package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"practicum-middle/config"
	"practicum-middle/pkg/handler"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	} else {
		fmt.Println("Environment variables loaded from .env")
	}

	opt := config.ParseFlags()

	router := gin.Default()

	h := handler.NewHandler(opt)

	router.POST("/", h.HandleShortenURL)
	router.GET("/:shortID", h.HandleRedirect)
	if err := router.Run(opt.FlagRunAddr); err != nil {
		panic(err)
	}
}
