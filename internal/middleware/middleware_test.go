package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGzipDecompress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(GzipDecompress())
	router.POST("/test", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.String(http.StatusOK, string(body))
	})

	// Create a gzip compressed request body
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, _ = gz.Write([]byte(`{"key":"value"}`))
	gz.Close()

	req, _ := http.NewRequest(http.MethodPost, "/test", &buf)
	req.Header.Set("Content-Encoding", "gzip")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"key":"value"}`, w.Body.String())
}

func TestGzipCompress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(GzipCompress())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"key": "value"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))

	// Decompress the response body
	gz, err := gzip.NewReader(w.Body)
	assert.NoError(t, err)
	defer gz.Close()
	body, err := io.ReadAll(gz)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"key":"value"}`, string(body))
}
