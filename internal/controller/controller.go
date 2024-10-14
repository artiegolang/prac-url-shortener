package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"practicum-middle/config"
	"practicum-middle/internal/logger"
	"practicum-middle/internal/middleware"
	"practicum-middle/pkg/handler"
)

type BaseController struct {
	Router *gin.Engine
	Opt    *config.Options
	Logger *zap.SugaredLogger
}

func NewBaseController() *BaseController {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	} else {
		fmt.Println("Environment variables loaded from .env")
	}

	// Parse flags and environment variables
	opt := config.ParseFlags()

	// Initialize the logger
	log := logger.NewLogger()

	// Initialize the router
	router := gin.New()

	router.Use(middleware.GzipDecompress())
	router.Use(middleware.GzipCompress())

	// Add the request logger middleware
	router.Use(middleware.RequestLogger(log))
	router.Use(gin.Recovery())

	// Initialize handlers
	h := handler.NewHandler(opt, log, opt.FileStoragePath)

	// Set up routes
	router.POST("/", h.HandleShortenURL)
	router.GET("/:shortID", h.HandleRedirect)
	router.POST("/api/shorten", h.HandleShortenURLJSON)

	return &BaseController{
		Router: router,
		Opt:    opt,
		Logger: log,
	}
}

func (bc *BaseController) Run() error {
	bc.Logger.Info("Starting the server on %s", bc.Opt.FlagRunAddr)
	return bc.Router.Run(bc.Opt.FlagRunAddr)
}
