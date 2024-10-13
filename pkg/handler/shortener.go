package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"practicum-middle/config"
	"sync"
)

var (
	UrlStore = make(map[string]string)
	Mu       sync.RWMutex
)

type Handler struct {
	opt *config.Options
}

func NewHandler(opt *config.Options) *Handler {
	return &Handler{opt: opt}
}

func (h *Handler) HandleRedirect(c *gin.Context) {
	shortID := c.Param("shortID")

	if shortID == "" {
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}

	Mu.RLock()
	longURL, ok := UrlStore[shortID]
	Mu.RUnlock()

	if !ok {
		c.String(http.StatusNotFound, "Not Found")
		return
	}

	c.Redirect(http.StatusMovedPermanently, longURL)
}

func (h *Handler) HandleShortenURL(c *gin.Context) {
	if c.GetHeader("Content-Type") != "text/plain" {
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}
	longURL := string(body)

	shortID := GenerateShortID(longURL)

	Mu.Lock()
	UrlStore[shortID] = longURL
	Mu.Unlock()

	// shortURL := "http://localhost:8085/" + shortID

	shortURL := h.opt.BaseURL + "/" + shortID

	c.String(http.StatusCreated, shortURL)
}
