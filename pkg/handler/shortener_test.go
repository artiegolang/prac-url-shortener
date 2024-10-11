package handler_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"practicum-middle/pkg/handler"
	"strings"
	"testing"
)

func TestHandleShortener(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        string
	}
	tests := []struct {
		name    string
		request *http.Request
		want    want
	}{
		{
			name: "Valid",
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader("http://google.com"))
				req.Header.Set("Content-Type", "text/plain")
				return req
			}(),
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				body:        "http://localhost:8085/" + handler.GenerateShortID("http://google.com"),
			},
		},
		{
			name: "Invalid",
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader("http://google.com"))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        "Bad Request\n",
			},
		},
		{
			name: "Empty",
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				req.Header.Set("Content-Type", "text/plain")
				return req
			}(),
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        "Bad Request\n",
			},
		},
		{
			name: "InvalidJSON",
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader("{"))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        "Bad Request\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.Mu.Lock()
			handler.UrlStore = make(map[string]string)
			handler.Mu.Unlock()

			w := httptest.NewRecorder()
			handler.HandleShortenURL(w, tt.request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.want.statusCode)

			assert.Equal(t, res.Header.Get("Content-Type"), tt.want.contentType)

			respBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.body, string(respBody))

			if tt.want.statusCode == http.StatusCreated {
				handler.Mu.Lock()
				defer handler.Mu.Unlock()

				shortId := handler.GenerateShortID("http://google.com")
				assert.Equal(t, "http://google.com", handler.UrlStore[shortId])
			}
		})
	}
}

func TestHandleRedirect(t *testing.T) {
	type want struct {
		status int
		header string
		body   string
	}
	tests := []struct {
		name     string
		urlStore map[string]string
		path     string
		want     want
	}{
		{
			name: "Valid",
			urlStore: map[string]string{
				"abc123": "google.com",
			},
			path: "/abc123",
			want: want{
				status: http.StatusFound,
				header: "/google.com",
			},
		},
		{
			name: "Invalid",
			urlStore: map[string]string{
				"abc123": "google.com",
			},
			path: "/abc12",
			want: want{
				status: http.StatusNotFound,
				header: "",
				body:   "Not Found\n",
			},
		},
		{
			name:     "Empty",
			urlStore: map[string]string{},
			path:     "/",
			want: want{
				status: http.StatusBadRequest,
				body:   "Bad Request\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.Mu.Lock()
			handler.UrlStore = tt.urlStore
			handler.Mu.Unlock()

			req, _ := http.NewRequest(http.MethodGet, tt.path, nil)

			w := httptest.NewRecorder()

			handler.HandleRedirect(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.want.status)

			if tt.want.status == http.StatusFound {
				assert.Equal(t, res.Header.Get("Location"), tt.want.header)
			}

			if tt.want.body != "" {
				respBody, err := io.ReadAll(res.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.want.body, string(respBody))
			}
		})
	}
}
