package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"practicum-middle/config"
	"practicum-middle/internal/controller"
	"practicum-middle/pkg/database"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения из .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	} else {
		fmt.Println("Environment variables loaded from .env")
	}

	opt := config.ParseConfig()

	// Теперь переменные окружения доступны для database.NewDB()
	db := database.NewDB()
	defer db.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("Received signal: %v\n", sig)
		db.Close()
		os.Exit(0)
	}()

	// Initialize the base controller
	baseController := controller.NewBaseController(db, opt)
	defer baseController.Logger.Sync()

	// Run the server
	if err := baseController.Run(); err != nil {
		panic(err)
	}
}
