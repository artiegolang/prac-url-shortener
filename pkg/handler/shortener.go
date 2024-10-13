package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"sync"
)

var (
	UrlStore = make(map[string]string)
	Mu       sync.RWMutex
)

func HandleRedirect(c *gin.Context) {
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

func HandleShortenURL(c *gin.Context) {
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

	shortURL := "http://localhost:8085/" + shortID

	c.String(http.StatusCreated, shortURL)
}
