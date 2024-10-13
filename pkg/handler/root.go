package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RootHandler(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodPost:
		if c.Request.URL.Path == "/" {
			HandleShortenURL(c)
		} else {
			c.String(http.StatusBadRequest, "Bad Request")
		}
	case http.MethodGet:
		HandleRedirect(c)
	default:
		c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
	}
}
