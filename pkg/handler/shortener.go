package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"practicum-middle/config"
	"sync"
)

type URLRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Handler struct {
	opt             *config.Options
	logger          *zap.SugaredLogger
	mu              sync.Mutex
	UrlStore        map[string]string
	fileStoragePath string
}

func NewHandler(opt *config.Options, logger *zap.SugaredLogger, filestoragePath string) *Handler {
	h := &Handler{
		opt:             opt,
		logger:          logger,
		UrlStore:        make(map[string]string),
		fileStoragePath: opt.FileStoragePath,
	}
	h.loadURLsFromFile()
	return h
}

func (h *Handler) saveURLsToFile() {
	if h.fileStoragePath == "" {
		return
	}

	file, err := os.Create(h.fileStoragePath)
	if err != nil {
		h.logger.Errorf("Error creating file: %v", err)
		return
	}
	defer file.Close()

	h.mu.Lock()
	defer h.mu.Unlock()

	encoder := json.NewEncoder(file)
	for shortID, longURL := range h.UrlStore {
		record := URLRecord{
			UUID:        shortID,
			ShortURL:    h.opt.BaseURL + "/" + shortID,
			OriginalURL: longURL,
		}
		if err := encoder.Encode(&record); err != nil {
			h.logger.Errorf("Error encoding record: %v", err)
		}
	}
}

func (h *Handler) loadURLsFromFile() {
	if h.fileStoragePath == "" {
		return
	}

	file, err := os.Open(h.fileStoragePath)
	if err != nil {
		h.logger.Errorf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	h.mu.Lock()
	defer h.mu.Unlock()

	decoder := json.NewDecoder(file)

	for {
		var record URLRecord
		if err := decoder.Decode(&record); err == io.EOF {
			break
		} else if err != nil {
			h.logger.Errorf("Error decoding record: %v", err)
			break
		}
		h.UrlStore[record.UUID] = record.OriginalURL
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

	h.mu.Lock()
	h.UrlStore[shortID] = longURL
	h.mu.Unlock()

	shortURL := h.opt.BaseURL + "/" + shortID

	h.logger.Infof("Shortened URL %s for long URL %s", shortURL, longURL)
	c.JSON(http.StatusCreated, gin.H{"result": shortURL})

	h.saveURLsToFile()

}

func (h *Handler) HandleRedirect(c *gin.Context) {
	shortID := c.Param("shortID")

	if shortID == "" {
		h.logger.Infof("Bad Request: shortID is empty")
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}

	h.mu.Lock()
	longURL, ok := h.UrlStore[shortID]
	h.mu.Unlock()

	if !ok {
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

	h.mu.Lock()
	h.UrlStore[shortID] = longURL
	h.mu.Unlock()

	// shortURL := "http://localhost:8085/" + shortID

	shortURL := h.opt.BaseURL + "/" + shortID

	h.logger.Infof("Shortened URL %s for long URL %s", shortURL, longURL)
	c.String(http.StatusCreated, shortURL)

	h.saveURLsToFile()
}
