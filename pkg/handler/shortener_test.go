package handler_test

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"practicum-middle/pkg/handler"
	"strings"
	"testing"
)

func createTestContext(method, url string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(method, url, bytes.NewBuffer(body))
	return c, w
}

func TestHandleShortenURL(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		wantStatus  int
		wantBody    string
	}{
		{
			name:        "Valid",
			contentType: "text/plain",
			body:        "http://google.com",
			wantStatus:  http.StatusCreated,
			wantBody:    "http://localhost:8085/" + handler.GenerateShortID("http://google.com"),
		},
		{
			name:        "Invalid Content-Type",
			contentType: "application/json",
			body:        "http://google.com",
			wantStatus:  http.StatusBadRequest,
			wantBody:    "Bad Request",
		},
		{
			name:        "Empty Body",
			contentType: "text/plain",
			body:        "",
			wantStatus:  http.StatusBadRequest,
			wantBody:    "Bad Request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.Mu.Lock()
			handler.UrlStore = make(map[string]string)
			handler.Mu.Unlock()

			c, w := createTestContext(http.MethodPost, "/", []byte(tt.body))
			c.Request.Header.Set("Content-Type", tt.contentType)
			handler.HandleShortenURL(c)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)
			respBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.wantBody, string(respBody))
		})
	}
}

func TestHandleRedirect(t *testing.T) {
	tests := []struct {
		name       string
		urlStore   map[string]string
		path       string
		wantStatus int
		wantHeader string
		wantBody   string
	}{
		{
			name: "Valid",
			urlStore: map[string]string{
				"abc123": "http://google.com",
			},
			path:       "/abc123",
			wantStatus: http.StatusMovedPermanently,
			wantHeader: "http://google.com",
		},
		{
			name: "Invalid",
			urlStore: map[string]string{
				"abc123": "http://google.com",
			},
			path:       "/invalid",
			wantStatus: http.StatusNotFound,
			wantBody:   "Not Found",
		},
		{
			name:       "Empty",
			urlStore:   map[string]string{},
			path:       "/",
			wantStatus: http.StatusBadRequest,
			wantBody:   "Bad Request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.Mu.Lock()
			handler.UrlStore = tt.urlStore
			handler.Mu.Unlock()

			c, w := createTestContext(http.MethodGet, tt.path, nil)
			c.Params = gin.Params{gin.Param{Key: "shortID", Value: strings.TrimPrefix(tt.path, "/")}}
			handler.HandleRedirect(c)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)

			if tt.wantStatus == http.StatusMovedPermanently {
				assert.Equal(t, tt.wantHeader, res.Header.Get("Location"))
			}

			if tt.wantBody != "" {
				respBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, tt.wantBody, string(respBody))
			}
		})
	}
}
