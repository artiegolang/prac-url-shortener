package handler

import "C"
import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"practicum-middle/config"
	"practicum-middle/internal/repository"
)

type URLRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Handler struct {
	opt     *config.Options
	logger  *zap.SugaredLogger
	urlRepo *repository.URLRepository
}

func NewHandler(opt *config.Options, logger *zap.SugaredLogger, urlRepo *repository.URLRepository) *Handler {
	return &Handler{
		opt:     opt,
		logger:  logger,
		urlRepo: urlRepo,
	}
}

func (h *Handler) HandleShortenURLJSON(c *gin.Context) {
	var requestBody struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		h.logger.Infof("Bad Request: Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	longURL := requestBody.URL
	if longURL == "" {
		h.logger.Infof("Bad Request: URL is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}
	shortID := GenerateShortID(longURL)

	// Сохранение в базу данных через репозиторий
	existingShortID, exists, err := h.urlRepo.SaveURL(c.Request.Context(), shortID, longURL)
	if err != nil {
		h.logger.Errorf("Error saving URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	shortURL := h.opt.BaseURL + "/" + existingShortID

	if exists {
		// URL уже существует
		h.logger.Infof("URL already exists for long URL %s", longURL)
		c.JSON(http.StatusConflict, gin.H{"result": shortURL})
		return
	}

	h.logger.Infof("Shortened URL %s for long URL %s", shortURL, longURL)
	c.JSON(http.StatusCreated, gin.H{"result": shortURL})
}

func (h *Handler) HandleRedirect(c *gin.Context) {
	shortID := c.Param("shortID")

	if shortID == "" {
		h.logger.Infof("Bad Request: shortID is empty")
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}

	// Получаем longURL из репозитория
	longURL, err := h.urlRepo.GetOriginalURL(c.Request.Context(), shortID)
	if err != nil {
		h.logger.Infof("Not Found: shortID %s not found", shortID)
		c.String(http.StatusNotFound, "Not Found")
		return
	}

	h.logger.Infof("Redirecting to %s for shortID %s", longURL, shortID)
	c.Redirect(http.StatusMovedPermanently, longURL)
}

func (h *Handler) HandleShortenURL(c *gin.Context) {
	if c.GetHeader("Content-Type") != "text/plain" {
		h.logger.Infof("Bad Request: Content-Type is not text/plain")
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		h.logger.Infof("Bad Request: Error reading body or body is empty")
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}
	longURL := string(body)

	shortID := GenerateShortID(longURL)

	// Сохранение в базу данных
	existingShortID, exists, err := h.urlRepo.SaveURL(c.Request.Context(), shortID, longURL)
	if err != nil {
		h.logger.Errorf("Error saving URL: %v", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	shortURL := h.opt.BaseURL + "/" + existingShortID

	if exists {
		// URL уже существует
		h.logger.Infof("URL already exists for long URL %s", longURL)
		c.String(http.StatusConflict, shortURL)
		return
	}

	h.logger.Infof("Shortened URL %s for long URL %s", shortURL, longURL)
	c.String(http.StatusCreated, shortURL)
}

func (h *Handler) PingDB(c *gin.Context) {
	err := h.urlRepo.Ping(c.Request.Context())
	if err != nil {
		h.logger.Errorf("Error pinging database: %v", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.String(http.StatusOK, "OK")

}

func (h *Handler) HandleShortenURLBatch(c *gin.Context) {
	var requestBody []struct {
		CorrelationID string `json:"correlation_id" binding:"required"`
		URL           string `json:"original_url" binding:"required"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		h.logger.Infof("Bad Request: Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	if len(requestBody) == 0 {
		h.logger.Infof("Bad Request: Empty batch")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	urlPairs := make([]repository.URLPair, len(requestBody))
	for i, record := range requestBody {
		shortID := GenerateShortID(record.URL)
		urlPairs[i] = repository.URLPair{ShortID: shortID, OriginalURL: record.URL}
	}

	urlMap, err := h.urlRepo.SaveURLsBatch(c.Request.Context(), urlPairs)
	if err != nil {
		h.logger.Errorf("Error saving URLs batch: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	result := make([]struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}, len(requestBody))

	for i, record := range requestBody {
		shortID := urlMap[record.URL]
		shortURL := h.opt.BaseURL + "/" + shortID
		result[i] = struct {
			CorrelationID string `json:"correlation_id"`
			ShortURL      string `json:"short_url"`
		}{
			CorrelationID: record.CorrelationID,
			ShortURL:      shortURL,
		}
	}

	c.JSON(http.StatusCreated, result)
}
