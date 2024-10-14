package main

import (
	"practicum-middle/internal/controller"
)

func main() {

	// Initialize the base controller
	baseController := controller.NewBaseController()
	defer baseController.Logger.Sync()

	// Run the server
	if err := baseController.Run(); err != nil {
		panic(err)
	}
}
