package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"practicum-middle/config"
	"practicum-middle/internal/logger"
	"practicum-middle/internal/middleware"
	"practicum-middle/internal/repository"
	"practicum-middle/pkg/database"
	"practicum-middle/pkg/handler"
)

type BaseController struct {
	Router *gin.Engine
	Opt    *config.Options
	Logger *zap.SugaredLogger
}

func NewBaseController(db *database.DB, opt *config.Options) *BaseController {
	// Инициализация логгера
	log := logger.NewLogger()

	// Инициализация репозитория
	urlRepo := repository.NewURLRepository(db)

	// Инициализация хендлера
	h := handler.NewHandler(opt, log, urlRepo)

	// Инициализация роутера
	router := gin.New()
	router.Use(middleware.GzipDecompress())
	router.Use(middleware.GzipCompress())
	router.Use(middleware.RequestLogger(log))
	router.Use(gin.Recovery())

	// Настройка маршрутов
	router.POST("/", h.HandleShortenURL)
	router.GET("/:shortID", h.HandleRedirect)
	router.POST("/api/shorten", h.HandleShortenURLJSON)
	router.GET("/ping", h.PingDB)
	router.POST("/api/shorten/batch", h.HandleShortenURLBatch)

	return &BaseController{
		Router: router,
		Opt:    opt,
		Logger: log,
	}
}

func (bc *BaseController) Run() error {
	bc.Logger.Infof("Starting the server on %s", bc.Opt.ServerAddress)
	return bc.Router.Run(bc.Opt.ServerAddress)
}
